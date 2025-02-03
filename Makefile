build:
	@echo "Building..."
	@go build -o CacheFlow.exe main.go

all:
	@go build -o CacheFlow.exe main.go
	@./CacheFlow.exe --port 8000 --origin https://dummyjson.com