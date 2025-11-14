up:
	docker compose up --build

upd:
	docker compose up --build -d

down: 
	docker compose down

downv:
	docker compose down -v

log:
	docker compose logs -f

lint:
	cd backend && GOEXPERIMENT=synctest go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run --fix ./...

lint-ci:
	cd backend && GOEXPERIMENT=synctest go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...

storybook:
	cd frontend && npm run storybook

storybook-stop:
	pkill -f "storybook dev" || echo "Storybookプロセスが見つかりませんでした"
