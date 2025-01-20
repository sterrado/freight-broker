.PHONY: build up down logs test clean

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

reset:
	docker-compose down -v
	docker-compose build
	docker-compose up -d