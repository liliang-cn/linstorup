GO_CMD=go
APP_NAME=linstorup
BUILD_DIR=./bin
MAIN_PACKAGE=./cmd/$(APP_NAME)

.PHONY: all build run test clean

all: build

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running tests..."
	$(GO_CMD) test ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f inventory.ini playbook.yml
	@echo "Clean complete."
