#!/usr/bin/env bash


rm metricsd_linux

set -e


cd /gopath/src/github.com/cloudrable/metricsd/

echo "Running go clean"
go clean
echo "Running go get"
go get
echo "Running go build"
go build
echo "Renaming output to _linux"
mv metricsd metricsd_linux
sudo cp  metricsd_linux /usr/bin/metricsd
