#!/bin/bash -e

./build.sh
go build -o ./out/checker ./checker.go
for i in $(seq -w 01 10); do
  echo "input/$i.txt --------------------------------------------------"
  ./out/scheduler <"./input/$i.txt" >"./out/$i.txt"
  cat "./input/$i.txt" "./out/$i.txt" | ./out/checker
done
