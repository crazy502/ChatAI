# ChatAI / AgentGo

一个基于 Go + Vue 3 的智能对话系统，包含用户注册登录、邮件验证码、会话管理、多模型聊天、流式响应和管理员监控面板。

当前仓库名为 `ChatAI`，前后端页面与系统名称主要使用 `AgentGo`。

## 项目概览

- 前端：Vue 3 + Vue Router + Axios
- 后端：Go 1.24 + Gin + Gorm
- 数据库：MySQL
- 缓存：Redis
- 队列：RabbitMQ
- AI 模型：Qwen、DeepSeek
- 邮件：QQ SMTP

## 当前已实现功能

- 登录 / 注册同页展示，注册与登录在同一张卡片内切换
- 邮箱验证码注册，验证码存入 Redis，并通过邮件发送
- 用户登录后签发 JWT
- 首次启动时自动检查并创建管理员账号
- AI 对话支持 `qwen` 与 `deepseek` 两种模型
- 支持普通对话与流式对话（SSE）
- 支持会话创建、会话历史、会话重命名、置顶、归档、搜索
- 支持管理员监控页，查看请求量、错误率、平均延迟、接口与模型健康状态
- 支持请求指标采集与模型调用指标采集
- RabbitMQ 不可用时，消息持久化会降级为直接写入 MySQL，减少系统不可用风险

## 真实目录结构

当前项目真正参与运行的目录主要是 `client` 和 `server`：

```text
F:\ChatAI
├── client/                  # Vue 3 前端
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── composables/
│   │   ├── router/
│   │   ├── styles/
│   │   ├── utils/
│   │   └── views/
│   ├── package.json
│   └── vue.config.js
├── server/                  # Go 后端
│   ├── config/
│   │   └── config.example.toml
│   ├── infra/
│   │   ├── cache/
│   │   ├── config/
│   │   ├── db/
│   │   ├── mail/
│   │   ├── metrics/
│   │   └── mq/
│   ├── internal/
│   │   ├── admin/
│   │   ├── ai/
│   │   ├── chat/
│   │   ├── middleware/
│   │   ├── router/
│   │   ├── session/
│   │   └── user/
│   ├── pkg/
│   │   ├── code/
│   │   ├── jwt/
│   │   ├── password/
│   │   ├── response/
│   │   └── utils/
│   ├── go.mod
│   └── main.go
├── GPTest/                  # 辅助目录，不参与当前前后端运行
└── GraduationThesis/        # 论文相关目录
```

## 页面说明

- `/login`：登录页，内部集成注册表单和聊天预览
- `/register`：已重定向到 `/login`
- `/menu`：系统菜单 / 控制台
- `/ai-chat`：智能对话页面
- `/admin-metrics`：管理员监控面板

## 后端模块说明

- `internal/user`：登录、注册、验证码、管理员初始化
- `internal/session`：会话列表、重命名、置顶、归档
- `internal/chat`：消息发送、流式响应、历史记录、消息入库
- `internal/admin`：监控接口聚合
- `internal/ai`：模型提供者、上下文管理、流式生成
- `internal/middleware`：鉴权、管理员权限、监控、恢复中间件
- `infra/db`：MySQL 初始化与自动迁移
- `infra/cache`：Redis 初始化与验证码缓存
- `infra/mq`：RabbitMQ 初始化与消息队列
- `infra/mail`：邮件发送
- `infra/metrics`：请求与模型指标采集

## 环境要求

建议环境如下：

- Go `1.24+`
- Node.js `18+`
- npm `9+`
- MySQL `8+`
- Redis `6+`
- RabbitMQ `3+`（推荐，可缺省启动为降级模式）

## 快速开始

### 1. 准备后端配置

当前后端通过相对路径读取 `config/config.toml`，所以后端命令建议在 `server` 目录执行。

在 PowerShell 中执行：

```powershell
cd F:\ChatAI\server
Copy-Item .\config\config.example.toml .\config\config.toml
```

然后根据本机环境修改 `F:\ChatAI\server\config\config.toml`。

### 2. 启动后端

```powershell
cd F:\ChatAI\server
go mod download
go run .
```

说明：

- 当前入口文件是 `F:\ChatAI\server\main.go`
- 也可以使用 `go run main.go`
- 但请在 `server` 目录下运行，否则会因为配置文件相对路径导致读取失败

后端默认监听：

```text
http://localhost:9090
```

### 3. 启动前端

```powershell
cd F:\ChatAI\client
npm install
npm run serve
```

前端默认访问地址：

```text
http://localhost:8080
```

### 4. 前后端联调说明

前端开发服务器已配置代理：

- 前端请求 `/api`
- 自动代理到 `http://localhost:9090/api/v1`

也就是说，前端代码里写的是：

```text
/api/user/login
```

最终会转发到后端：

```text
/api/v1/user/login
```

## 配置文件说明

`server/config/config.toml` 主要包含以下配置块：

- `mainConfig`：服务名、监听地址、端口
- `emailConfig`：QQ 邮箱账号与授权码
- `redisConfig`：Redis 连接配置
- `mysqlConfig`：MySQL 连接配置
- `jwtConfig`：JWT 过期时间、签发者、密钥
- `adminConfig`：管理员账号初始化信息
- `rabbitmqConfig`：RabbitMQ 连接配置
- `qwenConfig`：通义千问模型配置
- `deepseekConfig`：DeepSeek 模型配置

