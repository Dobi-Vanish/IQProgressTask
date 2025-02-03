SHELL=cmd.exe
FINANCIAL_BINARY=financialApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
run: build_financial
	@echo Stopping docker images (if running...)
	docker-compose down
	@echo Building (when required) and starting docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build_financial: builds the financial binary as a linux executable
build_financial:
	@echo Building financial binary...
	chdir ..\financial-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${FINANCIAL_BINARY} ./cmd/api
	@echo Done!

