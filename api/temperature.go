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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"lachut.net/gogs/dslachut/go-irleak/kb"
)

type temperaturePostBody struct {
	Token     string  `json:"token"`
	Sensor    string  `json:"sensor"`
	Value     float64 `json:"value"`
	Timestamp float64 `json:"timestamp"`
}

func TemperatureHandler(w http.ResponseWriter, r *http.Request, k kb.KB) {
	if r.Method == "POST" {
		temperaturePost(w, r, k)
		//} else if r.Method == "GET" {
		//	temperatureGet(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func temperaturePost(w http.ResponseWriter, r *http.Request, k kb.KB) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		requestFailed(w, http.StatusNoContent)
		log.Printf("empty or errant request body\n")
		return
	}

	rec := temperaturePostBody{}
	if json.Unmarshal(body, &rec) != nil {
		requestFailed(w, http.StatusInternalServerError)
		log.Printf("request not in json format\n")
		return
	}

	user, newToken, ok := checkToken(rec.Token, k)
	if !ok {
		requestFailed(w, http.StatusForbidden)
		log.Printf("invalid token")
		return
	}

	success := apiResponse{true, newToken}
	ok = k.AddTemperature(user, rec.Sensor, rec.Timestamp, rec.Value)
	if !ok {
		success = apiResponse{false, newToken}
	}

	payload, _ := json.Marshal(success)
	w.Write(payload)
}
