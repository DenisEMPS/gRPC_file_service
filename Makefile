CONFIG_PATH=./config/config.yml

.PHONY: run

run:
	CONFIG_PATH=${CONFIG_PATH} go run cmd/file-service/main.go

build:
	go build cmd/file-service/main.go
