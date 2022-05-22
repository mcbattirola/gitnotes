build:
	go build -o ./dist/gn main.go

run:
	go run main.go edit

test:
	go test ./...