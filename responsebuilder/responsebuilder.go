/* 
 * This file is part of the Ethereum-JSONRPC-to-REST distribution
 * Copyright (c) 2022-2023 Dario Casalinuovo.
 * 
 * This program is free software: you can redistribute it and/or modify  
 * it under the terms of the GNU General Public License as published by  
 * the Free Software Foundation, version 2.
 *
 * This program is distributed in the hope that it will be useful, but 
 * WITHOUT ANY WARRANTY; without even the implied warranty of 
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU 
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License 
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package responsebuilder

import (
	"encoding/json"
	"net/http"
)

type responseErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func ReplyOK(w http.ResponseWriter, r *http.Request, rawJSON interface{}) {
	replyJSON(w, http.StatusOK, rawJSON)
}

func ReplyError(err int, w http.ResponseWriter, r *http.Request, text string) {
	replyJSON(w, err, map[string]string{"error": text})
}

func replyJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(response))
}
