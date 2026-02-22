local:
	set -a; \
	. ./.env; \
	set +a; \
	go run cmd/projectservice/main.go --config=./config/local.yaml

build_proto:
	protoc --go_out=./proto/userservice --go_opt=paths=import --go-grpc_out=./proto/userservice --go-grpc_opt=paths=import ./proto/userservice/user.proto
	protoc --go_out=./proto/projectservice --go_opt=paths=import --go-grpc_out=./proto/projectservice --go-grpc_opt=paths=import ./proto/projectservice/project.proto

migrate_all_up:
	set -a; \
	. ./.env; \
	set +a; \
	migrate -path ./migrations -database "postgres://$$MIG_NAME:$$MIG_PASS@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_MODE" upMessageName