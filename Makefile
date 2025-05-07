APP_NAME = Provider-Pairing-System
CMD_DIR = ./cmd
BIN_DIR = ./bin
BIN_FILE = $(BIN_DIR)/$(APP_NAME)

# ANSI colors
GREEN = \033[0;32m
BLUE = \033[0;34m
NC = \033[0m  # No Color

.PHONY: default build run clean

# Default target: build then run
default: build run

# Build the application
build:
	@echo "$(BLUE)==> Tidying modules...$(NC)"
	@go mod tidy
	@echo "$(BLUE)==> Creating bin directory...$(NC)"
	@mkdir -p $(BIN_DIR)
	@echo "$(BLUE)==> Building binary...$(NC)"
	@go build -o $(BIN_FILE) $(CMD_DIR)
	@echo "$(GREEN)==> Build complete: $(BIN_FILE)$(NC)"

# Run the application (without building binary)
run:
	@echo "$(BLUE)==> Running application...$(NC)"
	@go run $(CMD_DIR)
	@echo "$(GREEN)==> Run complete$(NC)"

# Clean up built files
clean:
	@echo "$(BLUE)==> Cleaning up...$(NC)"
	@rm -rf $(BIN_DIR)
	@echo "$(GREEN)==> Clean complete$(NC)"
