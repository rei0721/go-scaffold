# 贡献指南

感谢您对 Go Scaffold 项目的关注！我们欢迎所有形式的贡献，包括但不限于代码贡献、文档改进、问题报告和功能建议�?

## 🤝 贡献方式

### 1. 代码贡献
- 修复 Bug
- 添加新功�?
- 性能优化
- 代码重构

### 2. 文档贡献
- 改进现有文档
- 添加使用示例
- 翻译文档
- 修正错别�?

### 3. 问题反馈
- 报告 Bug
- 提出功能请求
- 改进建议
- 使用体验反馈

### 4. 社区参与
- 回答问题
- 代码审查
- 分享使用经验
- 推广项目

## 🚀 快速开�?

### 1. Fork 项目

点击 GitHub 页面右上角的 "Fork" 按钮，将项目 Fork 到您的账户下�?

### 2. 克隆代码

```bash
# 克隆�?Fork 的仓�?
git clone https://github.com/YOUR_USERNAME/go-scaffold.git
cd go-scaffold

# 添加上游仓库
git remote add upstream https://github.com/rei0721/go-scaffold.git
```

### 3. 创建分支

```bash
# �?main 分支创建新的功能分支
git checkout -b feature/your-feature-name

# 或者修复分�?
git checkout -b fix/your-fix-name
```

### 4. 开发环境设�?

```bash
# 安装依赖
go mod download

# 复制配置文件
cp configs/config.example.yaml configs/config.yaml
cp .env.example .env

# 初始化数据库
go run ./cmd/server initdb

# 运行测试确保环境正常
make test
```

## 📝 开发流�?

### 1. 开发前准备

#### 同步上游代码
```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

#### 创建功能分支
```bash
git checkout -b feature/add-user-management
```

### 2. 开发过�?

#### 编写代码
- 遵循项目的代码规�?
- 添加必要的注释和文档
- 确保代码的可读性和可维护�?

#### 编写测试
```bash
# 为新功能编写单元测试
go test ./internal/service/user -v

# 运行所有测�?
make test

# 检查测试覆盖率
make test-coverage
```

#### 代码检�?
```bash
# 代码格式�?
make fmt

# 代码检�?
make lint

# 安全检�?
make security
```

### 3. 提交代码

#### 提交信息规范
遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范�?

```bash
# 功能添加
git commit -m "feat: add user management API"

# Bug 修复
git commit -m "fix: resolve database connection timeout issue"

# 文档更新
git commit -m "docs: update API documentation"

# 代码重构
git commit -m "refactor: improve error handling in auth service"

# 性能优化
git commit -m "perf: optimize database query performance"

# 测试相关
git commit -m "test: add unit tests for user service"
```

#### 推送代�?
```bash
git push origin feature/add-user-management
```

### 4. 创建 Pull Request

1. 访问您的 Fork 仓库页面
2. 点击 "New Pull Request" 按钮
3. 选择目标分支（通常�?`main`�?
4. 填写 PR 标题和描�?
5. 提交 Pull Request

## 📋 Pull Request 模板

创建 PR 时，请使用以下模板：

```markdown
## 变更类型
- [ ] Bug 修复
- [ ] 新功�?
- [ ] 代码重构
- [ ] 性能优化
- [ ] 文档更新
- [ ] 测试改进

## 变更描述
简要描述此次变更的内容和目的�?

## 相关 Issue
关闭 #issue_number

## 测试
- [ ] 已添加单元测�?
- [ ] 已添加集成测�?
- [ ] 已运行现有测试套�?
- [ ] 已手动测试功�?

## 检查清�?
- [ ] 代码遵循项目规范
- [ ] 已添加必要的文档
- [ ] 已更新相关的 README 或文�?
- [ ] 变更向后兼容
- [ ] 已添加适当的日志记�?

## 截图（如适用�?
如果�?UI 相关的变更，请提供截图�?

## 其他说明
任何其他需要说明的内容�?
```

## 🔍 代码审查

### 审查标准

1. **功能�?*
   - 代码是否实现了预期功�?
   - 是否存在逻辑错误
   - 边界条件是否处理正确

2. **代码质量**
   - 代码是否清晰易读
   - 是否遵循项目规范
   - 是否有适当的注�?

3. **性能**
   - 是否存在性能问题
   - 资源使用是否合理
   - 是否有内存泄漏风�?

4. **安全�?*
   - 是否存在安全漏洞
   - 输入验证是否充分
   - 敏感信息是否正确处理

5. **测试**
   - 测试覆盖率是否充�?
   - 测试用例是否合理
   - 是否包含边界测试

### 审查流程

1. **自动检�?*
   - CI/CD 流水线自动运�?
   - 代码格式检�?
   - 单元测试执行
   - 安全扫描

2. **人工审查**
   - 至少需要一个维护者审�?
   - 复杂变更需要多人审�?
   - 审查者提供建设性反�?

3. **修改和完�?*
   - 根据反馈修改代码
   - 回复审查意见
   - 重新请求审查

## 🐛 问题报告

### Bug 报告模板

```markdown
## Bug 描述
清晰简洁地描述 Bug�?

