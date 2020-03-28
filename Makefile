gen.proto:
	protoc --go_out=plugins=grpc:pkg/calendar api/api.proto

postgres.run:
	docker run --name calendar-postgres -e POSTGRES_DB=calendar -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:12.2

postgres.start:
	docker start calendar-postgres

postgres.stop:
	docker stop calendar-postgres

migrate:
	docker run -v $(shell pwd)/internal/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgresql://postgres:password@localhost:5432/calendar?sslmode=disable up
