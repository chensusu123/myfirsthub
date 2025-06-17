#!/bin/bash
mode=dev HOSTNAME=maze-main-server-c0-g5-0 S_GROUP_ID=5 MAZE_REDIS_ADDR_5=10.101.110.231:9004 LOCAL_DEV=true go run main.go
