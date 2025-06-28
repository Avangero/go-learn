.PHONY: up down build logs

reload:
	docker-compose up --build
up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose build

logs:
	docker-compose logs -f 