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

package backend

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

var infuraConnection *rpc.Client = nil

func init() {
	conn, err := rpc.Dial("https://mainnet.infura.io/v3/YOUR-API-KEY")
	if err != nil {
		fmt.Println("Could not connect to Infura: %v", err)
		return
	}
	infuraConnection = conn
}

func Dial(rawJSON interface{}, methodName string, arg ...interface{}) error {
	if (infuraConnection == nil) {
		return errors.New("Endpoint error")
	}

	err := infuraConnection.Call(&rawJSON, methodName, arg...)

	if err != nil {
		fmt.Println("Cannot call the procedure:", err)
		return err
	}
	return nil
}
