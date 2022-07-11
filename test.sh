#!/bin/bash

./start.sh "$1"
go clean -testcache
BASE_URL=http://localhost:8080 API_BASE_URL="$1" go test -v ./... && API_BASE_URL="$1" docker compose down
