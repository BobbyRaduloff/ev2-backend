#!/bin/bash

POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=verifier

# Create a directory on your host to store PostgreSQL data
mkdir -p $HOME/docker/volumes/postgres

# Run the PostgreSQL container with a volume mount
docker run --name verifier-postgres \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_DB=$POSTGRES_DB \
    -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data \
    -p 5432:5432 \
    --restart unless-stopped \
    -d postgres
