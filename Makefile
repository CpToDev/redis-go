APP_NAME = kvs
BUILD_DIR = build
OUTPUT = $(BUILD_DIR)/$(APP_NAME)

all: build

test:
	go test 

build: test
	go build -o $(OUTPUT) 

run: build
	./$(OUTPUT)

lint:
	golangci-lint run

clean:
	rm -f $(OUTPUT)
