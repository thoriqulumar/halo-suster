# README

DB migration article https://www.freecodecamp.org/news/database-migration-golang-migrate/

# How to create migration file

migrate create -ext sql -dir database/migration/ -seq {migration_name} ex: migrate create -ext sql -dir database/migration/ -seq cat

# Run Migration

run the migration migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose up