#!/bin/bash
echo "building server"
go build -o server.exe cmd/server/main.go 

echo "building agent"
go build -o agent.exe cmd/agent/main.go 
