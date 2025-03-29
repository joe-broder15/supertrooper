.PHONY: build clean certs

build:
	@if [ ! -d "./build" ]; then mkdir -p ./build; fi
	GOOS=windows GOARCH=amd64 go build -o ./build/server.exe ./cmd/server/main.go 
	GOOS=windows GOARCH=amd64 go build -o ./build/agent.exe ./cmd/agent/main.go 
	upx ./build/agent.exe

clean:
	rm -rf ./build/*