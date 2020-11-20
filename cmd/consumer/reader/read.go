/*
  Copyright (C) 2019 - 2021 MWSOFT
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
	"go.uber.org/zap"
	"time"

	"github.com/superhero-match/consumer-match-delete/internal/consumer/model"
	dbm "github.com/superhero-match/consumer-match-delete/internal/db/model"
	fm "github.com/superhero-match/consumer-match-delete/internal/firebase/model"
)

// Read consumes the Kafka topic and stores the choice made by superhero
// to DB and if it is a like to Cache as well.
func (r *Reader) Read() error {
	ctx := context.Background()

	for {
		fmt.Println("before FetchMessage")
		m, err := r.Consumer.Consumer.FetchMessage(ctx)
		fmt.Println("after FetchMessage")
		if err != nil {
			r.Logger.Error(
				"failed to fetch message",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)

			err = r.Consumer.Consumer.Close()
			if err != nil {
				r.Logger.Error(
					"failed to close consumer",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

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
				r.Logger.Error(
					"failed to unmarshal JSON to Match consumer model",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				err = r.Consumer.Consumer.Close()
				if err != nil {
					r.Logger.Error(
						"failed to close consumer",
						zap.String("err", err.Error()),
						zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
					)

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
			r.Logger.Error(
				"failed to delete match from database",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)

			err = r.Consumer.Consumer.Close()
			if err != nil {
				r.Logger.Error(
					"failed to close consumer",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				return err
			}

			return err
		}

		token, err := r.Cache.GetFirebaseMessagingToken(fmt.Sprintf(r.Cache.TokenKeyFormat, match.MatchedSuperheroID))
		if err != nil || token == nil {
			r.Logger.Error(
				"failed to fetch Firebase token from cache",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)

			err = r.Consumer.Consumer.Close()
			if err != nil {
				r.Logger.Error(
					"failed to close consumer",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				return err
			}

			return err
		}

		err = r.Firebase.PushDeleteMatchNotification(fm.Request{
			Token:       token.Token,
			SuperheroID: match.SuperheroID,
		})
		if err != nil {
			r.Logger.Error(
				"failed to push Firebase notification to delete match",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)

			err = r.Consumer.Consumer.Close()
			if err != nil {
				r.Logger.Error(
					"failed to close consumer",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				return err
			}

			return err
		}

		err = r.Consumer.Consumer.CommitMessages(ctx, m)
		if err != nil {
			r.Logger.Error(
				"failed to commit message",
				zap.String("err", err.Error()),
				zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
			)
			
			err = r.Consumer.Consumer.Close()
			if err != nil {
				r.Logger.Error(
					"failed to close consumer",
					zap.String("err", err.Error()),
					zap.String("time", time.Now().UTC().Format(r.TimeFormat)),
				)

				return err
			}

			return err
		}
	}
}
