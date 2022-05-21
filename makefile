build:
	go build -o ./dist/gn main.go

test:
	go test ./...