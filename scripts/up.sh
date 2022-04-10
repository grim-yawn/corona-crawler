#!/usr/bin/env bash

export DOCKER_IMAGE_TAG=corona-crawler/crawler:dev
export COMPOSE_PROJECT_NAME=corona_crawler_dev

function cleanup() {
  docker-compose down
}

trap cleanup EXIT

docker-compose up -d

docker-compose logs -f server crawler-history crawler-latest
