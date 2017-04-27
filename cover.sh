#!/bin/bash -e

mkdir -p out
go test -v -covermode=count -coverprofile=./out/coverage.out ./src
go tool cover -html=./out/coverage.out
