LOCATIONS_DB_URL := postgres://root:shadow@localhost:5432/locations?sslmode=disable
HISTORY_DB_URL := postgres://root:shadow@localhost:5432/history?sslmode=disable

migrate_users_up:
	migrate -database ${LOCATIONS_DB_URL} -path db/migrations/locations up

migrate_history_up:
	migrate -database ${HISTORY_DB_URL} -path db/migrations/history up

migrate_users_down:
	migrate -database ${LOCATIONS_DB_URL} -path db/migrations/locations down

migrate_history_down:
	migrate -database ${HISTORY_DB_URL} -path db/migrations/history down

protoc_gen:
	sh ./scripts/protoc-gen.sh

run_locations:
	go run ./cmd/locations

.PHONY: migrate_users_up migrate_users_down migrate_history_up migrate_history_down protoc_gen run_locations
