#!/bin/bash

mkdir -p bin

echo "打包file_server"
go build -o bin/file_server cmd/file/main.go

echo "打包logic_server"
go build -o bin/logic_server cmd/logic/main.go

echo "打包connect_server"
go build -o bin/connect_server cmd/connect/main.go

cd bin

echo "停止file_server服务"
pkill file_server

echo "停止logic_server服务"
pkill logic_server

echo "停止connect_server服务"
pkill connect_server

sleep 2

rm -rf *.out *.log

echo "启动file_server服务"
nohup ./file_server &

echo "启动logic_server服务"
nohup ./logic_server &

echo "启动connect_server服务"
nohup ./connect_server &

