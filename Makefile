export GO111MODULE=on
ABS_PATH=`pwd`
N?=1
TEST_DB_CONTAINER=test_db
DB_IMAGE=postgres:11.5-alpine
DB_HOST=localhost
DB_USER=postgres
NETWORK=mail_me_all_mailmeapp

ifneq (, $(findstring test, $(MAKECMDGOALS)))
	DB_NAME=mailmeapp_test
	DB_PORT=5439
else
	DB_NAME=mailmeapp
	DB_PORT=5436
endif

POSTGRES_URL:=postgres://$(DB_USER)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
POSTGRES_URL_INTERNAL:=$(DB_USER)://postgres@$(TEST_DB_CONTAINER):5432/$(DB_NAME)?sslmode=disable

export PG_HOST=$(DB_HOST)
export PG_PORT=$(DB_PORT)
export PG_DATABASE=$(DB_NAME)
export PG_USER=$(DB_USER)

.PHONY: all restart build-backend restart-backend up-backend restart-all migrate-up migrate-down proto test-backend

all: build-backend up-backend 
restart: restart-all
build-backend:
	$(info Running target $(MAKECMDGOALS))
	docker-compose -f docker-compose.dev.yaml build twproxy 
	docker-compose -f docker-compose.dev.yaml build backend
restart-backend:
	$(info Running target $(MAKECMDGOALS))
	docker-compose -f docker-compose.dev.yaml restart twproxy
	docker-compose -f docker-compose.dev.yaml restart backend
up-backend:
	docker-compose -f docker-compose.dev.yaml up -d backend
restart-all:
	$(info Running target $(MAKECMDGOALS))
	docker-compose -f docker-compose.dev.yaml up -d
	docker-compose -f docker-compose.dev.yaml restart
create-migration:
	$(info Running target $(MAKECMDGOALS))
	docker run --rm -v $(ABS_PATH)/backend/migrations:/migrations migrate create -ext sql -seq -dir /migrations $(NAME) 
migrate-up:
	$(info Running target $(MAKECMDGOALS) with $(POSTGRES_URL))
	docker-compose -f docker-compose.dev.yaml up -d postgresql
	./scripts/wait-for-pq.sh
	docker run --rm --net=host -v $(ABS_PATH)/backend/migrations:/migrations migrate -database $(POSTGRES_URL) -path /migrations up $(N)
migrate-down:
	$(info Running target $(MAKECMDGOALS) with $(POSTGRES_URL))
	docker-compose -f docker-compose.dev.yaml up -d postgresql
	./scripts/wait-for-pq.sh
	docker run --rm --net=host -v $(ABS_PATH)/backend/migrations:/migrations migrate -database $(POSTGRES_URL) -path /migrations down $(N)
test-backend:
	$(info Running target $(MAKECMDGOALS) with $(POSTGRES_URL))
	TEST_DB_CONTAINER=$(TEST_DB_CONTAINER) NETWORK=$(NETWORK) DB_NAME=$(DB_NAME) DB_PORT=$(DB_PORT) DB_IMAGE=$(DB_IMAGE) ABS_PATH=$(ABS_PATH) POSTGRES_URL_INTERNAL=$(POSTGRES_URL_INTERNAL) POSTGRES_URL=$(POSTGRES_URL) ./scripts/run-test.sh
proto: 
	$(info Running target $(MAKECMDGOALS))
        cd backend && protoc -I rpc rpc/twproxy.proto --go_out=plugins=grpc:rpc
