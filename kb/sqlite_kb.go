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
	done    chan bool
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
		select {
		case q := <-kb.inbound:
			if q.result != nil {
				doInsert(db, q)
			} else if q.rows != nil {
				doQuery(db, q)
			} else {
				continue
			}
		case <-kb.done:
			return
		}
	}
}

func doInsert(db *sql.DB, q *query) {
	stmt, err := db.Prepare(q.queryString)
	if err != nil {
		log.Println(err)
		close(q.result)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(q.arguments...)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(q.arguments...)
	if err != nil {
		close(q.rows)
		log.Println(err)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		close(q.rows)
		log.Println(err)
		return
	}

	outRows := make([]map[string]interface{}, 0)
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

	_, err = tx.Exec(sqlite_createTemperature)
	if err != nil {
		log.Fatal(err)
	}
}

func (k *sqliteKB) Stop() {
	close(k.done)
}

func (k *sqliteKB) GetHash(user string) ([]byte, bool) {
	q := &query{
		queryString: sqlite_getHash,
		arguments:   []interface{}{user},
		rows:        make(chan []map[string]interface{}),
		result:      nil,
	}
	k.inbound <- q
	rows := <-q.rows
	if len(rows) == 1 {
		return rows[0]["hash"].([]byte), true
	} else {
		return nil, false
	}
}

func (k *sqliteKB) AddToken(user, token string, expiration int64) bool {
	q := &query{
		queryString: sqlite_addToken,
		arguments:   []interface{}{user, token, expiration},
		rows:        nil,
		result:      make(chan sql.Result),
	}
	k.inbound <- q
	_, err := (<-q.result).RowsAffected()
	if err != nil {
		return false
	}
	return true
}

func (k *sqliteKB) AddUser(user, hash string) bool {
	q := &query{
		queryString: sqlite_addUser,
		arguments:   []interface{}{user, hash},
		rows:        nil,
		result:      make(chan sql.Result),
	}
	k.inbound <- q

	res, ok := <-q.result
	if !ok {
		return false
	}

	_, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return true
}

func (k *sqliteKB) GetUser(token string) (string, int64, bool) {
	q := &query{
		queryString: sqlite_getUser,
		arguments:   []interface{}{token},
		rows:        make(chan []map[string]interface{}),
		result:      nil,
	}
	k.inbound <- q
	rows := <-q.rows
	if len(rows) == 1 {
		//		log.Printf("%T.%T\n", rows[0]["user"], rows[0]["exp"])
		return string(rows[0]["user"].([]byte)), rows[0]["exp"].(int64), true
	} else {
		return "", 0, false
	}
}

func (k *sqliteKB) AddTemperature(user, sensor string, timestamp, value float64) bool {
	q := &query{
		queryString: sqlite_addTemperature,
		arguments:   []interface{}{user, sensor, timestamp, value},
		rows:        nil,
		result:      make(chan sql.Result),
	}
	k.inbound <- q

	res, ok := <-q.result
	if !ok {
		return false
	}

	_, err := res.RowsAffected()
	if err != nil {
		return false
	}

	return true
}
