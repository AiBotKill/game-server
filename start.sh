#!/bin/bash

docker build -t game-server .
docker stop game-server
docker rm game-server
docker run -d --name game-server --link gnatsd:gnatsd game-server
