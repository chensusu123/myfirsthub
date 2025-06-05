#!/bin/bash

set -e

echo "[INFO] Creating output directory ../servers/maze_main_server/bin if it doesn't exist..."
mkdir -p ../servers/maze_main_server/bin/

echo "[INFO] Cross-compiling maze-server for Linux amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../servers/maze_main_server/bin/maze-server ../servers/maze_main_server/main.go

echo "[OK] Build succeeded: ./bin/maze-server"

echo "[INFO] Building Docker image..."
docker-compose  -f ../docker/docker-compose.yml up

echo "[INFO] All done. You can now run with: docker-compose up"

