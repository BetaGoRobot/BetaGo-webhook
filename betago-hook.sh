#!/bin/sh

docker login -u "$DOCKER_USER_NAME_TENCENT" -p "$DOCKER_PASSWORD_TENCENT"

docker stop betago

docker rm betago

docker pull ccr.ccs.tencentyun.com/kevinmatt/betago:latest

docker run -d --network betago --name betago ccr.ccs.tencentyun.com/kevinmatt/betago:latest