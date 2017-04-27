#!/bin/bash -e

mkdir -p ./out
go test -v -cpuprofile ./out/cpu.prof -bench . -benchtime 3s -run '^$' ./src
mv ./src.test ./out/main.bench
(cd ./out && go tool pprof ./main.bench ./cpu.prof)
