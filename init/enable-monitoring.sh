#!/bin/bash
set -x

# enable monitoring
docker-machine scp daemon.json gdr1:/etc/docker/ && 
docker-machine ssh gdr1 systemctl restart docker &