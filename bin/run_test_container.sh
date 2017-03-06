#!/usr/bin/env bash
docker pull cloudrable/metricsd:latest
docker run  -it --name runner2  \
-p 5514:514/udp \
-v `pwd`:/gopath/src/github.com/cloudrable/metricsd \
cloudrable/metricsd
docker rm runner2
