APP_NAME=music_app
CMD_PATH=./cmd/music_app
BIN_DIR=bin

.PHONY: run build clean test docker-build docker-run

run:
	go run $(CMD_PATH)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./...

docker-build:
	docker build -t albums-api .

docker-run:
	docker run -p 8080:8080 albums-api
