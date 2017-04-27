#!/bin/bash -e

if [[ $(git status --porcelain) ]]; then
  echo 'Please commit all changes'
  exit 1
fi

COMMIT=$(git rev-parse --short HEAD)

mkdir -p ./out
{
  echo "// https://github.com/slava-sh/university-scheduler"
  echo "// commit ${COMMIT}"
  (cd ./src && bundle -prefix ' ' .) | tail -n +3
} >./out/submission.go
