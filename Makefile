.PHONY: build up down logs test

build:
	docker build -t .

up:
	docker-compose --env-file .env -f docker-compose.yml up -d --build

down:
	docker-compose --env-file .env -f docker-compose.yml down

logs:
	docker-compose -f docker-compose.yml logs -f

test:
	go test ./...