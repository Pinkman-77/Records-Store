migrator:
migrate -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" -path migrations up
