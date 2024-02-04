#!/bin/bash

# No knowledge of migrations for Go, simple script to update the database schema
# Should not be used in production, only for development

docker cp ./db/schema.sql services-db-1:/schema.sql
docker exec -it services-db-1 psql -U postgres -d tasks -f /schema.sql
