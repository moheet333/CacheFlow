build:
	@echo "Building..."
	@go build -o CacheFlow.exe main.go

demo:
	@go build -o CacheFlow.exe main.go
	@echo "Running..."
	@"CacheFlow.exe" --port 8000 --origin https://dummyjson.com && exit

clear:
	@go build -o CacheFlow.exe main.go
	@echo "Running..."
	@"CacheFlow.exe" --port 8000 --origin https://dummyjson.com --clear-cache && exit

test:
	@go clean -testcache
	@echo "Testing..."
	@go test ./... -v

docker-run:
	@docker compose up --build

docker-down:
	@docker compose down