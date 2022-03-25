# Kubernetes Hyperledger Fabric with REST API

## Prerequisites 

- [Docker](https://www.docker.com)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [jq](https://stedolan.github.io/jq/)

## Quickstart

Create a local Kubernetes cluster:
```shell
./network kind
```

Launch the network, create a channel, and deploy the [smart contract](./chaincode/fabric-ccaas-asset-transfer-basic/README.md) : 
```shell
./network up
./network channel create
./network chaincode deploy
```

Invoke and query chaincode:
```shell
./network chaincode invoke '{"Args":["InitLedger"]}' 
./network chaincode query '{"Args":["GetAllAssets"]}'
```

Access the blockchain with a [REST API](./rest-api-typescript/README.md): 
```
./network rest-easy
```

Tear down the test network: 
```shell
./network down 
```

Tear down the cluster: 
```shell
./network unkind
```