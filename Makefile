.PHONY: test-db-unit test-db-integration test-db-all up down clean

# Запуск юнит-тестов (требует локальной БД)
test-db-unit:
	go test -v ./internal/db/... -tags=unit

# Запуск интеграционных тестов в Docker
test-db-integration:
	docker-compose up --build --exit-code-from dbtest

# Запуск всех тестов
test-db-all: test-db-unit test-db-integration

# Запуск контейнеров для ручного тестирования
up:
	docker-compose up -d

# Остановка контейнеров
down:
	docker-compose down

# Очистка
clean: down
	docker volume rm order-service_pgdata