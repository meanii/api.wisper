BINARY_NAME=api.wisper
BIN_DIR=bin

build:
	@mkdir -p ${BINARY_NAME}
	@echo "Building ${BINARY_NAME} for current platform..."
	@go build -o ${BINARY_NAME}/${BINARY_NAME} main.go

compile:
	@mkdir -p ${BIN_DIR}
	@echo "Building ${BINARY_NAME} for multiple platforms..."
  GOARCH=amd64 GOOS=darwin go build -o ${BIN_DIR}/${BINARY_NAME}-darwin main.go
	GOARCH=arm64 GOOS=darwin go build -o ${BIN_DIR}/${BINARY_NAME}-darwin-arm64 main.go
  GOARCH=amd64 GOOS=linux go build -o ${BIN_DIR}/${BINARY_NAME}-linux main.go
	GOARCH=arm64 GOOS=linux go build -o ${BIN_DIR}/${BINARY_NAME}-linux-arm64 main.go
  GOARCH=amd64 GOOS=windows go build -o ${BIN_DIR}/${BINARY_NAME}-windows main.go

run: build
	@echo "Running ${BINARY_NAME}..."
	@./${BINARY_NAME}/${BINARY_NAME}

dep:
	@echo "Installing dependencies..."
	@go mod tidy

clean:
	@echo "Cleaning up..."
	@rm -rf ${BIN_DIR}
	@go clean

watch:
	echo "Watching for changes..."
	@nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go

