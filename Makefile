CONFIG_PATH=./config/config.yml

.PHONY: run, build

run:
	CONFIG_PATH=${CONFIG_PATH} go run cmd/file-service/main.go

build:
	go build cmd/file-service/main.go
