#!/bin/bash

USERS_SERVICE_IMAGE_URL="$1" docker compose up -d --build
