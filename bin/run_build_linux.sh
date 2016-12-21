#!/usr/bin/env bash
source ~/.bash_profile
docker pull advantageous/metricsd:latest
docker run -it --name build \
-v `pwd`:/gopath/src/github.com/advantageous/metricsd \
advantageous/metricsd \
/bin/sh -c "/gopath/src/github.com/advantageous/metricsd/bin/build_linux.sh"
docker rm build


