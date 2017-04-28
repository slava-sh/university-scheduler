#!/bin/bash -e

if [[ $(git status --porcelain) ]]; then
  echo 'Please commit all changes'
  echo
fi

PACKAGE='github.com/slava-sh/university-scheduler'
COMMIT=$(git rev-parse --short HEAD)

mkdir -p ./out
{
  echo "// https://${PACKAGE}"
  echo "// commit ${COMMIT}"
  (cd ./src && bundle -prefix ' ' "${PACKAGE}/src") \
    | tail -n +3 \
    | perl -pe 's/^(\t+local += +)true$/$1false/'
} >./out/submission.go
