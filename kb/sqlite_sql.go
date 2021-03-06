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

const sqlite_createAuth = `CREATE TABLE IF NOT EXISTS auth (
	user TEXT PRIMARY KEY,
	hash TEXT NOT NULL
)`

const sqlite_createToken = `CREATE TABLE IF NOT EXISTS tokens (
	user TEXT REFERENCES auth (user),
	token TEXT UNIQUE NOT NULL,
	exp INTEGER NOT NULL
)`

const sqlite_createTemperature = `CREATE TABLE IF NOT EXISTS temperatures(
	user TEXT REFERENCES auth (user),
	sensor TEXT NOT NULL,
	timestamp NUMERIC NOT NULL,
	value NUMERIC NOT NULL,
	PRIMARY KEY (user,sensor,timestamp)
)`

const sqlite_createLocation = `CREATE TABLE IF NOT EXISTS location(
	l_id INTEGER AUTOINCREMENT,
	user TEXT REFERENCES auth (user),
	place_name TEXT NOT NULL,
	lat TEXT NOT NULL,
	lon TEXT NOT NULL,
	PRIMARY KEY (l_id),
	UNIQUE (user, place_name)
)`

const sqlite_createWeather = `CREATE TABLE IF NOT EXISTS weather(
	l_id                 INTEGER REFERENCES location (l_id),
	timestamp            NUMERIC,
	sun_up               BOOLEAN,
	temperature          NUMERIC,
	apparent_temperature NUMERIC,
	cloud_cover          NUMERIC,
	humidity             NUMERIC,
	pressure             NUMERIC,
	precib_probability   NUMERIC,
	PRIMARY KEY          (l_id, timestamp)
)`

// auth functions

const sqlite_getHash = `SELECT hash FROM auth WHERE user=?`

const sqlite_addToken = `INSERT INTO tokens VALUES (?, ?, ?)`

const sqlite_addUser = `INSERT INTO auth VALUES (?, ?)`

const sqlite_getUser = `SELECT user, exp FROM tokens WHERE token = ?`

const sqlite_expireToken = `UPDATE tokens SET exp=0 where token=?`

const sqlite_purgeTokens = `DELETE FROM tokens WHERE exp < ?`

// data functions

const sqlite_addTemperature = `INSERT OR IGNORE INTO temperatures VALUES (?, ?, ?, ?)`

const sqlite_addWeather = `REPLACE INTO weather VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

const sqlite_getTemperatures = `SELECT timestamp, value FROM temperatures WHERE user=? and sensor=? and timestamp>=? and timestamp<=?`

const sqlite_getTemperatureSensors = `SELECT DISTINCT sensor FROM temperatures WHERE user=? and timestamp>=? and timestamp<=?`

// location functions

const sqlite_getCoordinates = `SELECT l_id, lat, long FROM location`
const sqlite_addLocation = `INSERT INTO location VALUES (?, ?, ?, ?)`
