#!/bin/bash

APP="simple-tee-proxy"
SOURCE="cmd/$APP.go"

GOOS=linux GOARCH=arm go build -o "dist/$APP-linux-arm" "$SOURCE"
GOOS=linux GOARCH=amd64 go build -o "dist/$APP-linux-amd64" "$SOURCE" 
GOOS=freebsd GOARCH=amd64 go build -o "dist/$APP-freebsd-amd64" "$SOURCE"
GOOS=windows GOARCH=amd64 go build -o "dist/$APP-windows-amd64.exe" "$SOURCE"