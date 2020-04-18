#!/bin/bash
for secret in ./prod_secrets/*; do
    secret_name=$(basename $secret);
    cat $secret | docker secret create $secret_name -
done