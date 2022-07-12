#!/bin/bash

./start.sh "$1"
go clean -testcache
BASE_URL=http://localhost:8080 USERS_SERVICE_BASE_URL=http://localhost:8081 go test -v ./... && USERS_SERVICE_IMAGE="$1" docker compose down
