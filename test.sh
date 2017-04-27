#!/bin/bash -e

if type godebug >/dev/null; then
  # Use godebug if available.
  GO=godebug
else
  GO=go
fi

$GO test ./src
