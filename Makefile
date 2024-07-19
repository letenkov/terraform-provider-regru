# Версия провайдера
VERSION := 0.2.1

# Определение пути для установки плагина
OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/letenkov/regru/$(VERSION)/$(OS_ARCH)
PLUGIN_NAME := terraform-provider-regru_v$(VERSION)
BUILD_OUTPUT := $(PLUGIN_DIR)/$(PLUGIN_NAME)

# Компиляция и установка провайдера
.PHONY: build
build:
	@echo "Building Terraform provider..."
	mkdir -p $(PLUGIN_DIR)
	go build -o $(BUILD_OUTPUT)
	@echo "Build completed and installed at $(BUILD_OUTPUT)"

# Получение версии Go из go.mod
.PHONY: go-version
go-version:
	@grep ^go go.mod | awk '{ print $$2 }'

# Запуск тестов
.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -v
	@echo "Tests completed"

# Форматирование кода
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Линтинг кода
.PHONY: lint
lint:
	@echo "Linting code..."
	go vet ./...
	@echo "Linting completed"

# Очистка сборочных артефактов
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BUILD_OUTPUT)
	@echo "Cleanup completed"

# Установка всех зависимостей
.PHONY: install-deps
install-deps:
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Dependencies installed"

# Помощь
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build          Build and install the Terraform provider"
	@echo "  make go-version Get the Go version from go.mod"
	@echo "  make test           Run tests"
	@echo "  make fmt            Format the code"
	@echo "  make lint           Lint the code"
	@echo "  make clean          Clean up build artifacts"
	@echo "  make install-deps   Install all dependencies"
	@echo "  make help           Display this help message"
