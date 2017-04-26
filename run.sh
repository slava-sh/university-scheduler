#!/bin/bash -e

mkdir -p out
(cd src && go build -o ../out/main)
./out/main <"$1"
