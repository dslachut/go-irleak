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

// Create the tables

const mysql_createAuth = `CREATE TABLE IF NOT EXISTS auth (
	user TINYTEXT PRIMARY KEY,
	hash VARBINARY(128) NOT NULL
)`

const mysql_createToken = `CREATE TABLE IF NOT EXISTS tokens (
	user  TINYTEXT REFERENCES auth (user),
	token CHAR(32) UNIQUE NOT NULL,
	exp   BIGINT NOT NULL
)`

const mysql_createTemperature = `CREATE TABLE IF NOT EXISTS temperatures(
	user        TINYTEXT REFERENCES auth (user),
	sensor      TINYTEXT NOT NULL,
	timestamp   DOUBLE NOT NULL,
	value       DOUBLE NOT NULL,
	PRIMARY KEY (user,sensor,timestamp)
)`

// auth functions

const mysql_getHash = `SELECT hash FROM auth WHERE user=?`

const mysql_addToken = `INSERT INTO tokens VALUES (?, ?, ?)`

const mysql_addUser = `INSERT INTO auth VALUES (?, ?)`

const mysql_getUser = `SELECT user, exp FROM tokens WHERE token = ?`

const mysql_purgeTokens = `DELETE FROM tokens WHERE exp < ?`

// data functions

const mysql_addTemperature = `INSERT INTO temperatures VALUES (?, ?, ?, ?)`
