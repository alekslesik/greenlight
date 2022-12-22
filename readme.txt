Installs:
1. Install posgre sql https://www.postgresql.org/download/
2. Install migrate https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
3. Download a pre-built binary and move it to a location on your system path https://github.com/golang-migrate/migrate/releases

Migrate:
migrate create -seq -ext .sql -dir ./migrations create_tokens_table
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
