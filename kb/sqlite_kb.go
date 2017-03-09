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

	_ "github.com/mattn/go-sqlite3"
)

type sqliteKB struct {
	inbound chan *query
}

type query struct {
	queryString string
	arguments   []interface{}
	rows        chan []map[string]interface{}
	result      chan sql.Result
}

func NewSQLiteKB(dbFile string, sqliteOpts map[string]string) KB {
	if sqliteOpts != nil && len(sqliteOpts) > 0 {
		dbFile = s.Join([]string{"file:", "?"}, dbFile)
		args := make([]string, 0, len(sqliteOpts))
		for k, v := range sqliteOpts {
			args = append(args, s.Join([]string{k, v}, "="))
		}
		allArgs := s.Join(args, "&")
		dbFile = s.Join([]string{dbFile, allArgs}, "")
	}

	newKB := &sqliteKB{inbound: make(chan *query)}
	go kbLoop(dbFile, newKB)
	return newKB
}

func kbLoop(dbFile string, kb *sqliteKB) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDB(db)

	for {
		q, ok := <-kb.inbound
		if !ok {
			break
		}

		if q.result != nil {
			doInsert(db, q)
		} else if q.rows != nil {
			doQuery(db, q)
		} else {
			continue
		}
	}
}

func doInsert(db *sql.DB, q *query) {
	stmt, err := db.Prepare(q.queryString)
	if err != nil {
		close(q.result)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(q.arguments...)
	if err != nil {
		close(q.result)
		return
	}

	q.result <- res
	close(q.result)
}

func doQuery(db *sql.DB, q *query) {
	stmt, err := db.Prepare(q.queryString)
	if err != nil {
		close(q.rows)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(q.arguments...)
	if err != nil {
		close(q.rows)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		close(q.rows)
		return
	}

	outRows := make([]map[string]interface{}, 1)
	for rows.Next() {
		newRow := make(map[string]interface{})
		vals := make([]interface{}, len(cols))
		valPtrs := make([]interface{}, len(cols))
		for i, _ := range cols {
			valPtrs[i] = &vals[i]
		}
		rows.Scan(valPtrs...)
		for i, col := range cols {
			newRow[col] = vals[i]
		}
		outRows = append(outRows, newRow)
	}
	q.rows <- outRows
	close(q.rows)
}

func initDB(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Commit()

	_, err = tx.Exec(sqlite_createAuth)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec(sqlite_createToken)
	if err != nil {
		log.Fatal(err)
	}
}

func (k *sqliteKB) GetHash(user string) string {
	q := new(query)
	q.queryString = sqlite_getHash
	q.arguments = []interface{}{user}
	q.rows = make(chan []map[string]interface{})
	q.result = nil
	k.inbound <- q
	rows := <-q.rows
	return rows[0]["hash"].(string)
}
