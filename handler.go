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

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"ethereum-jsonrpc-to-rest/backend"
	"ethereum-jsonrpc-to-rest/logger"
	"ethereum-jsonrpc-to-rest/responsebuilder"
	"ethereum-jsonrpc-to-rest/rpcparser"

	"github.com/gorilla/mux"
)

var listeners = make(map[string]rpcparser.RPCMethod)

func AddListener(method rpcparser.RPCMethod) {
	listeners[method.Name] = method
}

func handleRequest(name string, w http.ResponseWriter, r *http.Request) {
	// Check if the procedure exists
	method, exists := listeners[name]
	if !exists {
		responsebuilder.ReplyError(http.StatusBadRequest, w, r, "Bad Request")
		return
	}

	params := make([]interface{}, 0)
	for _, param := range method.Params {
		if param.Name == "" {
			// TODO: there are a couple of these to handle in the schema
			// we should resolve it to get the properties
			if param.Ref == "#/components/contentDescriptors/BlockNumber" {
				param.Name = "blockNumber"
			} else if param.Ref == "#/components/contentDescriptors/Transaction" {
				param.Name = "transaction"
			} else {
				continue
			}
		}

		urlParam := r.URL.Query().Get(param.Name)

		// TODO: we should use the params Schema type
		// to check the parameters data using some custom regex
		// and/or using the regex provided by the Param.Pattern

		if urlParam == "" && param.Required == true {
			responsebuilder.ReplyError(http.StatusBadRequest, w, r,
				"Missing required parameter")
			return
		}

		if param.Schema.Type == "boolean" {
			// TODO: we have to handle other types including param.Schema.Ref
			boolValue, err := strconv.ParseBool(urlParam)
			if err != nil {
				responsebuilder.ReplyError(http.StatusBadRequest, w, r,
					"Bad Request: "+err.Error())
				return
			}
			params = append(params, boolValue)
		} else {
			params = append(params, urlParam)
		}
	}

	var rawJSON interface{}
	err := backend.Dial(&rawJSON, name, params...)
	if err != nil {
		responsebuilder.ReplyError(http.StatusInternalServerError, w, r,
			"Bad Request: "+err.Error())
		return
	}

	responsebuilder.ReplyOK(w, r, rawJSON)
}

type methodHandler struct {
	h    http.Handler
	opts rpcparser.RPCMethod
}

func (m methodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleRequest(m.opts.Name, w, r)
}

func MethodHandler(o rpcparser.RPCMethod) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return methodHandler{h, o}
	}
}

func RunServer() {
	// Read the ethereum rpc spec file
	content, err := ioutil.ReadFile("./data/openrpc.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload rpcparser.RPCDef

	// This will create our parsed structure
	err = rpcparser.Parse(content, &payload)
	if err != nil {
		log.Fatal("Error when parsing file: ", err)
	}

	// We can now register the handlers for the automagically generated API
	router := mux.NewRouter().StrictSlash(true)
	for _, method := range payload.Methods {
		opts := method
		handler := MethodHandler(opts)

		AddListener(method)

		var h http.Handler
		router.
			// TODO set the http verb
			Methods("GET").
			Path("/api/" + method.Name).
			Name(method.Name).
			Handler(handler(logger.Logger(h, method.Name)))
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	RunServer()
}
