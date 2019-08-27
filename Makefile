ABS_PATH=`pwd`
POSTGRES_URL=postgres://postgres@localhost:5436/mailmeapp?sslmode=disable
N?=1

all: build-backend
restart: build-backend restart-backend
build-backend:
	docker-compose -f docker-compose.dev.yaml build backend
restart-backend:
	docker-compose -f docker-compose.dev.yaml up -d backend
create-migration:
	cd backend
	docker run -v $(ABS_PATH)/backend/migrations:/migrations migrate create -ext sql -seq -dir /migrations $(NAME) 
migrate-up:
	cd backend
	docker run --net=host -v $(ABS_PATH)/backend/migrations:/migrations migrate -database $(POSTGRES_URL) -path /migrations up $(N)
migrate-down:
	cd backend
	docker run --net=host -v $(ABS_PATH)/backend/migrations:/migrations migrate -database $(POSTGRES_URL) -path /migrations down $(N)
