# Set default target
.DEFAULT_GOAL := run-multi

# Detect OS and set binary name and commands
ifeq ($(OS),Windows_NT)
    BINARY_NAME=distribkv.exe
    RUN_BG=powershell -Command "Start-Process"
    KILL_CMD=taskkill /F /IM $(BINARY_NAME) 2>nul || echo No running instances to kill
else
    BINARY_NAME=distribkv
    RUN_BG=sh -c
    KILL_CMD=pkill -f "./$(BINARY_NAME)" || echo No running instances to kill
endif

# Build the Go binary
build:
	@go build -o $(BINARY_NAME) main.go

# Run a single instance
run: build
	@./$(BINARY_NAME) --db-location=./Delhi.db --http-addr=127.0.0.1:8080 --config-file=sharding.toml --shard=Delhi

# Run three instances in background
run-multi: build
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Delhi.db','--http-addr=127.0.0.1:8080','--config-file=sharding.toml','--shard=Delhi'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Mumbai.db','--http-addr=127.0.0.1:8081','--config-file=sharding.toml','--shard=Mumbai'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Hyderabad.db','--http-addr=127.0.0.1:8082','--config-file=sharding.toml','--shard=Hyderabad'"
else
	@./$(BINARY_NAME) --db-location=Delhi.db --http-addr=127.0.0.1:8080 --config-file=sharding.toml --shard=Delhi &
	@./$(BINARY_NAME) --db-location=Mumbai.db --http-addr=127.0.0.1:8081 --config-file=sharding.toml --shard=Mumbai &
	@./$(BINARY_NAME) --db-location=Hyderabad.db --http-addr=127.0.0.1:8082 --config-file=sharding.toml --shard=Hyderabad &
endif

# Stop all running instances
stop:
	@$(KILL_CMD)

# Run Go tests
test:
	@go test -v ./...

# Remove database files
remove:
ifeq ($(OS),Windows_NT)
	@del /Q Delhi.db 2>nul || echo Delhi.db not found
	@del /Q Mumbai.db 2>nul || echo Mumbai.db not found
	@del /Q Hyderabad.db 2>nul || echo Hyderabad.db not found
else
	@rm -f Delhi.db Mumbai.db Hyderabad.db
endif
	@echo "Database files removed."
