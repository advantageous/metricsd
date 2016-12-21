#!/usr/bin/env bash
docker pull advantageous/metricsd:latest
docker run  -it --name runner2  \
-p 5514:514/udp \
-v `pwd`:/gopath/src/github.com/advantageous/metricsd \
advantageous/metricsd
docker rm runner2
