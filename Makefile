APP_NAME = kvs
BUILD_DIR = build
OUTPUT = $(BUILD_DIR)/$(APP_NAME)

all: build

test:
	go test 

build:
	go build -o OUTPUT .

run: test build
	./$(OUTPUT)

lint:
	golangci-lint run

clean:
	rm -f $(APP_NAME)
