build:
	@echo "Building..."
	@go build -o CacheFlow.exe main.go

demo:
	@go build -o CacheFlow.exe main.go
	@echo "Running..."
	@"CacheFlow.exe" --port 8000 --origin https://dummyjson.com && exit

test:
	@echo "Testing..."
	@go test ./... -v