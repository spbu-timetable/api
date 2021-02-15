build:
	go build ./cmd/server

run:
	go run ./cmd/server

swagger:
	$$HOME/go/bin/swag init -d ./cmd/server/


.DEFAULT_GOAL: run