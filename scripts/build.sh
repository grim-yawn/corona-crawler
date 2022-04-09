#!/usr/bin/env bash

export DOCKER_IMAGE_TAG=corona-crawler/crawler:dev

# I would use proper tag here with version or commit hash inside CI
# This is a main reason why this script exists to provide consistent docker image for all steps
docker build -t ${DOCKER_IMAGE_TAG} .
