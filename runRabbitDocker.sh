#!/bin/bash

docker run --name verifier-rabbit \
  -d \
  --restart always \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=admin \
  -p 5672:5672 -p 15672:15672 rabbitmq
