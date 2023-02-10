jsonRPC to REST Server
==================

This project is a relay, written in Go, that allows to convert the
Ethereum (and probably clones) json-rpc API to REST endpoints.

It is interesting because the logic is built around the OpenRPC specification.

It is quite unfinished right now and it has been built to connect to the infura API,
even though the same code will work with any Ethereum node as backend.

Prerequisites
--------------------

This application needs docker compose to be built.

It is possible to prepare and run the docker image running:

`docker-compose up --build`

At this point the REST server should be available on localhost.

It is also possible to run the program natively but it's not suggested as docker create a clean environment and takes care of dependencies:

`go build`

The application is designed so that it is possible to include it in a microservice infrastructure.

The REST server
--------------------

This REST server uses the OpenRPC specification to build a service that can be used to relay jsonRPC to a REST API.

The application is able to parse the Ethereum OpenRPC specification and dynamically build the required handlers to serve an API as similar as possible to the jsonRPC counterpart.

The general structure is as follows:

`host/api/procedureName?arg1=...&arg2=...`

For example:

```
http://localhost:8080/api/eth_getTransactionByBlockNumberAndIndex?blockNumber=0xD59F80&index=0x10

http://localhost:8080/api/eth_getBlockByNumber?blockNumber=latest&includeTransactions=true

http://localhost:8080/api/eth_getBlockByNumber?blockNumber=0xD59F80&includeTransactions=false
```

The procedures available are:

```
web3_clientVersion

web3_sha3

net_listening

net_peerCount

net_version

eth_blockNumber

eth_call

eth_chainId

eth_coinbase

eth_estimateGas

eth_gasPrice

eth_getBalance

eth_getBlockByHash

eth_getBlockByNumber

eth_getBlockTransactionCountByHash

eth_getBlockTransactionCountByNumber

eth_getCode

eth_getFilterChanges

eth_getFilterLogs

eth_getRawTransactionByHash

eth_getRawTransactionByBlockHashAndIndex

eth_getRawTransactionByBlockNumberAndIndex

eth_getLogs

eth_getStorageAt

eth_getTransactionByBlockHashAndIndex

eth_getTransactionByBlockNumberAndIndex

eth_getTransactionByHash

eth_getTransactionCount

eth_getTransactionReceipt

eth_getUncleByBlockHashAndIndex

eth_getUncleByBlockNumberAndIndex

eth_getUncleCountByBlockHash

eth_getUncleCountByBlockNumber

eth_getProof

eth_getWork

eth_hashrate

eth_mining

eth_newBlockFilter

eth_newFilter

eth_newPendingTransactionFilter

eth_pendingTransactions

eth_protocolVersion

eth_sendRawTransaction

eth_submitHashrate

eth_submitWork

eth_syncing

eth_uninstallFilter
```

However since this is a Proof of Concept many procedures haven't been tested.

See `data/targets.txt` for some procedures known to work.

Load testing
--------------------

I have implemented a load testing script that uses `vegeta` to generate
nice reports. The script will create a directory called `loadtest`
where the tool will be downloaded and the reports stored.

After running the `loadtest.sh` file, html report files will be generated
in the loadtest directory. 

Load testing has been done using two main sample cases:

1) eth_getTransactionByBlockNumberAndIndex
2) eth_getBlockByNumber

The main parameters are, 30 seconds of Duration and 5 messages sent per second in each burst.

Using a representative run the functions are able to serve respectively:

1) eth_getTransactionByBlockNumberAndIndex

```
Requests      [total, rate, throughput]         150, 5.03, 5.00

Duration      [total, attack, wait]             29.999s, 29.8s, 199.018ms

Latencies     [min, mean, 50, 90, 95, 99, max]  133.033ms, 177.521ms, 169.534ms, 209.324ms, 213.54ms, 416.608ms, 606.849ms

Bytes In      [total, mean]                     89700, 598.00

Bytes Out     [total, mean]                     0, 0.00

Success       [ratio]                           100.00%

Status Codes  [code:count]                      200:150  
```

By looking at the latencies the max response time has been under 600ms and the best response has been 133 ms, this
indicates that the server is able to easily serve 30 requests per second.

2) eth_getBlockByNumber

```
Requests      [total, rate, throughput]         150, 5.03, 5.01

Duration      [total, attack, wait]             29.962s, 29.8s, 161.824ms

Latencies     [min, mean, 50, 90, 95, 99, max]  161.52ms, 262.926ms, 259.09ms, 275.907ms, 362.744ms, 414.132ms, 521.087ms

Bytes In      [total, mean]                     25150494, 167669.96

Bytes Out     [total, mean]                     0, 0.00

Success       [ratio]                           100.00%

Status Codes  [code:count]                      200:150
```

By looking at the latencies the max response time has been under 700ms and the best response has been 133 ms, this
indicates that the server has been able to easily serve 30 requests per second.

Future improvements
--------------------

* Extend the support to parts of the specification not covered in this PoC
* Generate documentation from the specs to swagger
* Improve data validation
* Support native ethereum nodes
* Improve testing/benchmarking, automatic generation of benchmarks and tests for each function
