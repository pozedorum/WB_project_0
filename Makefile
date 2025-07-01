.PHONY: test-db-unit test-db-integration test-db-all up down clean


test-db-unit:
	go test -v ./internal/db/... 

up:
	docker-compose up --build
	sleep 10  # Ждем инициализации Kafka
	make create-topic

create-topic:
	docker exec kafka_container kafka-topics --create \
		--bootstrap-server kafka:9092 \
		--replication-factor 1 \
		--partitions 1 \
		--topic orders

list-topics:
	docker exec kafka kafka-topics --list --bootstrap-server localhost:9092

describe-topic:
	docker exec kafka kafka-topics --describe --bootstrap-server localhost:9092 --topic orders