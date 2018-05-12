#!/bin/bash
echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin quay.io
docker push thorfour/stocktopus
