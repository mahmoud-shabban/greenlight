#!/bin/bash


# I forgot to be merged with the local one
# Create jaeger container for tracing
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest

docker run -d --name db2 \
    -e POSTRGRES_PASSWORD=pass \
    -p 
    postgrsql