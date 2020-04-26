#!/bin/sh

cd /go/src/app || exit 1
apt-get update && apt-get install uuid-runtime
make test
