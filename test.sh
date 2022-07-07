#!/bin/bash

docker compose up -d --build
USERS_SERVICE_BASE_URL=http://localhost:8081 go test -v ./... && docker compose down
