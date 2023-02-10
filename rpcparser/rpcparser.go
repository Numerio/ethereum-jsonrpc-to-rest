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

package rpcparser

import (
	"errors"
	"encoding/json"
	"log"
)

type RPCParamSchema struct {
	Title   string
	Type    string
	Pattern string

	Ref string `json:"$ref"`
}

type RPCParam struct {
	Name        string
	Required    bool
	Description string
	Summary     string

	Schema      RPCParamSchema

	Ref string `json:"$ref"`
}

type RPCMethod struct {
	Name        string
	Description string
	Summary     string

	Params []RPCParam
}

type RPCDef struct {
	OpenRPC string

	Methods []RPCMethod
}

func Parse(content []byte, output *RPCDef) error {
	err := json.Unmarshal(content, output)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return err
	}

	if output.OpenRPC != "1.0.0" {
		return errors.New("Unsupported OpenRPC version")
	}

	return nil
}
