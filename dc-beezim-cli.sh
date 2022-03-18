#!/bin/sh

docker-compose run --rm \
  --user $(id -u):$(id -g) \
  beezim ./bin/beezim-cli ${@:1}
