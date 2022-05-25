CONFIG_FILE_PATH=$(HOME)/.config/gitnotes
CONFIG_FILE=gn.config

./dist/gn:
	make build

.PHONY: build
build:
	go build -o ./dist/gn main.go

run:
	go run main.go edit

test:
	go test ./...

# TODO make a default config file and copy it to .config
.PHONY: install
install: ./dist/gn
	sudo cp ./dist/gn /usr/local/bin/gn
	mkdir -p $(CONFIG_FILE_PATH)
	touch $(CONFIG_FILE_PATH)/$(CONFIG_FILE)

.PHONY: lint
ling:
	golangci-lint run . --enable-all