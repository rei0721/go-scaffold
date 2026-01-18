# Go Scaffold Makefile
# 提供常用的开发、构建、测试和部署命令

# 变量定义
APP_NAME := go-scaffold
VERSION := 0.1.2
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GO_VERSION := $(shell go version | awk '{print $$3}')

# 构建标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# 目录定义
BIN_DIR := bin
DOCS_DIR := docs
CONFIGS_DIR := configs
LOGS_DIR := logs

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "Go Scaffold v$(VERSION) - 可用命令:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "示例:"
	@echo "  make dev          # 开发模式启动"
	@echo "  make build        # 构建应用"
	@echo "  make test         # 运行测试"
	@echo "  make docker       # 构建 Docker 镜像"

# 开发相关命令
.PHONY: dev
dev: ## 开发模式启动 (热重载)
	@echo "🚀 启动开发模式..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "⚠️  Air 未安装，使用普通模式启动"; \
		go run cmd/server/main.go server; \
	fi

.PHONY: run
run: ## 运行应用
	@echo "🚀 启动应用..."
	go run cmd/server/main.go server

.PHONY: initdb
initdb: ## 初始化数据库
	@echo "🗄️ 初始化数据库..."
	go run cmd/server/main.go initdb

.PHONY: test-db
test-db: ## 测试数据库连接
	@echo "🔍 测试数据库连接..."
	go run cmd/server/main.go tests

# 构建相关命令
.PHONY: build
build: clean ## 构建应用
	@echo "🔨 构建应用..."
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME) cmd/server/main.go
	@echo "✅ 构建完成: $(BIN_DIR)/$(APP_NAME)"

.PHONY: build-linux
build-linux: clean ## 构建 Linux 版本
	@echo "🔨 构建 Linux 版本..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME)-linux cmd/server/main.go
	@echo "✅ Linux 版本构建完成: $(BIN_DIR)/$(APP_NAME)-linux"

.PHONY: build-windows
build-windows: clean ## 构建 Windows 版本
	@echo "🔨 构建 Windows 版本..."
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME).exe cmd/server/main.go
	@echo "✅ Windows 版本构建完成: $(BIN_DIR)/$(APP_NAME).exe"

.PHONY: build-all
build-all: build build-linux build-windows ## 构建所有平台版本
	@echo "✅ 所有平台构建完成"

# 测试相关命令
.PHONY: test
test: ## 运行所有测试
	@echo "🧪 运行测试..."
	go test -v -race ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "🧪 运行测试覆盖率..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 覆盖率报告生成: coverage.html"

.PHONY: test-short
test-short: ## 运行快速测试 (跳过集成测试)
	@echo "🧪 运行快速测试..."
	go test -short ./...

.PHONY: test-integration
test-integration: ## 运行集成测试
	@echo "🧪 运行集成测试..."
	go test -tags=integration ./...

.PHONY: benchmark
benchmark: ## 运行基准测试
	@echo "⚡ 运行基准测试..."
	go test -bench=. -benchmem ./...

# 代码质量命令
.PHONY: fmt
fmt: ## 格式化代码
	@echo "🎨 格式化代码..."
	go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi

.PHONY: lint
lint: ## 代码检查
	@echo "🔍 代码检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint 未安装，跳过代码检查"; \
		echo "安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: vet
vet: ## Go vet 检查
	@echo "🔍 Go vet 检查..."
	go vet ./...

.PHONY: security
security: ## 安全检查
	@echo "🔒 安全检查..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec 未安装，跳过安全检查"; \
		echo "安装命令: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

.PHONY: check
check: fmt vet lint test ## 运行所有检查 (格式化、检查、测试)

# 依赖管理命令
.PHONY: deps
deps: ## 下载依赖
	@echo "📦 下载依赖..."
	go mod download

.PHONY: deps-update
deps-update: ## 更新依赖
	@echo "📦 更新依赖..."
	go get -u ./...
	go mod tidy

.PHONY: deps-verify
deps-verify: ## 验证依赖
	@echo "🔍 验证依赖..."
	go mod verify

