#!/bin/bash

VERSION="v0.4"
DOCKER_REGISTRY="docker.io"

rm -rf bin

GOOS=darwin  GOARCH=amd64 go build -o bin/macos/ip-resolver
GOOS=linux   GOARCH=amd64 go build -o bin/linux/ip-resolver
GOOS=windows GOARCH=amd64 go build -o bin/windows/ip-resolver.exe

#tar czf bin/ip-resolver-$VERSION-macos.tar.gz --directory=bin/macos ip-resolver
#tar czf bin/ip-resolver-$VERSION-linux.tar.gz --directory=bin/linux ip-resolver
#zip     bin/ip-resolver-$VERSION-windows.zip -j bin/windows/ip-resolver.exe

mkdir bin/docker
cp Dockerfile bin/docker/Dockerfile
cp bin/linux/ip-resolver bin/docker/ip-resolver

docker build --no-cache -t "$DOCKER_REGISTRY/maxmanuylov/ip-resolver:$VERSION" bin/docker
docker tag "$DOCKER_REGISTRY/maxmanuylov/ip-resolver:$VERSION" "$DOCKER_REGISTRY/maxmanuylov/ip-resolver:latest"

#docker push "$DOCKER_REGISTRY/maxmanuylov/ip-resolver:$VERSION"
#docker push "$DOCKER_REGISTRY/maxmanuylov/ip-resolver:latest"
