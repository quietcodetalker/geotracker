LOCATIONS_DB_URL := postgres://root:shadow@localhost:5432/locations?sslmode=disable
HISTORY_DB_URL := postgres://root:shadow@localhost:5432/history?sslmode=disable

migrate_locations_up:
	migrate -database ${LOCATIONS_DB_URL} -path db/migrations/locations up

migrate_history_up:
	migrate -database ${HISTORY_DB_URL} -path db/migrations/history up

migrate_locations_down:
	migrate -database ${LOCATIONS_DB_URL} -path db/migrations/locations down

migrate_history_down:
	migrate -database ${HISTORY_DB_URL} -path db/migrations/history down

protoc_gen:
	sh ./scripts/protoc-gen.sh

run_locations:
	go run ./cmd/locations

run_history:
	go run ./cmd/history

build_history_image:
	DOCKER_BUILDKIT=0 docker build -t registry.gitlab.com/spacewalker/locations/history:latest --tag history:latest -f ./deployments/history/Dockerfile .

build_locations_image:
	DOCKER_BUILDKIT=0 docker build -t registry.gitlab.com/spacewalker/locations/locations:latest --tag locations:latest -f ./deployments/locations/Dockerfile .

.PHONY: migrate_locations_up \
		migrate_locations_down \
		migrate_history_up \
		migrate_history_down \
		protoc_gen \
		run_locations \
		run_history \
		build_history_image \
		build_locations_image
