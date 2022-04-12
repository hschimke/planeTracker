#! /bin/bash
# MUST BE RUN FROM ROOT OF REPOSITORY
readonly GIT_VERSION=`git rev-parse HEAD`
docker build -t planes/web-serv -f docker/Dockerfile.api-server --build-arg "GIT_VERSION=${GIT_VERSION}" .