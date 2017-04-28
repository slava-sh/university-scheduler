#!/bin/bash -e

mkdir -p ./out
clang++ -Wall -std=c++14 -O3 -o ./out/main ./src/main.cpp
go build -o ./out/checker ./src/checker.go
for i in $(seq -w 01 10); do
  echo "input/$i.txt --------------------------------------------------"
  ./out/main <"./input/$i.txt" >"./out/$i.txt"
  cat "./input/$i.txt" "./out/$i.txt" | ./out/checker
done
