#!/bin/bash

./start.sh
go clean -testcache
BASE_URL=http://localhost:8080 USERS_SERVICE_BASE_URL="$1" go test -v ./... && docker compose down
