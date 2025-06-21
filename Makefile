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
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Delhi-R.db','--http-addr=127.0.0.22:8080','--config-file=sharding.toml','--shard=Delhi' --replica"

	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Mumbai.db','--http-addr=127.0.0.3:8080','--config-file=sharding.toml','--shard=Mumbai'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Mumbai-R.db','--http-addr=127.0.0.33:8080','--config-file=sharding.toml','--shard=Mumbai' --replica"

	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Hyderabad.db','--http-addr=127.0.0.4:8080','--config-file=sharding.toml','--shard=Hyderabad'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Hyderabad-R.db','--http-addr=127.0.0.44:8080','--config-file=sharding.toml','--shard=Hyderabad' --replica"

	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Chennai.db','--http-addr=127.0.0.5:8080','--config-file=sharding.toml','--shard=Chennai'"
	@powershell -Command "Start-Process $(BINARY_NAME) -ArgumentList '--db-location=Chennai-R.db','--http-addr=127.0.0.55:8080','--config-file=sharding.toml','--shard=Chennai' --replica"
	
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
	@del /Q Chennai.db 2>nul || echo Chennai.db not found
	
else
	@rm -f Delhi.db Mumbai.db Hyderabad.db
endif
	@echo "Database files removed."

# Load test using curl to hit shard endpoints with unique data
load-test:
	@for SHARD in 127.0.0.1:8080 127.0.0.1:8081; do \
		for i in $$(seq 1 1000); do \
			RANDOM_NUM=$$(od -An -N2 -i /dev/urandom | tr -d ' '); \
			echo "[$$i] Sending to $$SHARD with key=key-$$RANDOM_NUM"; \
			curl -s "http://$$SHARD/set?key=key-$$RANDOM_NUM&value=value-$$RANDOM_NUM" > /dev/null; \
		done; \
	done



# [352] Sending to 127.0.0.1:8080 with key=key-50791	1
# [353] Sending to 127.0.0.1:8080 with key=key-19242
# [354] Sending to 127.0.0.1:8080 with key=key-58472
# [355] Sending to 127.0.0.1:8080 with key=key-54727	0
# [356] Sending to 127.0.0.1:8080 with key=key-53263

# [996] Sending to 127.0.0.1:8081 with key=key-53541	1
# [997] Sending to 127.0.0.1:8081 with key=key-29244
# [998] Sending to 127.0.0.1:8081 with key=key-59131
# [999] Sending to 127.0.0.1:8081 with key=key-43965
# [1000] Sending to 127.0.0.1:8081 with key=key-16017	0	2