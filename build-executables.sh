#!/bin/bash

TAG=$1

# Create release directory
mkdir -p release

# Build for Linux amd64
GOOS=linux GOARCH=amd64 go build -o release/bf4db-linux-amd64-$TAG

# Build for Linux arm64
GOOS=linux GOARCH=arm64 go build -o release/bf4db-linux-arm64-$TAG

# Build for Windows amd64
GOOS=windows GOARCH=amd64 go build -o release/bf4db-windows-amd64-$TAG.exe