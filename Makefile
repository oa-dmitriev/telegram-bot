all: build 

build: .
	go build -o ./bin/mybot ./cmd/mybot/main.go

run: build
	./bin/mybot
