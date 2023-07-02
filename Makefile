BINARY_NAME=api.wisper

build:
	go build -o bin/main main.go

compile:
 GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
 GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
 GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

run:
	go run main.go

clean:
	go clean
	rm -f bin/*

watch:
	air