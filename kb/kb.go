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
)

type KB interface {
	Stop()
	GetHash(user string) ([]byte, bool)
	AddToken(user string, token string, expiration int64) bool
	AddUser(user string, hash string) bool
	GetUser(token string) (string, int64, bool)
	ExpireToken(token string) bool
	PurgeTokens(expiration int64) bool

	AddTemperature(user, sensor string, timestamp, value float64) bool
	AddWeather(location, timestamp int64, sunUp bool, temperature, apparentTemperature, cloudCover, humidity, pressure, precipProbability float64) bool

	GetCoordinates() ([][]string, []int64, bool)
}

type query struct {
	queryString string
	arguments   []interface{}
	rows        chan []map[string]interface{}
	result      chan sql.Result
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
