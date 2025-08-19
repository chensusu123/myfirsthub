#!/bin/bash
mode=dev POLARIS_CONF_ADRR=10.101.110.52:8093,10.101.110.53:8093,10.101.110.54:8093 POLARIS_SERVER_ADRR=10.101.110.52:8091,10.101.110.53:8091,10.101.110.54:8091 HOSTNAME=maze-main-server-c0-g5-0 S_GROUP_ID=5 APPNAME=maze-main-server SERVER_ID=5 NAMESPACE=adl-cs POD_LOCAL_IP=127.0.0.1 go run main.go conf.d/config.ini logs/maze-main-server maze-main-server
