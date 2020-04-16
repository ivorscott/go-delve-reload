#!/bin/bash
set -x

# enable monitoring
for server in 1; do
  docker-machine scp daemon.json gdr${server}:/etc/docker/ && 
  docker-machine ssh gdr${server} systemctl restart docker &
done