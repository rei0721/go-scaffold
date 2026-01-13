# 配置说明

> 完整配置参数详解

---

## Purpose

说明所有配置项的含义、默认值和推荐设置。

## Scope

`configs/config.yaml` 和 `.env` 环境变量。

## Status

**Active**

---

## 配置优先级

```
环境变量 > .env 文件 > config.yaml 默认值
```

配置文件支持 `${VAR:default}` 语法，优先使用环境变量，无则使用默认值。

---

## Server 配置

| 参数                  | 环境变量      | 默认值    | 说明                           |
| --------------------- | ------------- | --------- | ------------------------------ |
| `server.host`         | -             | `0.0.0.0` | 监听地址                       |
| `server.port`         | -             | `9999`    | 监听端口                       |
| `server.mode`         | `SERVER_MODE` | `debug`   | 运行模式（debug/release/test） |
| `server.readTimeout`  | -             | `10`      | 读取超时（秒）                 |
| `server.writeTimeout` | -             | `10`      | 写入超时（秒）                 |

---

## Database 配置

| 参数                    | 环境变量      | 默认值      | 说明                          |
| ----------------------- | ------------- | ----------- | ----------------------------- |
| `database.driver`       | `DB_DRIVER`   | `postgres`  | 驱动（postgres/mysql/sqlite） |
| `database.host`         | `DB_HOST`     | `localhost` | 主机地址                      |
| `database.port`         | `DB_PORT`     | `5432`      | 端口                          |
| `database.user`         | `DB_USER`     | `postgres`  | 用户名                        |
| `database.password`     | `DB_PASSWORD` | -           | 密码 ⚠️ 生产必须设置          |
| `database.dbname`       | `DB_NAME`     | `rei0721`   | 数据库名                      |
| `database.maxOpenConns` | -             | `100`       | 最大打开连接数                |
| `database.maxIdleConns` | -             | `10`        | 最大空闲连接数                |

---

## Redis 配置

| 参数                 | 环境变量         | 默认值      | 说明               |
| -------------------- | ---------------- | ----------- | ------------------ |
| `redis.enabled`      | `REDIS_ENABLED`  | `true`      | 是否启用           |
| `redis.host`         | `REDIS_HOST`     | `localhost` | 主机地址           |
| `redis.port`         | `REDIS_PORT`     | `6379`      | 端口               |
| `redis.password`     | `REDIS_PASSWORD` | -           | 密码               |
| `redis.db`           | `REDIS_DB`       | `0`         | 数据库索引（0-15） |
| `redis.poolSize`     | -                | `30`        | 连接池大小         |
| `redis.minIdleConns` | -                | `10`        | 最小空闲连接数     |
| `redis.maxRetries`   | -                | `3`         | 最大重试次数       |
| `redis.dialTimeout`  | -                | `5`         | 连接超时（秒）     |
| `redis.readTimeout`  | -                | `3`         | 读取超时（秒）     |
| `redis.writeTimeout` | -                | `3`         | 写入超时（秒）     |

---

## Logger 配置

| 参数                    | 环境变量     | 默认值         | 说明                              |
| ----------------------- | ------------ | -------------- | --------------------------------- |
| `logger.level`          | `LOG_LEVEL`  | `debug`        | 日志级别（debug/info/warn/error） |
| `logger.format`         | `LOG_FORMAT` | `console`      | 默认格式（json/console）          |
| `logger.console_format` | -            | -              | 控制台专用格式                    |
| `logger.file_format`    | -            | `json`         | 文件专用格式                      |
| `logger.output`         | `LOG_OUTPUT` | `both`         | 输出目标（stdout/file/both）      |
| `logger.file_path`      | -            | `logs/app.log` | 日志文件路径                      |
| `logger.max_size`       | -            | `10`           | 单文件最大大小（MB）              |
| `logger.max_backups`    | -            | `10`           | 保留文件数量                      |
| `logger.max_age`        | -            | `30`           | 保留天数                          |

---

## I18n 配置

| 参数             | 环境变量       | 默认值         | 说明           |
| ---------------- | -------------- | -------------- | -------------- |
| `i18n.default`   | `I18N_DEFAULT` | `zh-CN`        | 默认语言       |
| `i18n.supported` | -              | `zh-CN, en-US` | 支持的语言列表 |

---

## Executor 配置

| 参数                           | 默认值  | 说明                                 |
| ------------------------------ | ------- | ------------------------------------ |
| `executor.enabled`             | `true`  | 是否启用执行器                       |
| `executor.pools[]`             | -       | 池配置列表（数组）                   |
| `executor.pools[].name`        | -       | 池名称（http/database/background）   |
| `executor.pools[].size`        | -       | 池容量（最大并发 worker 数）         |
| `executor.pools[].expiry`      | -       | worker 过期时间（秒）                |
| `executor.pools[].nonBlocking` | `false` | 非阻塞模式（池满时是否立即返回错误） |

**示例配置**：

```yaml
executor:
  enabled: true
  pools:
    - name: http
      size: 200
      expiry: 10
      nonBlocking: true
    - name: database
      size: 50
      expiry: 30
      nonBlocking: false
    - name: background
      size: 30
      expiry: 60
      nonBlocking: true
```

---

## 环境配置建议

### 开发环境

```yaml
server:
  mode: debug
logger:
  level: debug
  output: both
  format: console
```

### 生产环境

```yaml
server:
  mode: release
logger:
  level: info
  output: file
  format: json
```

---

## Evidence

- 配置模板：[configs/config.example.yaml](../configs/config.example.yaml)
- 环境变量模板：[.env.example](../.env.example)
- 配置加载：[internal/config/](../internal/config/)

## Related

- [快速开始](quickstart.md)
- [开发指南](development.md)

## Changelog

| 日期       | 变更                   |
| ---------- | ---------------------- |
| 2026-01-14 | 初始创建               |
| 2026-01-14 | 补充 Executor 配置章节 |
