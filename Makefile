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
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Delhi.db','--http-addr=127.0.0.2:8080','--config-file=sharding.toml','--shard=Delhi'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Delhi-R.db','--http-addr=127.0.0.22:8080','--config-file=sharding.toml','--shard=Delhi', '--replica'"

	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Mumbai.db','--http-addr=127.0.0.3:8080','--config-file=sharding.toml','--shard=Mumbai'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Mumbai-R.db','--http-addr=127.0.0.33:8080','--config-file=sharding.toml','--shard=Mumbai', '--replica'"

	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Hyderabad.db','--http-addr=127.0.0.4:8080','--config-file=sharding.toml','--shard=Hyderabad'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Hyderabad-R.db','--http-addr=127.0.0.44:8080','--config-file=sharding.toml','--shard=Hyderabad', '--replica'"

	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Chennai.db','--http-addr=127.0.0.5:8080','--config-file=sharding.toml','--shard=Chennai'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Chennai-R.db','--http-addr=127.0.0.55:8080','--config-file=sharding.toml','--shard=Chennai', '--replica'"
	
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
	@echo "Removing all .db files on Windows..."
	@powershell -Command "Get-ChildItem -Filter *.db | Remove-Item -Force"
else
	@echo "Removing all .db files on Unix/Linux/macOS..."
	@rm -f *.db
endif
	@echo "✅ All .db files removed."


# Load test using curl to hit shard endpoints with unique data
load-test:
	@for SHARD in 127.0.0.1:8080 127.0.0.1:8081; do \
		for i in $$(seq 1 1000); do \
			RANDOM_NUM=$$(od -An -N2 -i /dev/urandom | tr -d ' '); \
			echo "[$$i] Sending to $$SHARD with key=key-$$RANDOM_NUM"; \
			curl -s "http://$$SHARD/set?key=key-$$RANDOM_NUM&value=value-$$RANDOM_NUM" > /dev/null; \
		done; \
	done
