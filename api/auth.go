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
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/elithrar/simple-scrypt"
	"github.com/spf13/viper"
	"lachut.net/gogs/dslachut/go-irleak/kb"
)

type authPostBody struct {
	User string `json:"user"`
	Pass string `json:"password"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request, k kb.KB) {
	if r.Method == "POST" {
		authPost(w, r, k)
		//} else if r.Method == "GET" {
		//	authGet(w, r)
	} else {
		requestFailed(w, http.StatusBadRequest)
		log.Printf("bad request method %s\n", r.Method)
	}
}

func authPost(w http.ResponseWriter, r *http.Request, k kb.KB) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		requestFailed(w, http.StatusNoContent)
		log.Printf("empty or errant request body\n")
		return
	}

	rec := authPostBody{}
	if json.Unmarshal(body, &rec) != nil {
		requestFailed(w, http.StatusInternalServerError)
		log.Printf("request not in json format\n")
		return
	}

	if rec.Pass == "" {
		requestFailed(w, http.StatusBadRequest)
		log.Printf("Missing password\n")
		return
	}

	hash, ok := k.GetHash(rec.User)
	if !ok {
		requestFailed(w, http.StatusNotFound)
		log.Printf("user '%s' not found\n", rec.User)
		return
	}

	err = scrypt.CompareHashAndPassword(hash, []byte(rec.Pass))
	if err != nil {
		requestFailed(w, http.StatusForbidden)
		log.Printf("bad password\n")
		return
	}

	token, err := generateToken(rec.User, k)
	if err != nil || token == "" {
		requestFailed(w, http.StatusInternalServerError)
		log.Printf("token generating error\n")
		return
	}

	success := apiResponse{true, token}
	payload, _ := json.Marshal(success)
	w.Write(payload)
}

func generateToken(user string, k kb.KB) (string, error) {
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("%x", tokenBytes)
	exp := time.Now().Unix() + viper.GetInt64("exptoken")
	k.AddToken(user, token, exp)
	return token, nil
}

func checkToken(token string, k kb.KB) (user string, newToken string, ok bool) {
	now := time.Now().Unix()
	user, exp, ok := k.GetUser(token)
	if !ok || exp < now || user == "" {
		log.Printf("%v.%v.%v.%v\n", user, exp, ok, token)
		return "", "", false
	}

	newToken, err := generateToken(user, k)
	if err != nil {
		log.Println(err)
		return "", "", ok
	}

	k.ExpireToken(token)
	return
}

func PurgeTokens(k kb.KB, done chan bool) {
	tick := time.NewTicker(time.Second * time.Duration(viper.GetInt64("exptoken")))
	for {
		select {
		case <-tick.C:
			now := time.Now().Unix()
			k.PurgeTokens(now)
		case <-done:
			return
		}
	}
}
