#!/bin/bash

USERS_SERVICE_IMAGE="$1" docker compose up -d --build
