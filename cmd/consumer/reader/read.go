/*
  Copyright (C) 2019 - 2020 MWSOFT
  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.
  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.
  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/superhero-match/consumer-match-delete/internal/consumer/model"
	dbm "github.com/superhero-match/consumer-match-delete/internal/db/model"
	fm "github.com/superhero-match/consumer-match-delete/internal/firebase/model"
)

const (
	timeFormat = "2006-01-02T15:04:05"
)

// Read consumes the Kafka topic and stores the choice made by superhero
// to DB and if it is a like to Cache as well.
func (r *Reader) Read() error {
	ctx := context.Background()

	for {
		fmt.Println("before FetchMessage")
		m, err := r.Consumer.Consumer.FetchMessage(ctx)
		fmt.Println("after FetchMessage")
		fmt.Println("err: ")
		fmt.Println(err)
		if err != nil {
			err = r.Consumer.Consumer.Close()
			if err != nil {
				return err
			}

			return err
		}

		fmt.Printf(
			"message at topic/partition/offset \n%v/\n%v/\n%v: \n%s = \n%s\n",
			m.Topic,
			m.Partition,
			m.Offset,
			string(m.Key),
			string(m.Value),
		)

		var match model.Match

		if err := json.Unmarshal(m.Value, &match); err != nil {
			_ = r.Consumer.Consumer.Close()
			if err != nil {
				fmt.Println("Unmarshal")
				fmt.Println(err)
				err = r.Consumer.Consumer.Close()
				if err != nil {
					return err
				}

				return err
			}
		}

		err = r.DB.DeleteMatch(dbm.Match{
			SuperheroID:        match.SuperheroID,
			MatchedSuperheroID: match.MatchedSuperheroID,
			DeletedAt:          time.Now().UTC().Format(timeFormat),
		}, )
		if err != nil {
			fmt.Println("DB")
			fmt.Println(err)
			err = r.Consumer.Consumer.Close()
			if err != nil {
				return err
			}

			return err
		}

		token, err := r.Cache.GetFirebaseMessagingToken(fmt.Sprintf(r.Cache.TokenKeyFormat, match.MatchedSuperheroID))
		if err != nil || token == nil {
			fmt.Println("r.Cache.GetFirebaseMessagingToken")
			fmt.Println(err)
			fmt.Println(token)

			err = r.Consumer.Consumer.Close()
			if err != nil {
				return err
			}

			return err
		}

		err = r.Firebase.PushDeleteMatchNotification(fm.Request{
			Token:       token.Token,
			SuperheroID: match.SuperheroID,
		})
		if err != nil {
			fmt.Println("r.Firebase.PushNewMatchNotification")
			fmt.Println(err)

			err = r.Consumer.Consumer.Close()
			if err != nil {
				return err
			}

			return err
		}

		err = r.Consumer.Consumer.CommitMessages(ctx, m)
		if err != nil {
			err = r.Consumer.Consumer.Close()
			if err != nil {
				return err
			}

			return err
		}
	}
}