建议至少先正确配置以下内容：

- MySQL
- Redis
- 邮件账号与授权码
- JWT 密钥
- 至少一个 AI 模型的 `apiKey`

## 数据初始化说明

系统启动时会自动完成以下动作：

- 初始化 MySQL 连接
- 自动迁移数据表：`users`、`sessions`、`messages`
- 检查并创建管理员账号
- 初始化 Redis
- 初始化 RabbitMQ
- 启动消息消费者

其中 RabbitMQ 初始化失败时，系统会以降级模式继续运行；但 Redis 初始化失败会直接导致后端启动失败，因为注册验证码依赖 Redis。

## 默认管理员说明

管理员账号会在启动时根据 `adminConfig` 自动创建或校正。

如果配置为空，代码中的回退值为：

- 用户名：`admin`
- 密码：`admin`
- 邮箱：`admin@chatai.local`

实际使用时请务必在 `config.toml` 中改成强密码。

## API 概览

### 用户模块

- `POST /api/v1/user/captcha`：发送邮箱验证码
- `POST /api/v1/user/register`：注册
- `POST /api/v1/user/login`：登录

### 对话与会话模块

- `GET /api/v1/AI/chat/sessions`：获取当前用户会话列表
- `POST /api/v1/AI/chat/session/rename`：重命名会话
- `POST /api/v1/AI/chat/session/pin`：置顶 / 取消置顶
- `POST /api/v1/AI/chat/session/archive`：归档 / 恢复归档
- `POST /api/v1/AI/chat/send-new-session`：创建新会话并发送消息
- `POST /api/v1/AI/chat/send`：向已有会话发送消息
- `POST /api/v1/AI/chat/history`：获取历史消息
- `POST /api/v1/AI/chat/send-stream-new-session`：创建新会话并使用流式响应
- `POST /api/v1/AI/chat/send-stream`：已有会话流式响应

### 管理模块

- `GET /api/v1/admin/metrics/all`：获取全部监控快照

## 认证与权限

- 登录成功后后端返回 JWT
- 前端会将 token 存入 `localStorage`
- `/api/v1/AI/*` 路由需要登录
- `/api/v1/admin/*` 路由需要管理员权限

当前管理员权限判定逻辑是：

- token 中 `is_admin = true`
- 当前用户名与 `adminConfig.username` 一致

## 聊天能力说明

当前聊天模块支持：

- 模型切换：`qwen` / `deepseek`
- 普通回复
- SSE 流式回复
- 历史消息回填
- 会话标题自动生成
- 消息持久化

另外，后端在生成回复时会记录模型调用指标，包括：

- 请求次数
- 错误次数
- 错误率
- 平均延迟
- 最近成功时间
- 最近失败时间

## 监控能力说明

管理员监控页当前聚合了以下指标：

- 全局请求总数
- 全局错误总数
- 错误率
- 平均延迟
- 路由维度统计
- 模型维度统计
- 时间序列归档趋势

当前归档策略为：

- 约每 `30` 秒采样一次
- 默认保留近 `6` 小时数据

## 构建命令

### 前端构建

```powershell
cd F:\ChatAI\client
npm run build
```

### 后端编译

```powershell
cd F:\ChatAI\server
go build ./...
```

## 常见问题

### 1. 为什么后端不能直接在仓库根目录执行 `go run server/main.go`？

因为当前配置文件读取使用的是相对路径：

```text
config/config.toml
```

这个路径是相对于运行时工作目录解析的。当前项目应在 `F:\ChatAI\server` 目录下执行后端命令：

```powershell
cd F:\ChatAI\server
go run .
```

### 2. RabbitMQ 没启动还能聊天吗？

可以。当前代码在 RabbitMQ 不可用时会进入降级模式，消息会直接写 MySQL。

### 3. Redis 没启动可以注册吗？

不可以。验证码发送和校验依赖 Redis。

### 4. 注册后为什么提示“使用邮件中的账号登录”？

当前注册逻辑会为用户自动生成用户名，并通过邮件发送给注册邮箱。登录使用的是用户名，不是邮箱。

## 当前前端主题状态

当前前端已经统一为浅色主题基线：

- 页面背景：浅灰白
- 面板：白色
- 文本：冷黑 / 冷灰
- 强调色：绿色，仅用于按钮和高亮
- 危险色：红色

主题样式主要集中在：

- `F:\ChatAI\client\src\styles\sci-fi-theme.css`

## 开发建议

- 后端继续保持 `internal` 按业务域拆分
- README 中的运行命令请优先以 `server` / `client` 子目录为工作目录
- 如果后续要支持部署，建议把后端配置读取改成支持环境变量或绝对路径
- 如果后续要写论文，可以重点描述“流式对话 + 指标监控 + 队列降级持久化”这三部分

## 仓库说明

本 README 基于当前真实代码目录与已实现功能生成，不是通用模板文档。后续如果你继续调整目录、接口或部署方式，建议同步更新这里的启动说明和接口列表。
