USERS_DB_URL := postgres://root:shadow@localhost:5432/users?sslmode=disable
TRANSITIONS_DB_URL := postgres://root:shadow@localhost:5432/transitions?sslmode=disable

migrate_users_up:
	migrate -database ${USERS_DB_URL} -path db/migrations/users up

migrate_transitions_up:
	migrate -database ${TRANSITIONS_DB_URL} -path db/migrations/transitions up

migrate_users_down:
	migrate -database ${USERS_DB_URL} -path db/migrations/users down

migrate_transitions_down:
	migrate -database ${TRANSITIONS_DB_URL} -path db/migrations/transitions down

.PHONY: migrate_users_up migrate_users_down migrate_transitions_up migrate_transitions_down
