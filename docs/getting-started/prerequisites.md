# 环境准备

在开始使用 Go Scaffold 之前，请确保您的开发环境满足以下要求。

## 🔧 必需环境

### 1. Go 语言环境

**版本要求**: Go 1.24.6 或更高版本

**安装方法**:

#### macOS
```bash
# 使用 Homebrew
brew install go

# 或下载官方安装包
# https://golang.org/dl/
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL
sudo yum install golang

# 或使用官方二进制包
wget https://golang.org/dl/go1.24.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

#### Windows
1. 下载官方安装包: https://golang.org/dl/
2. 运行安装程序
3. 配置环境变量 `PATH`

**验证安装**:
```bash
go version
# 输出: go version go1.24.6 darwin/amd64
```

### 2. Git 版本控制

**版本要求**: Git 2.0 或更高版本

**安装方法**:
```bash
# macOS
brew install git

# Ubuntu/Debian
sudo apt install git

# CentOS/RHEL
sudo yum install git

# Windows
# 下载安装包: https://git-scm.com/download/win
```

**验证安装**:
```bash
git --version
# 输出: git version 2.39.0
```

## 🗄️ 数据库环境

### MySQL (推荐)

**版本要求**: MySQL 8.0 或更高版本

**安装方法**:
```bash
# macOS
brew install mysql
brew services start mysql

# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# 使用 Docker
docker run --name mysql-scaffold \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=scaffold \
  -p 3306:3306 \
  -d mysql:8.0
```

**配置数据库**:
```sql
-- 创建数据库
CREATE DATABASE scaffold CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（可选）
CREATE USER 'scaffold'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON scaffold.* TO 'scaffold'@'localhost';
FLUSH PRIVILEGES;
```

### PostgreSQL (可选)

**版本要求**: PostgreSQL 12 或更高版本

**安装方法**:
```bash
# macOS
brew install postgresql
brew services start postgresql

# Ubuntu/Debian
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 使用 Docker
docker run --name postgres-scaffold \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=scaffold \
  -p 5432:5432 \
  -d postgres:15
```

### SQLite (开发环境)

SQLite 是内置支持的，无需额外安装。适合开发和测试环境使用。

## 🚀 缓存服务

### Redis (可选但推荐)

**版本要求**: Redis 6.0 或更高版本

**安装方法**:
```bash
# macOS
brew install redis
brew services start redis

# Ubuntu/Debian
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server

# 使用 Docker
docker run --name redis-scaffold \
  -p 6379:6379 \
  -d redis:7-alpine
```

**验证安装**:
```bash
redis-cli ping
# 输出: PONG
```

## 🛠️ 开发工具

### 1. 代码编辑器

**推荐选择**:
- **VS Code** + Go 扩展
- **GoLand** (JetBrains)
- **Vim/Neovim** + vim-go
- **Emacs** + go-mode

**VS Code 配置**:
```json
{
  "go.toolsManagement.autoUpdate": true,
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v"],
  "go.coverOnSave": true
}
```

### 2. Go 工具链

**安装常用工具**:
```bash
# 代码格式化
go install golang.org/x/tools/cmd/goimports@latest

# 代码检查
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 热重载工具
go install github.com/cosmtrek/air@latest

# 依赖管理
go install golang.org/x/mod/cmd/go-mod-outdated@latest

# 测试覆盖率
go install github.com/axw/gocov/gocov@latest
go install github.com/matm/gocov-html@latest
```

### 3. Make 工具

**安装方法**:
```bash
# macOS (通常已预装)
xcode-select --install

# Ubuntu/Debian
sudo apt install build-essential

# CentOS/RHEL
sudo yum groupinstall "Development Tools"

# Windows
# 安装 MinGW 或使用 WSL
```

## 🐳 容器化环境 (可选)

### Docker

**版本要求**: Docker 20.0 或更高版本

**安装方法**:
```bash
# macOS
brew install --cask docker

# Ubuntu
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Windows
# 下载 Docker Desktop: https://www.docker.com/products/docker-desktop
```

### Docker Compose

**版本要求**: Docker Compose 2.0 或更高版本

通常随 Docker Desktop 一起安装，或单独安装：
```bash
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

## 🌐 网络工具

### cURL
用于 API 测试：
```bash
# macOS (通常已预装)
brew install curl

# Ubuntu/Debian
sudo apt install curl

# Windows
# 通常已预装，或使用 PowerShell 的 Invoke-RestMethod
```

### Postman (可选)
图形化 API 测试工具：
- 下载地址: https://www.postman.com/downloads/

## 📋 环境检查清单

在继续之前，请确认以下环境已正确安装：

- [ ] Go 1.24.6+ 已安装并配置
- [ ] Git 已安装
- [ ] 数据库已安装并运行 (MySQL/PostgreSQL/SQLite)
- [ ] Redis 已安装并运行 (可选)
- [ ] 代码编辑器已配置
- [ ] Go 工具链已安装
- [ ] Make 工具已安装
- [ ] Docker 已安装 (可选)

## 🔍 环境验证脚本

创建一个简单的验证脚本：

```bash
#!/bin/bash
# check-env.sh

echo "🔍 检查开发环境..."

# 检查 Go
if command -v go &> /dev/null; then
    echo "✅ Go: $(go version)"
else
    echo "❌ Go 未安装"
fi

# 检查 Git
if command -v git &> /dev/null; then
    echo "✅ Git: $(git --version)"
else
    echo "❌ Git 未安装"
fi

# 检查 Make
if command -v make &> /dev/null; then
    echo "✅ Make: $(make --version | head -n1)"
else
    echo "❌ Make 未安装"
fi

# 检查 Docker
if command -v docker &> /dev/null; then
    echo "✅ Docker: $(docker --version)"
else
    echo "⚠️  Docker 未安装 (可选)"
fi

# 检查数据库连接
echo "🗄️ 检查数据库连接..."

# MySQL
if command -v mysql &> /dev/null; then
    if mysql -u root -e "SELECT 1;" &> /dev/null; then
        echo "✅ MySQL 连接正常"
    else
        echo "⚠️  MySQL 连接失败"
    fi
else
    echo "⚠️  MySQL 未安装"
fi

# Redis
if command -v redis-cli &> /dev/null; then
    if redis-cli ping &> /dev/null; then
        echo "✅ Redis 连接正常"
    else
        echo "⚠️  Redis 连接失败"
    fi
else
    echo "⚠️  Redis 未安装 (可选)"
fi

echo "🎉 环境检查完成！"
```

运行验证：
```bash
chmod +x check-env.sh
./check-env.sh
```

## 🚨 常见问题

### Go 环境问题

**问题**: `go: command not found`
**解决**: 检查 `PATH` 环境变量是否包含 Go 的 bin 目录

```bash
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

**问题**: `GOPATH` 和 `GOROOT` 配置
**解决**: Go 1.11+ 使用模块模式，通常不需要设置 `GOPATH`

### 数据库连接问题

**问题**: MySQL 连接被拒绝
**解决**: 检查服务状态和防火墙设置

```bash
# 检查服务状态
sudo systemctl status mysql

# 检查端口监听
netstat -tlnp | grep 3306

# 重启服务
sudo systemctl restart mysql
```

### 权限问题

**问题**: Docker 权限不足
**解决**: 将用户添加到 docker 组

```bash
sudo usermod -aG docker $USER
# 重新登录或执行
newgrp docker
```

---

**下一步**: 查看 [安装指南](./installation.md) 开始安装项目