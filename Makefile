all: build-backend
restart: build-backend restart-backend
build-backend:
	docker-compose -f docker-compose.dev.yaml build backend
restart-backend:
	docker-compose -f docker-compose.dev.yaml up -d backend

