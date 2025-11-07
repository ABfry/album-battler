up:
	docker compose up --build

lint:
	cd backend && GOEXPERIMENT=synctest go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run --fix ./...

lint-ci:
	cd backend && GOEXPERIMENT=synctest go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...
