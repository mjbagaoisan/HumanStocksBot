include .env.local
export

# Status of the DB
status:
	goose -dir db/migrations status

# Apply all migrations (Up)
up:
	goose -dir db/migrations up

# Rollback one migration (Down)
down:
	goose -dir db/migrations down

# Create a new migration file (Usage: make create name=add_users)
create:
	goose -dir db/migrations create $(name) sql