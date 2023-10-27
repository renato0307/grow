BINARY_NAME=./bin/monitor

build:
	env GOOS=linux GOARCH=arm GOARM=6 go build -o ${BINARY_NAME}
	
publish:
	scp ${BINARY_NAME} growzero1:~/