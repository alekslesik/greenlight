# Include variables from the .envrc file
include .envrc

#=====================================#
# HELPERS #
#=====================================#

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

#=====================================#
# DEVELOPMENT #
#=====================================#

## run/api: run the cmd/api application
.PHONY: run
run:
	go run ./cmd/api/

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

#=====================================#
# QUALITY CONTROL #
#=====================================#

## audit: tidy dependencies and format, vet and test all code

## go fmt ./... : command to format all .go files in the project directory, according to the Go standard.
## go vet ./... : runs a variety of analyzers which carry out static analysis of your code and warn you
## go test -race -vet=off ./... : command to run all tests in the project directory
## staticcheck tool : to carry out some additional static analysis checks.
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## go mod tidy : prune any unused dependencies from the go.mod and go.sum files, and add any missing dependencies
## go mod verify : check that the dependencies on your computer (located in your module cache located at $GOPATH/pkg/mod)
## haven’t been changed since they were downloaded and that they match the cryptographic hashes in your go.sum file
## go mod vendor: copy the necessary source code from your module cache into a new vendor directory in your project root
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

#=====================================#
# BUILD #
#=====================================#

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build
build:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

#=====================================#
# PRODUCTION #
#=====================================#

production_host_ip = '82.146.47.139'
## production/connect: connect to the production server
.PHONY: connect 
connect:
	ssh greenlight@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: deploy
deploy:
	rsync -rP --delete ./bin/linux_amd64/api ./migrations greenlight@${production_host_ip}:~
	ssh -t greenlight@${production_host_ip} 'migrate -path ~/migrations -database $$GREENLIGHT_DB_DSN up'
	ssh -t greenlight@${production_host_ip} 'chmod 744 api'

## production/configure/api.service: configure the production systemd api.service file
.PHONY: production/configure/api.service
production/configure/api.service:
	rsync -P ./remote/production/api.service greenlight@${production_host_ip}:~
	ssh -t greenlight@${production_host_ip} 'sudo mv ~/api.service /etc/systemd/system/ && sudo systemctl enable api && sudo systemctl restart api'

## production/configure/caddyfile: configure the production Caddyfile
.PHONY: production/configure/caddyfile
production/configure/caddyfile:
	rsync -P ./remote/production/Caddyfile greenlight@${production_host_ip}:~
	ssh -t greenlight@${production_host_ip} 'sudo mv ~/Caddyfile /etc/caddy/ && sudo systemctl reload caddy'
