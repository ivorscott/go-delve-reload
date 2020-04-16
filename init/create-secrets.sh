#!/bin/bash
for secret in ./secrets/*; do
    secret_name=$(basename $secret);
    cat $secret | docker secret create $secret_name -
done