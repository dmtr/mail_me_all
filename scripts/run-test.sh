#!/bin/sh

docker run --rm --name $TEST_DB_CONTAINER --network $NETWORK -d -e POSTGRES_DB=$DB_NAME -p $DB_PORT:5432 $DB_IMAGE
$ABS_PATH/scripts/wait-for-pq.sh
docker run --rm --network $NETWORK -v $ABS_PATH/backend/migrations:/migrations migrate -database $POSTGRES_URL_INTERNAL -path /migrations up
cd $ABS_PATH/backend && MAILME_APP_DSN=$POSTGRES_URL go test -v ./...
retVal=$?
docker stop $TEST_DB_CONTAINER
exit $retVal
