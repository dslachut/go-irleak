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
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type authPostBody struct {
	User string `json:"user"`
	Pass string `json:"password"`
}

type authResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		authPost(w, r)
		//} else if r.Method == "GET" {
		//	authGet(w, r)
	} else {
		requestFailed(w, http.StatusBadRequest)
		log.Printf("bad request method %s\n", r.Method)
	}
}

func authPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		requestFailed(w, http.StatusNoContent)
		log.Printf("empty or errant request body\n")
		return
	}

	rec := authPostBody{}
	if json.Unmarshal(body, &rec) != nil {
		requestFailed(w, http.StatusInternalServerError)
		log.Printf("request not in json format")
		return
	}

	if rec.Pass == "" {
		requestFailed(w, http.StatusBadRequest)
		log.Printf("Missing password\n")
		return
	}

	hash, err := scrypt.GenerateFromPassword([]byte(rec.Pass), scrypt.DefaultParams)
	if err != nil {
		requestFailed(w, http.StatusInternalServerError)
		log.Printf("failed to hash password")
		return
	}
	out := rec.User + " " + rec.Pass + " " + string(hash)
	io.WriteString(w, out)
}

func requestFailed(w http.ResponseWriter, status int) {
	failure := authResponse{false, ""}
	payload, _ := json.Marshal(failure)
	w.WriteHeader(status)
	w.Write(payload)
}
