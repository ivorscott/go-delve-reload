#!/bin/bash
set -x

# delete server
docker-machine rm -y gdr${server} &

# delete all storage in DO (be sure you are ok deleting ALL storage in an account)
# doctl compute volume ls --format ID --no-header | while read -r id; do doctl compute volume rm -f "$id"; done