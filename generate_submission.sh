#!/bin/bash -e

if [[ $(git status --porcelain) ]]; then
  echo 'Please commit all changes'
  exit 1
fi

PACKAGE='github.com/slava-sh/university-scheduler'

mkdir -p ./out
{
  {
    echo "https://${PACKAGE}"
    echo
    git log --oneline --max-count=5
    echo '...'
  } | awk '{ print "// " $0 }'
  (cd ./src && bundle -prefix ' ' "${PACKAGE}/src") \
    | tail -n +3 \
    | perl -pe 's/^(\t+local += +)true$/$1false/'
} >./out/submission.go
