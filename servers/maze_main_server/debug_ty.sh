#!/bin/bash
mode=dev POLARIS_CONF_ADRR=10.101.110.52:8093,10.101.110.53:8093,10.101.110.54:8093 POLARIS_SERVER_ADRR=10.101.110.52:8091,10.101.110.53:8091,10.101.110.54:8091 HOSTNAME=test-main-server-c0-g4-0 S_GROUP_ID=4 APPNAME=test-main-server SERVER_ID=4 NAMESPACE=adl-ty POD_LOCAL_IP=127.0.0.1 go run main.go conf.d/config.ini logs/test-main-server test-main-server
