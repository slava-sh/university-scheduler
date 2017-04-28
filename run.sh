#!/bin/bash -e

if [[ $(git status --porcelain) ]]; then
  echo 'Please commit all changes'
  echo
fi

COMMIT=$(git rev-parse --short HEAD)

mkdir -p ./out ./reports
(cd ./src && go build -o ../out/main)
/usr/bin/time ./out/main | tee "./reports/${COMMIT}.txt"
