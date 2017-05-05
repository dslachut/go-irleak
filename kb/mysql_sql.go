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
	uname       VARCHAR(255) NOT NULL,
	hashval     VARBINARY(128) NOT NULL,
	PRIMARY KEY (uname)
);`

const mysql_createToken = `CREATE TABLE IF NOT EXISTS tokens (
	uname VARCHAR(255) REFERENCES auth (uname),
	token CHAR(32) UNIQUE NOT NULL,
	exp   BIGINT NOT NULL
)`

const mysql_createTemperature = `CREATE TABLE IF NOT EXISTS temperatures(
	uname       VARCHAR(255) REFERENCES auth (uname),
	sensor      VARCHAR(255) NOT NULL,
	timestamp   DOUBLE NOT NULL,
	value       DOUBLE NOT NULL,
	PRIMARY KEY (uname,sensor,timestamp)
)`

const mysql_createLocation = `CREATE TABLE IF NOT EXISTS location(
	l_id        INTEGER AUTO_INCREMENT,
	uname       VARCHAR(255) REFERENCES auth (uname),
	place_name  VARCHAR(255) NOT NULL,
	lat         VARCHAR(12) NOT NULL,
	long        VARCHAR(12) NOT NULL,
	PRIMARY KEY (l_id),
	UNIQUE KEY  (uname,place_name)
)`

const mysql_createWeather = `CREATE TABLE IF NOT EXISTS weather(
	l_id                 INTEGER REFERENCES location (l_id),
	timestamp            DOUBLE,
	sun_up               BOOLEAN,
	temperature          DOUBLE,
	apparent_temperature DOUBLE,
	cloud_cover          DOUBLE,
	humidity             DOUBLE,
	pressure             DOUBLE,
	precib_probability   DOUBLE,
	PRIMARY KEY          (l_id, timestamp)
)`

// auth functions

const mysql_getHash = `SELECT hashval FROM auth WHERE uname=?`

const mysql_addToken = `INSERT INTO tokens VALUES (?, ?, ?)`

const mysql_addUser = `INSERT INTO auth VALUES (?, ?)`

const mysql_getUser = `SELECT uname, exp FROM tokens WHERE token = ?`

const mysql_expireToken = `UPDATE tokens SET exp=0 where token=?`

const mysql_purgeTokens = `DELETE FROM tokens WHERE exp < ?`

// data functions

const mysql_addTemperature = `INSERT IGNORE INTO temperatures VALUES (?, ?, ?, ?)`

const mysql_addWeather = `REPLACE INTO weather VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

// location functions

const mysql_getCoordinates = `SELECT l_id, lat, long FROM location`
const mysql_addLocation = `INSERT INTO location VALUES (?, ?, ?, ?)`
