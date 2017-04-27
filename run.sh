#!/bin/bash -e

mkdir -p out
(cd src && go build -o ../out/main)
/usr/bin/time ./out/main <"$1"
