#!/bin/sh
go build
./ethereum-jsonrpc-to-rest &

mkdir -p loadtest
cd loadtest

if [[ ! -f vegeta ]]; then
wget -N https://github.com/tsenart/vegeta/releases/download/v12.8.3/vegeta-12.8.3-linux-amd64.tar.gz
tar -zxvf vegeta-12.8.3-linux-amd64.tar.gz
chmod +x vegeta
fi

echo "GET http://localhost:8080/api/eth_getTransactionByBlockNumberAndIndex?blockNumber=0xD59F80&index=0x10" | ./vegeta attack -duration=30s -rate=5 -output eth_getTransactionByBlockNumberAndIndex.bin
./vegeta report -type=text eth_getTransactionByBlockNumberAndIndex.bin
./vegeta plot -title=eth_getTransactionByBlockNumberAndIndex%20Results eth_getTransactionByBlockNumberAndIndex.bin > eth_getTransactionByBlockNumberAndIndex_results.html

echo "GET http://localhost:8080/api/eth_getBlockByNumber?blockNumber=latest&includeTransactions=true" | ./vegeta attack -duration=30s -rate=5 -output eth_getBlockByNumber.bin
./vegeta report -type=text eth_getBlockByNumber.bin
./vegeta plot -title=eth_getBlockByNumber%20Results eth_getBlockByNumber.bin > eth_getBlockByNumber_results.html

./vegeta attack --targets=../data/targets.txt -duration=5s -rate=5 -output eth_gasPrice.bin
./vegeta report -type=text eth_gasPrice.bin
./vegeta plot -title=eth_gasPrice%20Results eth_gasPrice.bin > eth_gasPrice_results.html

pkill ethereum-jsonrpc-to-rest
