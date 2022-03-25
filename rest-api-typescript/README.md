# Asset Transfer REST API Sample

### INTERNAL Gateway Client (In Kubernetes)

#### TODO: Deploy

```shell
./network rest-easy
```

#### Local Container Registry

Docker images built locally can be uploaded to the `localhost:5000` container registry for
immediate access within the Kube/KIND cluster.  In addition to providing fast turn-around to/from containers 
running in Kube, the use of a private container registry allows us to quickly iterate on code without uploading 
images to the Internet.  Even when using _private_ container registries, the use of a local server saves valuable 
time when loading images into the kind control plane.

e.g.: 
```shell
docker build -t localhost:5000/my-gateway-app . 
docker push localhost:5000/my-gateway-app 
```

Provided that the `imagePullPolicy` for the client deployment is not set to `IfNotPresent`, killing the current pod 
running the gateway client will force a refresh with the latest image layer available at the local registry. 
