.ONESHELL:

# Params
BINARY_PATH=bin
BINARY_NAME=sample_go

MOD="./src/..."
COVER="coverage.out"
# Verbose mode. Set `make -e VERB=-v command`
VERB=

# https://github.com/confluentinc/confluent-kafka-go/blob/master/kafka/README.md#build-tags
# -tags musl is used for librdkafka
build:
	CGO_ENABLED=1 go build -tags musl -ldflags '-linkmode external -w -extldflags "-static"' -o $(BINARY_PATH)/${BINARY_NAME} cmd/main.go

lint:
	golangci-lint --timeout=5m run

watch:
	ulimit -n 2048 && air

run:
	APPSERVER_PORT=8080 go run cmd/main.go run 
run-with-logs:
	@mkdir -p logs
	@LOG_FILE="logs/app-$$(date +'%Y-%m-%d').log" && \
	APPSERVER_PORT=8080 go run cmd/main.go run 2>&1 | tee -a "$$LOG_FILE"

generate:
	go generate $(VERB) ./...

generate-admin-token:
	go run script/generateToken.go

test-run:
	GIN_MODE=release go test -race $(VERB) $(MOD)

# Run test with coverage
test-cover:
	GIN_MODE=release CGO_ENABLED=1 go test -race -coverprofile=$(COVER) $(VERB) $(MOD)

html-coverage-report: test-cover
	go tool cover -html=$(COVER) -o coverage.html

migrate:
	go run cmd/main.go --migrate-all migrate

rollback-last:
	go run cmd/main.go --rollback-last migrate

seed-db:
	go run cmd/main.go seed

generate-migration:
	sh ./infra/dev/migration/migration_generator.sh ${name}

db-reset: SHELL := zsh
db-reset:
	PGPASSWORD=postgres psql -h localhost -U postgres -d postgres -c  'DROP SCHEMA IF EXISTS public CASCADE'
	PGPASSWORD=postgres psql -h localhost -U postgres -d postgres -c  'CREATE SCHEMA public'

dev-setup-up: dev-infra-up migrate seed-db

dev-infra-up:
	docker-compose -f infra/dev/docker/docker-compose.yaml up -d

dev-setup-destroy:
	docker-compose -f infra/dev/docker/docker-compose.yaml down -v

dev-setup-reset: dev-setup-destroy dev-setup-up

service-api-test:
	go test -coverpkg=./... -coverprofile="coverage.out" github.com/Buddy-Git/comms-svc/test/servicetest

api-test-coverage:
	go test -coverpkg=./... -coverprofile="api_coverage.out" github.com/Buddy-Git/comms-svc/test/servicetest/...
	scoverage -package=github.com/Buddy-Git/comms-svc -coverageOutPath=api_coverage.out -minCov=15

test-coverage-overall:
	GIN_MODE=release go test -race -coverprofile="unit_coverage.out" ./src/...
	go test -coverpkg=./... -coverprofile="api_coverage.out" github.com/Buddy-Git/comms-svc/test/servicetest/notifyconsumer/...
	go test -coverpkg=./... -coverprofile="api_coverage.out" github.com/Buddy-Git/comms-svc/test/servicetest/webserver/...
	go test -coverpkg=./... -coverprofile="api_coverage.out" github.com/Buddy-Git/comms-svc/test/servicetest/webhook/...
	go test -coverpkg=./... -coverprofile="api_coverage.out" github.com/Buddy-Git/comms-svc/test/servicetest/webhookconsumer/...
	./infra/dev/coverage-tools/merge_file.sh
	scoverage -package=github.com/Buddy-Git/comms-svc -minCov=75