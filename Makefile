.PHONY: up down restart logs test

# Start all services defined in docker-compose
up:
	docker compose up -d --build

# Stop and remove containers, networks, volumes, and images created by up
down:
	docker compose down

# Restart services
restart: down up

# Follow logs for all services
logs:
	docker compose logs -f

# Run all Go tests in the 'test' directory and its subdirectories
test:
	go test -v -race ./test/...