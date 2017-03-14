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
	"net/http"
)

type apiResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func requestFailed(w http.ResponseWriter, status int) {
	failure := apiResponse{false, ""}
	payload, _ := json.Marshal(failure)
	w.WriteHeader(status)
	w.Write(payload)
}
