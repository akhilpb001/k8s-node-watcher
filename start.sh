#!/bin/bash

IMAGE_REPOSITORY=${1:-docker-sandbox.infra.cloudera.com/apb/k8s-node-watcher}
IMAGE_VERSION=${2:-1.0.0-SNAPSHOT}

docker run --name health-monitor -it --rm ${IMAGE_REPOSITORY}:${IMAGE_VERSION}