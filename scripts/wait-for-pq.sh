#!/bin/sh

RETRIES=5

until psql -h $PG_HOST -U $PG_USER -d $PG_DATABASE -p $PG_PORT -c "select 1" > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
  echo "Waiting for postgres server $PG_USER @ $PG_HOST:$PG_PORT $PG_DATABASE"
  RETRIES=$((RETRIES-=1))
  sleep 1
done
