CONFIG_FILE_PATH=$(HOME)/.config/gitnotes
CONFIG_FILE=gn.conf

./dist/gn:
	make build

.PHONY: build
build:
	go build -o ./dist/gn main.go

run:
	go run main.go edit

test:
	go test ./...

# install gitnotes into /usr/local/bin
# needs sudo
.PHONY: install
install: ./dist/gn
	sudo cp ./dist/gn /usr/local/bin/gn
	mkdir -p $(CONFIG_FILE_PATH)

.PHONY: lint
ling:
	golangci-lint run . --enable-all