.PHONY: deps-clean
deps-clean: ## 清理依赖缓存
	@echo "🧹 清理依赖缓存..."
	go clean -modcache

# Docker 相关命令
.PHONY: docker
docker: ## 构建 Docker 镜像
	@echo "🐳 构建 Docker 镜像..."
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "✅ Docker 镜像构建完成: $(APP_NAME):$(VERSION)"

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	@echo "🐳 运行 Docker 容器..."
	docker run -p 8080:8080 --name $(APP_NAME) $(APP_NAME):latest

.PHONY: docker-compose-up
docker-compose-up: ## 启动 Docker Compose 服务
	@echo "🐳 启动 Docker Compose 服务..."
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## 停止 Docker Compose 服务
	@echo "🐳 停止 Docker Compose 服务..."
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## 查看 Docker Compose 日志
	@echo "📋 查看 Docker Compose 日志..."
	docker-compose logs -f

# 清理命令
.PHONY: clean
clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

.PHONY: clean-all
clean-all: clean ## 清理所有生成文件 (包括日志)
	@echo "🧹 清理所有文件..."
	@rm -rf $(LOGS_DIR)/*.log
	@rm -rf tmp/
	@echo "✅ 深度清理完成"

# 配置相关命令
.PHONY: config-init
config-init: ## 初始化配置文件
	@echo "⚙️ 初始化配置文件..."
	@if [ ! -f $(CONFIGS_DIR)/config.yaml ]; then \
		cp $(CONFIGS_DIR)/config.example.yaml $(CONFIGS_DIR)/config.yaml; \
		echo "✅ 配置文件已创建: $(CONFIGS_DIR)/config.yaml"; \
	else \
		echo "⚠️  配置文件已存在: $(CONFIGS_DIR)/config.yaml"; \
	fi
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ 环境变量文件已创建: .env"; \
	else \
		echo "⚠️  环境变量文件已存在: .env"; \
	fi

.PHONY: config-validate
config-validate: ## 验证配置文件
	@echo "🔍 验证配置文件..."
	go run cmd/server/main.go server --config=$(CONFIGS_DIR)/config.yaml --dry-run

# 文档相关命令
.PHONY: docs
docs: ## 生成文档
	@echo "📚 生成文档..."
	@echo "文档已存在于 $(DOCS_DIR)/ 目录"
	@echo "可以通过以下方式查看:"
	@echo "  - 在线查看: 启动应用后访问 http://localhost:8080/docs"
	@echo "  - 本地查看: 使用 Markdown 阅读器打开 $(DOCS_DIR)/README.md"

.PHONY: docs-serve
docs-serve: ## 启动文档服务器
	@echo "📚 启动文档服务器..."
	@if command -v python3 > /dev/null; then \
		echo "📖 文档服务器启动: http://localhost:8000"; \
		cd $(DOCS_DIR) && python3 -m http.server 8000; \
	elif command -v python > /dev/null; then \
		echo "📖 文档服务器启动: http://localhost:8000"; \
		cd $(DOCS_DIR) && python -m SimpleHTTPServer 8000; \
	else \
		echo "⚠️  Python 未安装，无法启动文档服务器"; \
	fi

# 安装工具命令
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "🔧 安装开发工具..."
	@echo "安装 Air (热重载)..."
	go install github.com/cosmtrek/air@latest
	@echo "安装 golangci-lint (代码检查)..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "安装 gosec (安全检查)..."
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "安装 goimports (导入整理)..."
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ 开发工具安装完成"

# 版本信息
.PHONY: version
version: ## 显示版本信息
	@echo "应用名称: $(APP_NAME)"
	@echo "版本号: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git 提交: $(GIT_COMMIT)"
	@echo "Go 版本: $(GO_VERSION)"

# 环境检查
.PHONY: env-check
env-check: ## 检查开发环境
	@echo "🔍 检查开发环境..."
	@echo "Go 版本:"
	@go version
	@echo ""
	@echo "Git 版本:"
	@git --version
	@echo ""
	@echo "Docker 版本:"
	@docker --version 2>/dev/null || echo "Docker 未安装"
	@echo ""
	@echo "Make 版本:"
	@make --version | head -n1
	@echo ""
	@echo "开发工具检查:"
	@command -v air > /dev/null && echo "✅ Air 已安装" || echo "❌ Air 未安装"
	@command -v golangci-lint > /dev/null && echo "✅ golangci-lint 已安装" || echo "❌ golangci-lint 未安装"
	@command -v gosec > /dev/null && echo "✅ gosec 已安装" || echo "❌ gosec 未安装"
	@command -v goimports > /dev/null && echo "✅ goimports 已安装" || echo "❌ goimports 未安装"

# 快速启动命令
.PHONY: quick-start
quick-start: config-init deps build ## 快速启动 (初始化配置、下载依赖、构建)
	@echo "🚀 快速启动完成！"
	@echo "运行以下命令启动应用:"
	@echo "  make run"
	@echo "或者:"
	@echo "  ./$(BIN_DIR)/$(APP_NAME) server"

# 生产部署准备
.PHONY: prod-build
prod-build: clean fmt vet test build-linux ## 生产环境构建 (清理、格式化、检查、测试、构建)
	@echo "🚀 生产环境构建完成！"
	@echo "可执行文件: $(BIN_DIR)/$(APP_NAME)-linux"

# CI/CD 相关命令
.PHONY: ci
ci: deps fmt vet lint test ## CI 流水线 (下载依赖、格式化、检查、测试)
	@echo "✅ CI 流水线执行完成"

.PHONY: cd
cd: ci build docker ## CD 流水线 (CI + 构建 + Docker)
	@echo "✅ CD 流水线执行完成"

# 监控和调试
.PHONY: pprof
pprof: ## 启动性能分析
	@echo "📊 启动性能分析..."
	@echo "确保应用正在运行，然后访问:"
	@echo "  CPU: go tool pprof http://localhost:8080/debug/pprof/profile"
	@echo "  内存: go tool pprof http://localhost:8080/debug/pprof/heap"
	@echo "  协程: go tool pprof http://localhost:8080/debug/pprof/goroutine"

.PHONY: health-check
health-check: ## 健康检查
	@echo "🏥 执行健康检查..."
	@curl -f http://localhost:8080/health || echo "❌ 健康检查失败"
	@echo ""
	@curl -f http://localhost:8080/health/detailed || echo "❌ 详细健康检查失败"

# 数据库相关
.PHONY: db-reset
db-reset: ## 重置数据库
	@echo "🗄️ 重置数据库..."
	@echo "⚠️  这将删除所有数据！"
	@read -p "确认继续? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	go run cmd/server/main.go initdb --force

# 日志相关
.PHONY: logs
logs: ## 查看应用日志
	@echo "📋 查看应用日志..."
	@if [ -f $(LOGS_DIR)/app.log ]; then \
		tail -f $(LOGS_DIR)/app.log; \
	else \
		echo "❌ 日志文件不存在: $(LOGS_DIR)/app.log"; \
	fi

.PHONY: logs-clean
logs-clean: ## 清理日志文件
	@echo "🧹 清理日志文件..."
	@rm -f $(LOGS_DIR)/*.log
	@echo "✅ 日志文件清理完成"

# 备份和恢复
.PHONY: backup
backup: ## 备份数据库
	@echo "💾 备份数据库..."
	@mkdir -p backups
	@mysqldump -u root -p scaffold > backups/scaffold_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "✅ 数据库备份完成"

.PHONY: restore
restore: ## 恢复数据库 (需要指定备份文件)
	@echo "🔄 恢复数据库..."
	@if [ -z "$(FILE)" ]; then \
		echo "❌ 请指定备份文件: make restore FILE=backups/scaffold_20240101_120000.sql"; \
		exit 1; \
	fi
	@mysql -u root -p scaffold < $(FILE)
	@echo "✅ 数据库恢复完成"