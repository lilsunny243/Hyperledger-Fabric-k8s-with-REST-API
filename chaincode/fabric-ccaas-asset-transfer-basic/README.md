## Build a Chaincode Docker Image

Before chaincode can be started in the network, it must be compiled, linked with the grpc `ChaincodeServer`,
embedded into a Docker image, and pushed to a container registry visible to the Kubernetes cluster.

By default, the `./network` script will launch the [asset-transfer-basic](../../asset-transfer-basic/chaincode-external)
chaincode.  When the test network installs this chaincode, there is no need to build a custom Docker image as it
has previously been uploaded to a public container registry.

As an exercise, we recommend making some updates to the asset transfer basic chaincode and then running the
modified smart contract on your local Kubernetes cluster.  For instance, the current version of the
[assetTransfer.go](../../asset-transfer-basic/chaincode-external/assetTransfer.go) code is completely
silent, printing nothing to the log when functions are invoked in the container.  Try adding some debugging
information to the stdout of this process, bundling into a Docker image, and pushing the docker
image to the local development container registry.

1. Add some print statements to assetTransfer.go.  E.g.:
```java 
    fmt.Printf("reading asset %s\n", id)
```

2.  Build the docker image locally with:
```shell
docker build -t asset-transfer-basic ../asset-transfer-basic/chaincode-external 
```

3. Override the test network's default chaincode image, pointing to our local container registry: 
```shell
export TEST_NETWORK_CHAINCODE_IMAGE=localhost:5000/asset-transfer-basic
```

3. Publish the custom image to the local registry: 

```shell
docker tag asset-transfer-basic $TEST_NETWORK_CHAINCODE_IMAGE 
docker push $TEST_NETWORK_CHAINCODE_IMAGE
```