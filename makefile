CONFIG_FILE_PATH=$(HOME)/.config/gitnotes
CONFIG_FILE=gn.conf

all: build

.PHONY: build
build:
	go build -o ./dist/gn main.go

# compile and run edit command
run:
	go run main.go edit

test:
	go test ./...

fmt:
	go fmt ./...

# install gitnotes into /usr/local/bin
# needs sudo
.PHONY: install
install: ./dist/gn
	sudo cp ./dist/gn /usr/local/bin/gn
	mkdir -p $(CONFIG_FILE_PATH)

.PHONY: lint
lint:
	golangci-lint run . --enable-all

.PHONY: clean
clean:
	rm ./dist/gn

./dist/gn:
	make build