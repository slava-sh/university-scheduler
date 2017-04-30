#!/bin/bash -e

./build.sh
/usr/bin/time ./out/scheduler <"$1"
