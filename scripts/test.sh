#!/usr/bin/env bash

export DOCKER_IMAGE_TAG=corona-crawler/crawler:dev
export COMPOSE_FILE=docker-compose.test.yml
export COMPOSE_PROJECT_NAME=corona_crawler_dev_test

function cleanup() {
    docker-compose down --volumes
}
trap cleanup EXIT

docker-compose up -d postgres
docker-compose run test
