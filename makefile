CONFIG_FILE_PATH=$HOME/.config/gitnotes
CONFIG_FILE=gn.config

build:
	go build -o ./dist/gn main.go

run:
	go run main.go edit

test:
	go test ./...

install: ./dist/gn
  cp ./dist/gn /user/local/bin/gn
	mkdir -p $(CONFIG_FILE_PATH)
	touch $(CONFIG_FILE_PATH)/$(CONFIG_FILE)
