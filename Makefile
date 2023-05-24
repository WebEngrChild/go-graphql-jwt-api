init:
	@make build
	@make up
build:
	docker-compose build --no-cache
up:
	docker-compose up -d
app:
	docker-compose exec app sh
mysql:
	docker-compose exec mysql mysql -uroot -puser
migrate:
	docker-compose exec app goose -dir ./build/db/migration mysql "root:user@tcp(mysql:3306)/golang" up # mysqlはコンテナ名
create-migration: # ファイル名は適宜変更すること
	docker-compose exec app goose -dir ./build/db/migration create insert_users sql
create-mock: # ファイル名は適宜変更すること
	docker-compose exec app mockgen -source=pkg/domain/repository/user_repository.go -destination pkg/lib/mock/mock_user.go
start:
	docker-compose exec app go run ./cmd/main.go
down:
	docker-compose down
stop:
	docker-compose stop
gqlgen:
	docker-compose exec app go run github.com/99designs/gqlgen generate
air:
	docker-compose exec app air -c .air.toml
app-dlv:  # コンテナ内でdlv直接実行
	docker-compose exec app dlv debug ./cmd/main.go
dlv:  #　Golandを使ったdlv
	docker-compose exec app dlv debug ./cmd/main.go --headless --listen=:2345 --api-version=2 --accept-multiclient
dry-fix:
	golangci-lint run ./...
fix:
	golangci-lint run --fix
test:
	go test ./... -coverprofile=coverage.out


.PHONY: build