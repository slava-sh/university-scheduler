#!/bin/bash -e

mkdir -p ./out
clang++ -Wall -std=c++14 -O3 -o ./out/main ./src/main.cpp
/usr/bin/time ./out/main <"$1"
