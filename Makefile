BINARY_NAME=logicmonitor
BUILD_DIR=build/
CMD=cmd/logicmonitor/logicmonitor

build:
	go mod vendor
	GOARCH=amd64 GOOS=linux go build -o ${BUILD_DIR}/${BINARY_NAME}_amd64_linux ${CMD}.go
	GOARCH=amd64 GOOS=windows go build -o ${BUILD_DIR}/${BINARY_NAME}_amd64_windows.exe ${CMD}.go

clean:
	go clean
	rm -rf ${BUILD_DIR}
