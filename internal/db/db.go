/*
  Copyright (C) 2019 - 2022 MWSOFT
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
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver.

	"github.com/superhero-match/consumer-match-delete/internal/config"
	"github.com/superhero-match/consumer-match-delete/internal/db/model"
)

// DB interface defines database methods.
type DB interface {
	DeleteMatch(m model.Match) error
}

// db holds the database connection.
type db struct {
	DB              *sql.DB
	stmtDeleteMatch *sql.Stmt
}

// NewDB returns database.
func NewDB(cfg *config.Config) (dbs DB, err error) {
	ds, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Name,
		),
	)
	if err != nil {
		return nil, err
	}

	stmtDel, err := ds.Prepare(`call delete_match(?,?,?)`)
	if err != nil {
		return nil, err
	}

	return &db{
		DB:              ds,
		stmtDeleteMatch: stmtDel,
	}, nil
}
