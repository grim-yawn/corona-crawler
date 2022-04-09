#!/usr/bin/env bash

export DOCKER_IMAGE_TAG=corona-crawler/crawler:dev
export COMPOSE_PROJECT_NAME=corona_crawler_dev

docker-compose up -d