## 复现步骤
1. 执行 '...'
2. 点击 '....'
3. 滚动�?'....'
4. 看到错误

## 预期行为
描述您期望发生的行为�?

## 实际行为
描述实际发生的行为�?

## 截图
如果适用，添加截图来帮助解释问题�?

## 环境信息
- OS: [例如 macOS 12.0]
- Go 版本: [例如 1.24.6]
- 项目版本: [例如 v0.1.2]

## 额外信息
添加任何其他有关问题的信息�?
```

### 功能请求模板

```markdown
## 功能描述
清晰简洁地描述您想要的功能�?

## 问题背景
描述这个功能要解决的问题�?

## 解决方案
描述您希望的解决方案�?

## 替代方案
描述您考虑过的替代解决方案�?

## 额外信息
添加任何其他有关功能请求的信息�?
```

## 📚 文档贡献

### 文档类型

1. **API 文档**
   - 接口说明
   - 参数描述
   - 响应示例
   - 错误码说�?

2. **使用指南**
   - 快速开�?
   - 配置说明
   - 最佳实�?
   - 故障排除

3. **开发文�?*
   - 架构设计
   - 代码规范
   - 测试指南
   - 部署说明

### 文档规范

1. **Markdown 格式**
   - 使用标准 Markdown 语法
   - 适当使用标题层级
   - 添加目录导航

2. **代码示例**
   - 提供完整可运行的示例
   - 添加必要的注�?
   - 使用正确的语法高�?

3. **图片和图�?*
   - 使用 Mermaid 绘制流程�?
   - 图片使用相对路径
   - 提供图片的替代文�?

## 🏷�?版本发布

### 版本号规�?

采用 [语义化版本](https://semver.org/lang/zh-CN/) 规范�?

- **主版本号**：不兼容�?API 修改
- **次版本号**：向下兼容的功能性新�?
- **修订�?*：向下兼容的问题修正

### 发布流程

1. **准备发布**
   - 更新版本�?
   - 更新 CHANGELOG
   - 运行完整测试套件

2. **创建发布**
   - 创建 Git 标签
   - 生成发布说明
   - 发布�?GitHub Releases

3. **发布�?*
   - 更新文档
   - 通知社区
   - 监控反馈

## 🎯 贡献者权�?

### 认可方式

1. **贡献者列�?*
   - README 中的贡献者部�?
   - 项目网站的贡献者页�?

2. **特殊徽章**
   - 核心贡献者徽�?
   - 文档贡献者徽�?
   - 社区贡献者徽�?

3. **参与决策**
   - 重要功能的讨�?
   - 项目方向的决�?
   - 技术选型的参�?

### 成为维护�?

满足以下条件的贡献者可以申请成为维护者：

1. **持续贡献**
   - 至少 3 个月的活跃贡�?
   - 提交了多个高质量�?PR
   - 参与了代码审�?

2. **技术能�?*
   - 熟悉项目架构和代�?
   - 具备良好�?Go 编程能力
   - 理解项目的设计理�?

3. **社区参与**
   - 积极回答社区问题
   - 帮助其他贡献�?
   - 参与项目讨论

## 📞 联系方式

### 获取帮助

1. **GitHub Issues**
   - 报告问题
   - 功能请求
   - 技术讨�?

2. **GitHub Discussions**
   - 一般性讨�?
   - 使用经验分享
   - 社区交流

3. **邮件联系**
   - 安全问题报告
   - 私人咨询
   - 合作洽谈

### 社区准则

1. **友善和尊�?*
   - 尊重不同观点
   - 友善地提供反�?
   - 避免人身攻击

2. **建设性沟�?*
   - 提供具体的建�?
   - 解释问题的背�?
   - 分享解决方案

3. **包容�?*
   - 欢迎新手参与
   - 帮助他人成长
   - 创造包容的环境

## 🙏 致谢

感谢所有为 Go Scaffold 项目做出贡献的开发者！您的每一个贡献都让这个项目变得更好�?

特别感谢�?
- 核心维护者团�?
- 活跃的贡献�?
- 文档贡献�?
- 问题报告�?
- 社区支持�?

---

**让我们一起构建更好的 Go Scaffold�?* 🚀
