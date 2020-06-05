#!/bin/bash

# docker.sh is used to avoid attempting to apt install docker before cloud-init has finished provisioning all apt source lists
until [[ -f /var/lib/cloud/instance/boot-finished ]]; do
    sleep 1
done
apt update
apt install -y docker.io
