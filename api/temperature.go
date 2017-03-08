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

package api

import (
	"io"
	"net/http"
)

func TemperatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		temperaturePost(w, r)
		//} else if r.Method == "GET" {
		//	temperatureGet(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func temperaturePost(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}
