## Skill 条目模板

使用本模板快速添加新 skill 到 skills-map。

### 表格条目格式

```markdown
| **[skill-name](file:///path/to/skill/SKILL.md)** | 简短描述 | 适用场景 | 附加信息 |
```

### Mermaid 图表节点格式

```mermaid
skill-name[skill-name<br/>简短描述]
```

### 完整示例

#### 添加到表格

**开发类 Skills**：

```markdown
| **[grpc-service](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/grpc-service/SKILL.md)** | 创建 gRPC 服务 | gRPC 服务端开发 | RPC |
```

**管理类 Skills**：

```markdown
| **[deployment](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/deployment/SKILL.md)** | 部署指南 | 应用部署流程 | 可选 |
```

#### 添加到 Mermaid 图表

```mermaid
grpc[grpc-service<br/>gRPC 服务]
```

### 变更历史条目格式

```markdown
### YYYY-MM-DD

- ✅ 新增 `skill-name` - 简短描述
- 🔄 更新 `skill-name` - 更新说明
- ❌ 删除 `skill-name` - 删除原因
```

### 快速索引条目格式

```markdown
| 我想要... | 使用这个 Skill |
| --------- | -------------- |
| 场景描述  | skill-name     |
```
