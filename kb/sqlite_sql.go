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

const sqlite_createAuth = `CREATE TABLE IF NOT EXISTS auth (
	user TEXT PRIMARY KEY,
	hash TEXT NOT NULL
)`

const sqlite_createToken = `CREATE TABLE IF NOT EXISTS token (
	user TEXT REFERENCES auth (user),
	token TEXT UNIQUE NOT NULL,
	exp INTEGER NOT NULL
)`

const sqlite_getHash = `SELECT hash FROM auth WHERE user=?`
