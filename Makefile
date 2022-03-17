BINARY_NAME=t-mirror
build:
	go build -o t-mirror .

run:
	go run main.go


compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME} .
	GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME} .
	GOOS=freebsd GOARCH=386 go build -o ${BINARY_NAME} .

all: build compile run
