ABS_PATH=`pwd`

all: build-backend
restart: build-backend restart-backend
build-backend:
	docker-compose -f docker-compose.dev.yaml build backend
restart-backend:
	docker-compose -f docker-compose.dev.yaml up -d backend
create-migration:
	cd backend
	docker run -v $(ABS_PATH)/backend/migrations:/migrations migrate create -ext sql -seq -dir /migrations $(NAME) 

