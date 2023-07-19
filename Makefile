build:
	docker build --tag mediasoft-statistics .

migrateUp:
	migrate -path ./migrations -database "postgresql://$PG_USER:$PG_PWD@$PG_HOST:$PG_PORT/$PG_DATABASE?sslmode=disable" up
