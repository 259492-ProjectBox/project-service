include .env
gen:
	openapi-generator-cli generate -i docs/swagger.yaml -g typescript-axios -o ../frontend/generated/server
migratecreate:
	migrate create -ext sql -dir db/migrations -seq $(name)
migrateup:
	migrate -path db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_NAME}?sslmode=${POSTGRES_SSLMODE}" -verbose up 1
migratedown:
	migrate -path db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_NAME}?sslmode=${POSTGRES_SSLMODE}" -verbose down 1 
migrateclean:
	migrate -path db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_NAME}?sslmode=${POSTGRES_SSLMODE}" force 2
.PHONY: seed migrateup migratedown	