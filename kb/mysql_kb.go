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
	"fmt"
	"log"
	s "strings"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlKB struct {
	db *sql.DB
}

func NewMysqlKB(user, password, dbname string, params map[string]string) *mysqlKB {
	k := new(mysqlKB)
	var err error
	if params != nil {
		args := make([]string, 0, len(params))
		for k, v := range params {
			args = append(args, s.Join([]string{k, v}, "="))
		}
		allArgs := s.Join(args, "&")
		k.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?%s", user, password, dbname, allArgs))
	} else {
		k.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, password, dbname))
	}
	if err != nil {
		log.Fatal(err)
	}
	initMysqlDB(k.db)
	return k
}

func initMysqlDB(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("tx start")
		log.Fatal(err)
	}
	defer tx.Commit()

	_, err = tx.Exec(mysql_createAuth)
	if err != nil {
		log.Println("create auth")
		log.Fatal(err)
	}
	_, err = tx.Exec(mysql_createToken)
	if err != nil {
		log.Println("create token")
		log.Fatal(err)
	}

	_, err = tx.Exec(mysql_createTemperature)
	if err != nil {
		log.Println("create temp")
		log.Fatal(err)
	}
}

func (k *mysqlKB) Stop() {
	k.db.Close()
}

func (k *mysqlKB) GetHash(user string) ([]byte, bool) {
	q := &query{
		queryString: mysql_getHash,
		arguments:   []interface{}{user},
		rows:        make(chan []map[string]interface{}),
		result:      nil,
	}
	go doQuery(k.db, q)
	rows := <-q.rows
	if len(rows) == 1 {
		return rows[0]["hashval"].([]byte), true
	} else {
		return nil, false
	}
}

func (k *mysqlKB) AddToken(user, token string, expiration int64) bool {
	q := &query{
		queryString: mysql_addToken,
		arguments:   []interface{}{user, token, expiration},
		rows:        nil,
		result:      make(chan sql.Result),
	}
	go doInsert(k.db, q)
	_, err := (<-q.result).RowsAffected()
	if err != nil {
		return false
	}
	return true
}

func (k *mysqlKB) AddUser(user string, hash string) bool {
	q := &query{
		queryString: mysql_addUser,
		arguments:   []interface{}{user, hash},
		rows:        nil,
		result:      make(chan sql.Result),
	}
	go doInsert(k.db, q)

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

func (k *mysqlKB) GetUser(token string) (string, int64, bool) {
	q := &query{
		queryString: mysql_getUser,
		arguments:   []interface{}{token},
		rows:        make(chan []map[string]interface{}),
		result:      nil,
	}
	go doQuery(k.db, q)
	rows := <-q.rows
	if len(rows) == 1 {
		//              log.Printf("%T.%T\n", rows[0]["user"], rows[0]["exp"])
		return string(rows[0]["uname"].([]byte)), rows[0]["exp"].(int64), true
	} else {
		return "", 0, false
	}
}

func (k *mysqlKB) ExpireToken(token string) bool {
	q := &query{
		queryString: mysql_expireToken,
		arguments:   []interface{}{token},
		rows:        nil,
		result:      make(chan sql.Result),
	}
	go doInsert(k.db, q)

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

func (k *mysqlKB) PurgeTokens(expiration int64) bool {
	q := &query{
		queryString: mysql_purgeTokens,
		arguments:   []interface{}{expiration},
		rows:        nil,
		result:      make(chan sql.Result),
	}
	go doInsert(k.db, q)

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

func (k *mysqlKB) AddTemperature(user, sensor string, timestamp, value float64) bool {
	q := &query{
		queryString: mysql_addTemperature,
		arguments:   []interface{}{user, sensor, timestamp, value},
		rows:        nil,
		result:      make(chan sql.Result),
	}
	go doInsert(k.db, q)

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
