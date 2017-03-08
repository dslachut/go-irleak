// Copyright © 2017 David Lachut <dslachut@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kb

import (
	"database/sql"
	"log"
	s "strings"

	_ "github.github.com/mattn/go-sqlite3"
)

type sqliteKB struct {
	inbound chan *query
}

type query struct {
	queryString string
	arguments   []interface{}
	rows        chan []interface{}
	result      chan sql.Result
}

func NewSQliteKB(dbFile string, sqliteOpts map[string]string) *sqliteKB {
	if sqliteOpts != nil && len(sqliteOpts) > 0 {
		dbFile = s.Join([]string{"file:", "?"}, dbFile)
		args := make([]string, 0, len(sqliteOpts))
		for k, v := range sqliteOpts {
			append(args, s.Join([]string{k, v}, "="))
		}
		allArgs = s.Join(args, "&")
		dbFile = s.Join([]string{dbFile, allArgs}, "")
	}

	kb := &sqliteKB{inbound: make(chan *query)}
	go kbLoop(dbFile, kb)
	return kb
}

func kbLoop(dbFile string, kb *sqliteKB) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for {
		q, ok := <-kb.inbound
		if !ok {
			break
		}

		if q.rows == nil && q.result != nil {
			doInsert(db, q)
		} else if q.rows != nil && q.result == nil {
			doQuery(db, q)
		} else {
			continue
		}
	}
}

func doInsert(db *sql.DB, q *query) {
	stmt, err := db.Prepare(q.queryString)
	if err != nil {
		Close(q.result)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(q.arguments...)
	if err != nil {
		Close(q.result)
		return
	}
}
