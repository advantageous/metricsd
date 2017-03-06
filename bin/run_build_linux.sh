#!/usr/bin/env bash
source ~/.bash_profile
docker pull cloudrable/metricsd:latest
docker run -it --name build \
-v `pwd`:/gopath/src/github.com/cloudrable/metricsd \
cloudrable/metricsd \
/bin/sh -c "/gopath/src/github.com/cloudrable/metricsd/bin/build_linux.sh"
docker rm build


