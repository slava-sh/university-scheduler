#!/bin/bash -e

if [[ $(git status --porcelain) ]]; then
  echo 'Please commit all changes'
  exit 1
fi

PACKAGE='github.com/slava-sh/university-scheduler'

mkdir -p ./out
{
  echo "// https://${PACKAGE}"
  echo '//'
  git log --oneline --max-count=5 | awk '{ print "// " $0 }'
  echo '// ...'
  echo
  cat ./src/main.cpp
} >./out/submission.cpp
