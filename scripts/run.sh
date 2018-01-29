#!/usr/bin/env bash

#start udnr service...

function help()
{
        echo " ./run.sh env"
}

echo $(pwd)

if [ $# -ne 1 ];then
        help
        exit 1
fi

env=$1

echo $env

cd /root/account-auth-service

if [ $env = "test" ];then
    cp ./config/config_test.json ./config/config.json
elif [ $env = "online" ];then
    cp ./config/config_online.json ./config/config.json
else
    help
    exit 1
fi

nohup ./account_auth -c ./config/config.json



