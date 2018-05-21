# Beautiful Things Server

In this repo you'll be able to find all the code for the Beautiful Things server along with a client implementation.

This is in a very early development stage, so everything is in flux and not really considered production.

# Build

## Docker

  docker build -t beautifulthings:latest docker
  docker tag

## Helm

  helm install --name beautifulthings ./kubernetes/
  kubectl describe pods