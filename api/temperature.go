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
	"math"
	"net/http"

	"lachut.net/gogs/dslachut/go-irleak/kb"
)

type temperaturePostBody struct {
	Token     string      `json:"token"`
	Sensor    string      `json:"sensor"`
	Value     float64     `json:"value"`
	Timestamp float64     `json:"timestamp"`
	Points    []tempPoint `json:"points"`
}

type tempPoint struct {
	Value     float64 `json:"value"`
	Timestamp float64 `json:"timestamp"`
}

func TemperatureHandler(w http.ResponseWriter, r *http.Request, k kb.KB) {
	if r.Method == "POST" {
		temperaturePost(w, r, k)
	} else if r.Method == "GET" {
		temperatureGet(w, r, k)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

type tempResponse struct {
	apiResponse
	Sensors []temps `json:"sensors"`
}
type temps struct {
	Name   string              `json:"name"`
	Values map[float64]float64 `json:"values"`
}

func temperatureGet(w http.ResponseWriter, r *http.Request, k kb.KB) {
	params := r.URL.Query()
	reqToken, ok := params["token"]
	if !ok {
		requestFailed(w, http.StatusForbidden)
	}
	user, newToken, ok := checkToken(reqToken, k)
	if !ok {
		requestFailed(w, http.StatusForbidden)
	}
	start_ts, is_start := params["start"]
	if !is_start {
		start_ts = 0
	}
	end_ts, is_end := params["end"]
	if !is_end {
		end_ts = math.MaxFloat64
	}

	sens := make([]string, 0, 1)
	sensor, is_sensor := params["sensor"]
	if is_sensor {
		sens = append(sens, sensor)
	} else {
		sens = k.GetTemperatureSensors(user, start_ts, end_ts)
	}

	tr := new(tempResponse)
	tr.Token = newToken
	tr.Success = true
	tr.Sensors = make([]temps, 0, len(sens))
	for _, sensor = range sens {
		tr.Sensors = append(tr.Sensors, temps{sensor, k.GetTemperatures(user, sensor, start_ts, end_ts)})
	}
	payload, _ := json.Marshal(tr)
	w.Write(payload)
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
	if len(rec.Points) == 0 {
		ok = k.AddTemperature(user, rec.Sensor, rec.Timestamp, rec.Value)
		if !ok {
			success = apiResponse{false, newToken}
		}
	} else {
		for _, pt := range rec.Points {
			ok = k.AddTemperature(user, rec.Sensor, pt.Timestamp, pt.Value)
			success.Success = ok && success.Success
		}
	}

	payload, _ := json.Marshal(success)
	w.Write(payload)
}
