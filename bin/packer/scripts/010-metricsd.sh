#!/bin/bash
set -e

echo Install metricsd agent -------------------------------
mkdir /tmp/metricsd
cd /tmp/metricsd
sudo mv /home/centos/etc/metricsd.conf /etc/metricsd.conf
sudo mv /home/centos/etc/systemd/system/metricsd.service /etc/systemd/system/metricsd.service
sudo chmod 664 /etc/systemd/system/metricsd.service

sudo systemctl enable metricsd.service
sudo rm -rf /tmp/metricsd
echo DONE installing metricsd agent -------------------------------
