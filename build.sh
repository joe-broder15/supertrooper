#!/bin/bash
echo "building super-c2"
go build -o super-c2.exe cmd/super-c2/main.go 

echo "building super-agent"
go build -o super-agent.exe cmd/super-agent/main.go 
