# AI 应用服务平台

现代 Web 应用开发，通过集成 AI 技术和前沿功能（AI 聊天助手和图像识别）

项目采用全栈结构

后端
- 使用 Gin 框架构建 RestfulAPI


前端
- 采用 Vue. js 实现响应式界面

数据库
- 集成 MySQL 和 redis

消息队列
- RabbitMQ

图像识别模型（需进行改动，Windos 不行）
- 通过 ONNX Runtime 运行 MobileNetV 2

实现端到端的智能化应用开发体验


项目通过集成 AI 助手（支持多会话管理和流式输出）、图像分类识别等前沿技术，全栈应用开发。


学习
- 如何在 Go 中处理 JWT 认证
- 异步消息存储
- SSE 流式响应
- 与外部 AI 模型集成

使用流行开源框架开发

目前已有模块
- 用户管理
- 会话管理
- AI 对话
- 图像处理模块
- 后续可进行进一步的模块化设计如将 AI 助手抽象为可复用的 Helper 类，支持不太模型类型的扩展
---


# 1：开篇

本项目旨利用Gin框架（强大Web服务能力）和EINO框架（开源且强大的AI服务能力），简化开发流程，实现AI应用的快速迭代。

## 项目目标

本项目旨在帮助学习Go语言的开发者快速上手现代Web应用开发，通过集成AI技术和前沿功能（如AI聊天助手和图像识别），避免单纯的语法学习陷阱，让开发者从实际项目中掌握Go语言的工程化应用。

项目采用全栈架构
- 后端使用Gin框架构建RESTful API
- 前端采用Vue.js实现响应式界面
- 数据库层集成MySQL和Redis
- 消息队列使用RabbitMQ
- 图像识别模型通过ONNX Runtime运行MobileNetV2，实现端到端的智能化应用开发体验。

项目帮助Go初学者拓宽视野，不局限于基础CRUD操作，而是通过集成AI助手（支持多会话管理和流式输出）、图像分类识别等前沿技术，进行全栈应用开发。

开发者将学习如何在Go中处理JWT认证、异步消息存储、SSE流式响应，以及与外部AI模型的集成，打破传统教程的局限性，培养解决实际问题的能力。

此外，项目让Go开发者掌握使用流行开源框架进行高效开发，避免每次从零构建基础设施。

通过完整的项目模板，包括用户管理、会话管理、AI对话和图像处理模块，开发者可以在紧急项目中节省时间，提升生产力。

项目鼓励模块化设计，如将AI助手抽象为可复用的Helper类，支持不同模型类型的扩展，为后续开发提供坚实基础。

最终，通过这个项目，开发者不仅能构建一个功能完整的AI聊天和图像识别应用，还能深入理解Go在微服务、并发处理和AI集成方面的最佳实践，为职业发展奠定基础。项目名称：GopherAI （AI应用服务平台） 


基于本项目，开发的一个应用：GopherAI智能助手平台

### 模块大致如下 ：

- 登录注册（基于JWT）：用户通过邮箱注册，系统生成11位随机账号并发送至邮箱。注册时需验证邮箱验证码（Redis存储，2分钟有效）。登录使用账号+密码，成功后返回JWT Token（包含用户ID和用户名，有效期可配置）。中间件验证Token，支持Authorization头或URL参数传Token，验证失败返回错误。
	
- AI聊天对话（包含流式输出与全量输出的选择，多会话管理）：支持多会话，每个会话独立管理消息历史。用户可创建新会话或在现有会话中继续对话。提供同步和流式两种输出模式：同步返回完整AI回复，流式通过SSE实时推送内容片段（以[DONE]结束）。消息通过RabbitMQ异步存储到数据库。支持获取用户所有会话列表和特定会话的历史记录。AI模型类型可配置（如GPT等），使用上下文管理器管理不同用户的AI助手实例。
	
- 图像识别：上传图片文件（支持JPEG/PNG/GIF），使用MobileNetV2 ONNX模型进行分类识别。模型输入224x224，输出1000类ImageNet标签中的最高概率类别。识别过程包括图片解码、缩放、归一化、ONNX推理，返回识别的类别名称。



# 2：本项目需要的基础

完成本项目需要多久不同基础的同学完成本项目周期会有明显差异，下面给出大体列一下：

|   |   |   |   |
|---|---|---|---|
|学员基础类型|典型背景|项目完成周期|完成方式说明|
|具备 Go 基础 + 做过 Web 开发|后端开发/转语言同学|3~7 天|直接跟代码，可快速掌握 AI 模型集成、流式传输、会话管理|
|有后端经验但不熟 Go|Java / Python / Node / C++ 后端|7~14 天|先用 2~4 天熟悉 Go 语法 + Gin，再进入工程部分|
|只有算法或刷题基础，但无后端经验|刷过 LeetCode / 有编程基础的学生|14~21 天|从 HTTP/REST、数据库、鉴权体系逐步带入，完成节奏可控|
|完全从零开始 Web 开发|没做过项目的新手|21-35 天|按课程安排由浅入深，重点在 Web 服务、数据库、Redis 与模型集成|

## 开发本项目需要的技术基础 （Go 版 AI 应用服务平台）

本项目面向有一定编程基础、希望系统掌握 AI 应用后端开发能力的同学。

项目从可运行 Demo 到可用于生产环境的工程化体系，循序渐进构建，因此无需一开始就完全掌握所有内容。

### 1 . Go 语言基础（必要）

需要达到能够阅读与编写基础业务逻辑的水平，包括：

|   |   |
|---|---|
|知识点|要求说明|
|Go 基础语法|变量、结构体、接口、切片、map、错误处理|
|Go 并发模型|goroutine、context 使用场景及常用模式|
|网络与 Web 基础|HTTP 请求处理、JSON 序列化、RESTful API 基础|
|模块化开发|会使用 go mod，理解项目结构划分|

不要求掌握高深的 Go 底层原理，但需要能读懂代码并能调试。

### 2 . Web 后端开发基础（必要）

项目提供完整的后端 API 服务，需要对 Web 服务开发有基本理解。

| 组件     | 用法                  |
| ------ | ------------------- |
| Gin 框架 | 路由、中间件、请求响应处理       |
| JWT 鉴权 | 登录态维护和访问权限控制        |
| 配置管理   | 可基于 toml / env 进行配置 |

如你做过任何语言的后端（例如 Java、Node、Python），则可快速迁移入门。

### 3 . 数据存储与缓存（建议掌握）本项目使用：

|   |   |   |
|---|---|---|
|技术|用途|要求程度|
|MySQL|用户数据、会话记录存储|会 CRUD + 了解事务与索引|
|Redis|缓存验证码等|会 set/get 以及 TTL、连接池|

ORM 使用 GORM，可边用边学，不要求提前精通。

### 4 . AI 模型调用基础（有就行，不需提前懂）项目会集成：

|   |   |
|---|---|
|模型来源|对接方式|
|OpenAI|REST API 调用|
|本地模型 Ollama|HTTP 调用 + 流式输出|

只需要理解：什么是 Token、上下文、Prompt同步请求和 流式输出 的区别模型切换和配置不会 AI 原理也完全没关系，项目会从 API 调用层教学。

### 5 . 前端基础（可选，会更好）

本项目前端采用 Vue3 + Vite。

|   |   |
|---|---|
|要求|程度|
|能理解组件结构、路由、API 请求|推荐|
|如果完全不会|也可直接用本项目提供的前端代码，不影响学习后端|

### 6 . 消息队列与异步任务（加分项，不会也可学）

使用 RabbitMQ 做异步消息处理。如果你之前没有接触过消息队列，本项目会带你从 0 使用，不需要提前掌握。

### 学完后将具备的能力

通过本项目，你将：
1. 掌握 Go 后端完整工程化能力（配置、路由、中间件、日志、错误处理）
2. 会从 0 集成 AI 模型（包括云端模型和本地模型）
3. 理解 AI 流式响应 + 并发处理的设计模式
4. 能独立开发一个可商用的 AI Web 服务
5. 具备向 C++ / Java / Python / Rust 等语言迁移的抽象能力

---

# 3：项目介绍

# 项目概述

该项目是一个使用 Gin 框架构建的 Go 语言 Web 应用服务平台，Gin 是一个高性能的 HTTP Web 框架。

该平台旨在高效处理 HTTP 请求和响应，集成 AI 助手、图像识别等功能，为构建现代 Web 应用程序提供基础。

项目视频演示注意：这里的验证码还有账号，是会发送到你的个人邮箱当中的，这边视频填入了当时发送的信息

# 技术架构图

![[技术架构.png]]

这张架构图展示了 GopherAI 整个系统的核心组成部分，包括：
- 业务服务
- AI 推理
- 第三方平台
- 基础设施
- 消息队列
- 数据库
- 以及前后端交互流程。

你可以把它理解为： “用户从输入一个问题 → 后端处理 → AI 推理 → 数据落库 → 前端实时显示” 的全链路流程图。

下面我们从外往内、从上到下解析整个系统。

---

## 1. 用户请求入口：前端 UI / 客户端

左侧的电脑图标代表 **用户客户端**（浏览器 / 前端应用）。

用户通过浏览器进行：

- 登录 / 注册
- 发送聊天消息
- 上传图片进行识别
- 获取 AI 回复（流式输出）

这些请求都会发送到右侧的 **业务服务**。

## 2. 第三方平台（大模型提供方）

图右上角的“阿里云”与“QQ”代表：

- 通义千问（阿里云）
- QQ邮箱服务
- 其他你以后可以扩展的 LLM

这里的作用是：

👉 **GopherAI 需要从这里获取 API-Key 和模型调用能力。**

只要配置好 API-Key，就可以直接把这些模型接入你的服务中。

## 3 . 业务服务（整个系统的大脑）

业务服务区块分为三大服务：

### **① AI 聊天服务（LLM Chat Service）**

负责：

- 多轮对话
- 流式输出（SSE 实时推送）
- 调用 LLM 生成回复
- 将消息写入 RabbitMQ 队列
- AI 会话历史管理

这是 GopherAI 的核心业务之一。

---

### **② 用户服务（User Service）**

负责：

- 用户注册
- 邮箱验证码
- JWT 登录鉴权
- Token 校验
- 用户信息查询

提供身份管理能力。

---

### **③ AI 图像识别服务（Image Recognition Service）**

负责：

- 接收图片上传
- 调用 ONNXRuntime 执行推理
- 返回图像分类结果

这是你的第二大核心业务。

---

## 4. 业务服务所使用的基础工具（绿色模块）

图中绿色框表示你项目中使用的 Go 生态工具：

### **Gin**

👉 Go 最流行的 Web 框架，用于处理 HTTP 请求。

### **Gorm**

👉 ORM 框架，用于操作数据库（MySQL）。

### **go-redis / redis/v8**

👉 操作 Redis（验证码、缓存、session）。

### **AMQP**

👉 RabbitMQ AMQP 协议客户端，用于队列通信。

### **SSE（Server-Sent Events）**

👉 用于 AI 聊天的流式响应。

### **http**

👉 底层网络交互组件。

你可以看到，整个业务模块都是围绕这些组件构建的。

---

## 5. 消息队列（RabbitMQ）

- AI 消息写入队列
- 后台异步消费消息
- 避免阻塞主线程
- 提高系统吞吐量

RabbitMQ 在架构中承担 **解耦 + 削峰 + 异步处理** 的角色。

这是你项目的一个亮点，也是面试官非常喜欢问的点。

---

## 6. 数据库层（MySQL + Redis）

### **MySQL：核心业务数据**

- 用户表
- 会话表
- 会话消息表
- 图片识别记录
- ……

所有结构化数据都持久化在 MySQL 中。

---

### **Redis：缓存与验证码**

- 邮箱验证码
- 用户 Token
- 可做热点数据缓存
- 存储临时会话数据（选做）

Redis 用于高性能、短时效的快速读写场景。

---

## 7. 推理引擎（ONNX Runtime）

右侧 ONNX 模块：

- 使用 MobileNetV2 模型
- 部署在 onnxruntime-linux-x64 推理引擎中
- 完成图像分类任务

这是系统中的 **AI 视觉能力模块**。

---

## 8. Docker 部署（容器化环境）

底部的 Docker 区域表示：

你可以用 Docker Compose 一键启动：

- 后端服务
- Redis
- MySQL
- RabbitMQ
- ONNXRuntime

实现完整的可运行环境。

这让你的项目具备 “真正能上线” 的能力。

---
# 关键组件

## 1. **Router 和 Controller**

- **Router 模块**（如 `router/user.go`、`router/AI.go`、`router/Image.go`）  
  - 使用 **Gin 框架**定义 API 路由分组  
  - 支持 **RESTful 设计**  
  - 后端路由包括：
    - 用户认证：`/user/login`、`/user/register`
    - AI 聊天：`/chat/sessions`、`/chat/send`（支持**流式**和**同步**模式）
    - 图像识别：`/image/recognize`

- **Controller 层**（如 `controller/user/user.go`、`controller/session/session.go`、`controller/image/image.go`）  
  - 处理 HTTP 请求参数绑定  
  - 调用业务逻辑  
  - 格式化响应  
  - 实现**统一的错误处理**和**标准化响应结构**

---

## 2. **Service 和 DAO**

- **Service 层**（如 `service/user/user.go`、`service/session/session.go`、`service/image/image.go`）  
  封装核心业务逻辑，包括：
  - **用户注册登录**：邮箱验证码验证、JWT 生成  
  - **会话管理**：创建会话、AI 回复生成、多模型支持  
  - **图像识别**：基于 ONNX 模型推理  

- **DAO 层**（如 `dao/user/user.go`、`dao/session/session.go`）  
  - 使用 **GORM** 进行数据库操作  
  - 支持：用户查询、会话创建、消息存储  
  - 实现**数据访问的抽象与复用**

---

## 3. **Middleware**

- **JWT 中间件**（`middleware/jwt/jwt.go`）  
  - 拦截请求并验证 Token  
  - 支持两种方式：
    - `Bearer` 请求头
    - URL 参数  
  - 解析 Token 获取用户名，并注入请求上下文  
  - 确保 API 安全性  
  - 设计简洁，**可快速集成到任意路由组**

---

## 4. **Common 模块**

提供高度复用的基础设施组件：

| 组件 | 功能说明 |
|------|--------|
| **AI 助手模块** (`common/aihelper/`) | 采用**工厂模式**管理不同用户的 AI 实例；支持**同步/流式生成**和**消息历史管理** |
| **图像识别器** (`common/image/image_recognizer.go`) | 基于 **ONNX Runtime** 封装 MobileNetV 2 模型推理 |
| **数据库连接** (`common/mysql/mysql.go`) | 初始化 GORM |
| **Redis 模块** (`common/redis/`) | 处理验证码缓存 |
| **RabbitMQ** (`common/rabbitmq/`) | 实现消息**异步存储** |
| **邮件模块** (`common/email/email.go`) | 支持验证码发送 |

> 所有组件均支持**高并发**与**异步处理**，显著提升系统性能。

---

## 5. **Model**

定义核心数据结构：

- `model/user.go` → `User`：用户账号、邮箱  
- `model/session.go` → `Session`：会话 ID、标题  
- `model/message.go` → `Message`：消息内容、角色标识  

> 模型设计兼容 **GORM 映射**，确保**数据一致性**与**未来扩展性**。

---

## 6. **Frontend（Vue. js）**

- 采用 **Vue. js** 构建 SPA（单页应用）  
- 包含：
  - 路由配置：`vue-frontend/src/router/index.js`  
  - 视图组件：
    - `Login.vue`
    - `Register.vue`
    - `AIChat.vue`
    - `ImageRecognition.vue`
    - `Menu.vue`
- 使用 **Axios** 调用后端 API，实现：
  - 用户认证  
  - AI 对话界面（支持**流式消息显示**）  
  - 图像上传与识别  
- **前后端分离架构**，符合现代 Web 开发规范

---

## 7. **Utils**

提供通用工具函数：

- **JWT 工具**（`utils/myjwt/jwt.go`）  
  - 封装 Token 生成与解析  
  - 基于 `golang-jwt` 库  
  - 支持自定义 Claims  

- **通用工具**（`utils/utils.go`）  
  - MD 5 哈希  
  - 随机数生成  

>  模块化设计，便于在各类组件中**安全复用**


---

# 目录结构：

```go

GopherAI-/
├── vue-frontend/  							# Vue.js前端项目目录
│   ├── vue.config.js  						# Vue项目配置文件，定义代理和构建选项
│   ├── src/  								# 前端源码目录
│   │   ├── views/  							# 页面视图组件
│   │   │   ├── Register.vue  				# 用户注册页面组件
│   │   │   ├── Menu.vue  					# 主菜单页面组件
│   │   │   ├── Login.vue  					# 用户登录页面组件
│   │   │   ├── ImageRecognition.vue  		# 图像识别页面组件
│   │   │   └── AIChat.vue  					# AI聊天对话页面组件
│   │   ├── utils/  							# 前端工具函数
│   │   │   └── api.js  						# API请求封装工具
│   │   ├── router/  						# 前端路由配置
│   │   │   └── index.js  					# Vue Router路由定义
│   │   ├── main.js  						# 前端应用入口文件
│   │   └── App.vue  						# 根组件
│   ├── public/  							# 静态资源目录
│   │   └── index.html  						# HTML模板文件
│   ├── package.json  						# 前端依赖配置
│   └── package-lock.json  					# 前端依赖锁定文件
├── utils/  									# 后端通用工具包
│   ├── utils.go  							# 通用工具函数，如MD5哈希和随机数生成
│   └── myjwt/  								# JWT工具包
│       └── jwt.go  							# JWT令牌生成和解析逻辑
├── service/  								# 业务逻辑服务层
│   ├── user/  								# 用户相关业务服务
│   │   └── user.go  						# 用户登录、注册、验证码发送业务逻辑
│   ├── session/  							# 会话相关业务服务
│   │   └── session.go  						# AI聊天会话管理、消息生成业务逻辑
│   └── image/  								# 图像相关业务服务
│       └── image.go  						# 图像识别业务逻辑
├── router/  								# 路由定义
│   ├── user.go  							# 用户相关API路由
│   ├── router.go  							# 主路由器配置
│   ├── Image.go  							# 图像识别API路由
│   └── AI.go  								# AI聊天API路由
├── model/  									# 数据模型定义
│   ├── user.go  							# 用户数据模型
│   ├── session.go  							# 会话数据模型
│   └── message.go  							# 消息数据模型
├── middleware/  							# 中间件
│   └── jwt/  								# JWT认证中间件
│       └── jwt.go  							# JWT令牌验证中间件
├── main.go  								# Go应用主入口文件
├── go.sum  									# Go模块依赖校验文件
├── go.mod  									# Go模块定义文件
├── dao/  									# 数据访问对象层
│   ├── user/  								# 用户数据访问
│   │   └── user.go  						# 用户数据库操作
│   ├── session/  							# 会话数据访问
│   │   └── session.go  						# 会话数据库操作
│   └── message/  							# 消息数据访问
│       └── message.go  						# 消息数据库操作
├── controller/  							# 控制器层
│   ├── user/  								# 用户控制器
│   │   └── user.go  						# 用户API请求处理
│   ├── session/  							# 会话控制器
│   │   └── session.go  						# AI聊天API请求处理
│   ├── image/  								# 图像控制器
│   │   └── image.go  						# 图像识别API请求处理
│   └── common.go  							# 通用响应结构定义
├── config/  								# 配置管理
│   ├── config.toml  						# 配置文件（TOML格式）
│   └── config.go  							# 配置加载和解析逻辑
└── common/  								# 通用组件库
    ├── redis/  								# Redis缓存组件
    │   ├── redis.go  						# Redis连接和操作
    │   └── key.go  							# Redis键名定义
    ├── rabbitmq/  							# RabbitMQ消息队列组件
    │   ├── rabbitmq.go  					# RabbitMQ连接和发布订阅
    │   ├── meesage.go  						# 消息结构定义
    │   └── init.go  						# RabbitMQ初始化
    ├── mysql/  								# MySQL数据库组件
    │   └── mysql.go  						# MySQL连接和GORM配置
    ├── image/  								# 图像处理组件
    │   └── image_recognizer.go  				# ONNX图像识别器实现
    ├── email/  								# 邮件发送组件
    │   └── email.go  						# 邮件发送逻辑
    ├── code/  								# 错误码定义
    │   └── code.go  							# 统一错误码管理
    └── aihelper/  							# AI助手组件
        ├── model.go  						# AI模型接口定义
        ├── manager.go  						# AI助手管理器
        ├── factory.go  						# AI助手工厂
        └── aihelper.go  					# AI助手核心逻辑


```

---

# 工作原理

> 基于代码分析，将工作原理扩展描述如下：

---

## 1. **请求解析**

- HTTP 请求由 **Gin 框架**接收并解析，提取以下信息：
  - **请求方法**：`GET` / `POST`
  - **路径**：如 `/user/login`
  - **头部信息**：`Authorization` 头（用于 JWT）
  - **请求体**：JSON 格式的参数
- Gin 的路由引擎**快速匹配请求**，支持：
  - 路径参数（如 `/user/:id`）
  - 查询参数（如 `?page=1`）

---

## 2. **路由和中间件**

- **Router 模块**定义的路由规则将请求分发到对应的 **Controller 方法**
- **JWT 中间件**在路由处理前拦截请求：
  - 验证 Token 有效性（解析 Claims 获取用户 ID）
  - 验证失败 → 返回 `401` 错误
- 中间件支持两种 Token 传递方式：
  - `Bearer` 请求头
  - URL 参数  
- 确保所有受保护 API 的**安全性**

---

## 3. **业务处理**

- **Controller 层**：
  - 将请求参数绑定到 Go 结构体
  - 调用 **Service 层**执行核心业务逻辑
- **Service 层**调用 **Common 模块**处理复杂任务：
  - **AI 助手生成回复**：调用外部 API 或本地模型
  - **图像识别**：通过 ONNX 执行模型推理
- **数据操作**通过 **DAO 层**完成：
  - 使用 **GORM** 与 MySQL 交互
  - 遵循 **MVC 模式**，实现关注点分离

---

## 4. **响应生成**

- 业务处理完成后，**Controller 构造统一响应结构**，包含：
  - 状态码
  - 消息（message）
  - 数据（data）
- 响应被序列化为 **JSON** 返回客户端
- 特殊场景支持：
  - **AI 聊天**：使用 **SSE（Server-Sent Events）** 实现**流式响应**，实时推送内容片段
  - **图像识别**：返回分类结果（如 `"cat"`, `"dog"`）

---

## 5. **连接管理**

- **Gin 的事件驱动架构**支持**高并发连接管理**
- 集成 **Redis** 实现：
  - 验证码缓存（TTL = 2 分钟）
  - 会话状态存储
- **RabbitMQ** 用于：
  - 消息的**异步持久化**
  - 避免阻塞主线程  
- 显著**提升系统响应性能**与吞吐量

---

## 6. **数据持久化**

- **MySQL** 存储核心结构化数据：
  - 用户注册信息 → `User` 表
  - 会话与消息 → 通过 RabbitMQ 异步写入 `Message` 表
- **Redis** 缓存热点数据：
  - 如邮箱验证码
- **数据一致性**由 MySQL 保障
- **AI 助手**在内存中维护**消息历史**，支持：
  - 多用户
  - 多会话
  - 完全隔离


> 💡 整个流程形成 **“请求 → 路由 → 鉴权 → 业务 → 响应 → 持久化”** 的闭环，兼顾**安全性、性能与可维护性**。

---

# 项目难点

> 基于代码分析，我将项目难点扩展描述如下：

---

## AI 集成与流式响应

- AI 助手需要处理**外部 API 调用**（如 OpenAI），实现**流式输出**（SSE 协议）以确保实时性，同时避免阻塞主线程。  
- **难点在于**：
  - 管理消息历史上下文
  - 处理 API 错误重试
  - 平衡同步 / 流式模式的性能差异  
-  代码中通过 `AIHelper` 封装模型接口，**支持扩展新 AI 服务**。

---

## **图像识别处理**

- 图像上传需处理**多格式文件**（JPEG / PNG / GIF），并调用 **ONNX Runtime** 进行推理。  
- **难点在于**：
  - 图像预处理（解码、缩放、归一化）
  - 模型加载优化
  - 高并发下的 GPU / CPU 资源管理  
- MobileNetV 2 模型需正确配置输入/输出张量，**避免内存泄漏**。

---

## 并发与异步处理

- **RabbitMQ 异步存储消息**需确保队列可靠性，防止消息丢失。  
- **难点在于**：
  - 消费者设计（确认机制）
  - 多会话并发访问 `AIHelper` 实例的**线程安全性**
  - Redis 缓存的 TTL 管理与过期清理  
- 代码使用 `sync.RWMutex` **保护共享数据**。

---

## 数据库与缓存优化

- **MySQL 与 Redis 集成**需处理数据一致性，例如验证码缓存与 DB 同步。  
- **难点在于**：
  - 连接池配置
  - GORM 事务确保操作**原子性**
  - Redis 键过期策略防止**内存溢出**  

---

## 安全认证与会话管理

- JWT 生成需包含必要 Claims，验证时需解析 Token 以**防止伪造**。  
- **难点在于**：
  - Token 刷新机制
  - 黑名单管理
  - 多设备登录的**会话隔离**  
-  代码中 JWT 中间件支持 `Bearer` 头验证，**密码哈希增强安全性**。

---

##  模块化与扩展性

- `Common` 模块需**解耦设计**，例如：
  - `AIHelper` 工厂支持新模型
  - `ImageRecognizer` 抽象推理接口  
- **难点在于**：
  - 接口定义的通用性
  - 依赖注入避免循环引用
  - 新功能集成时的**向后兼容性**  
- 代码通过**回调函数**和**配置管理**实现灵活扩展。

---

> 这些重难点反映了在构建集成 AI 和多组件的 Web 应用时需要考虑的关键问题。  
> 解决这些问题需要深入理解 **Go 并发模型**、**Web 开发** 和 **系统架构** 等方面的知识。  
> 通过合理的架构设计和代码实现，可以有效地应对这些挑战。

---
# 4：环境准备

# 从零开始配置环境 :

## 0 ：获取邮箱认证

参考如下文章，看到这一点就可以了(我们主要是获取authcode）

[获取授权码](https://blog.csdn.net/weixin_41957626/article/details/131386155)

后续填写如下信息

![[Pasted image 20251208212213.png]]

## 1：获取api-key

进入此链接，获取自己本账号的api-key，后续会用到

[https://bailian.console.aliyun.com/?spm=5176.29619931.J__Z58Z6CX7MY__Ll8p1ZOR.1.1369521crCDcVM&tab=api#/api](https://bailian.console.aliyun.com/?spm=5176.29619931.J__Z58Z6CX7MY__Ll8p1ZOR.1.1369521crCDcVM&tab=api#/api)

需要点击密钥管理这个按钮

## 2：下载Docker

[下载docker](https://blog.csdn.net/jackeydengjun/article/details/147185455)

---

# 梳理

系统功能划分为以下 **3 个核心模块**：

---

##  1. 登录注册（基于 JWT）

- 用户通过**邮箱注册**，系统生成 **11 位随机账号**并发送至邮箱。
- 注册时需验证**邮箱验证码**（Redis 存储，**2 分钟有效**）。
- 登录使用**账号 + 密码**，成功后返回 **JWT Token**：
    - 包含用户 ID 和用户名
    - 有效期可配置
- **中间件验证 Token**，支持：
    - `Authorization` 请求头
    - URL 参数传递 Token
- 验证失败 → 返回错误响应。

---

## 💬 2. AI 聊天对话（包含流式输出与全量输出的选择，多会话管理）

- 支持**多会话**，每个会话**独立管理消息历史**。
- 用户可：
    - 创建新会话
    - 在现有会话中继续对话
- 提供**两种输出模式**：
    - **同步**：返回完整 AI 回复
    - **流式**：通过 **SSE（Server-Sent Events）** 实时推送内容片段（以 `[DONE]` 结束）
- 消息通过 **RabbitMQ 异步存储**到数据库。
- 支持：
    - 获取用户所有会话列表
    - 查询特定会话的历史记录
- **AI 模型类型可配置**（如 GPT 等）
- 使用**上下文管理器**管理不同用户的 AI 助手实例。

---

## 3. 图像识别

- 上传图片文件（支持 **JPEG / PNG / GIF**）
- 使用 **MobileNetV 2 ONNX 模型**进行分类识别
    - 输入尺寸：**224×224**
    - 输出：**ImageNet 1000 类标签**中概率最高的类别
- 识别流程包括：
    - 图片解码
    - 缩放
    - 归一化
    - ONNX 推理
- 最终返回**识别的类别名称**。

---

> 三大模块分别聚焦 **身份认证**、**智能交互** 与 **视觉感知**，构成完整的 AI Web 应用能力闭环。

---

# 用户模块:

## 概述

用户模块是 GopherAI 项目中的核心组件，负责用户的注册、登录、认证和管理。

该模块采用 MVC 架构，使用 Gin 框架处理 HTTP 请求，集成 JWT（JSON Web Token）进行身份验证，并通过 Redis 缓存验证码、MySQL 存储用户数据。模块支持邮箱注册、账号登录，并提供安全的令牌机制确保用户会话的安全性。

在了解用户模块之前，我们先来了解一下什么是JWT？他和token这些有什么区别？
- [什么是 JWT](https://zhuanlan.zhihu.com/p/86937325)
- [token 和 JWT 区别](https://blog.csdn.net/qq_65951398/article/details/130959051)

---

## 工作流程图：

这里拿验证码来举例

![[Pasted image 20251208212625.png]]
- **前端**：Vue组件（Login.vue/Register.vue）发起HTTP请求，发送用户输入的账号密码或邮箱验证码。
- **路由层**：user.go中的Gin路由器匹配路径（如/user/login或/user/register），将请求分发到相应控制器。
- **中间件层**：jwt.go中间件拦截请求，验证JWT Token（注册时可能跳过，登录时必须验证）。
- **控制器层**：user.go控制器（Login/Register/HandleCaptcha）绑定请求参数，调用服务层方法。
- **服务层**：user.go服务（Login/Register/SendCaptcha）执行业务逻辑，如密码验证、验证码检查、生成随机账号。
- **数据访问层**：user.go DAO（IsExistUser/Register）与MySQL交互，查询用户存在性或插入新用户记录。
- **通用组件**：Redis缓存验证码（TTL=2分钟）、Email发送账号邮件、MyJWT生成HS256签名的Token。
- **外部服务**：SMTP邮件服务器接收并发送注册成功邮件。


箭头表示调用方向：前端 → 路由 → 中间件 → 控制器 → 服务 → DAO/通用组件 → 外部服务。颜色区分不同层，便于识别架构层次。

## 用户模型

用户数据模型定义在 `model/user.go` 中，使用 GORM ORM 进行数据库映射。模型结构如下：

```go

type User struct {
    ID        int64          `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"type:varchar(50)" json:"name"`
    Email     string         `gorm:"type:varchar(100);index" json:"email"`
    Username  string         `gorm:"type:varchar(50);uniqueIndex" json:"username"` // 唯一索引
    Password  string         `gorm:"type:varchar(255)" json:"-"`                   // 不返回给前端
    CreatedAt time.Time      `json:"created_at"`                                   // 自动时间戳
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 支持软删除
}

```

- **ID**：主键，自增。
- **Name**：用户昵称。
- **Email**：邮箱地址，支持索引，用于注册和验证码发送。
- **Username**：唯一用户名，用于登录（不同于邮箱）。
- **Password**：MD5 加密后的密码，不返回前端。
- **CreatedAt/UpdatedAt**：自动时间戳。
- **DeletedAt**：软删除支持。

---

## 注册流程

注册流程涉及邮箱验证码验证，确保安全性。流程如下：

1. **发送验证码**：用户提交邮箱，系统生成 6 位随机验证码，存储到 Redis（过期时间 2 分钟），并通过邮件发送。

- API：`POST /captcha`
- 请求体：`{"email": "user@example.com"}`
- 实现：`service/user.SendCaptcha()` 调用 `myredis.SetCaptchaForEmail()` 和 `myemail.SendCaptcha()`。

2. **注册账户**：用户提交邮箱、验证码和密码。

- API：`POST /register`
- 请求体：`{"email": "user@example.com", "captcha": "123456", "password": "password"}`
- 验证步骤：

- 检查用户是否已存在（通过邮箱）。
- 验证 Redis 中的验证码。
- 生成 11 位随机用户名。
- 将用户插入数据库（密码 MD5 加密）。
- 发送用户名到邮箱。
- 生成 JWT Token 并返回。

- 实现：`service/user.Register()` 调用 DAO 层 `user.Register()` 和 JWT 生成。

注册成功后，用户直接进入登录状态，返回 Token。

---

## 登录流程

登录使用用户名（非邮箱）和密码。
- API：`POST /login`
- 请求体：`{"username": "12345678901", "password": "password"}`
- 验证步骤：

	- 检查用户名是否存在。
	- 验证密码（MD5 比较）。
	- 生成 JWT Token 并返回。

- 实现：`service/user.Login()`。

```go

func Login(username, password string) (string, code.Code) {
	var userInformation *model.User
	var ok bool
	//1:判断用户是否存在
	if ok, userInformation = user.IsExistUser(username); !ok {

		return "", code.CodeUserNotExist
	}
	//2:判断用户是否密码账号正确
	if userInformation.Password != utils.MD5(password) {
		return "", code.CodeInvalidPassword
	}
	//3:返回一个Token
	token, err := myjwt.GenerateToken(userInformation.ID, userInformation.Username)

	if err != nil {
		return "", code.CodeServerBusy
	}
	return token, code.CodeSuccess
}

```

---

## JWT Token 机制

JWT 用于无状态认证，避免服务器存储会话。项目使用 `github.com/golang-jwt/jwt/v4` 库。

### Token 结构

Token 包含自定义 Claims 和标准注册声明：

```go

type Claims struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

```

- **ID**：用户 ID。
- **Username**：用户名。
- **RegisteredClaims**：包括过期时间 (exp)、发行者 (iss)、主题 (sub)、发行时间 (iat)。

### 生成 Token

- 函数：`utils/myjwt.GenerateToken(id, username)`
- 配置：从 `config/config.go` 读取过期时长（小时）、密钥、发行者、主题。
- 签名算法：HS256。
- 示例：Token 有效期为配置的 `ExpireDuration` 小时。

```go

func GenerateToken(id int64, username string) (string, error) {
	claims := Claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GetConfig().ExpireDuration) * time.Hour)),
			Issuer:    config.GetConfig().Issuer,
			Subject:   config.GetConfig().Subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().Key))
}

```

### 解析 Token

- 函数：`utils/myjwt.ParseToken(token)`
- 验证签名和过期时间。
- 返回用户名或失败。

```go

// ParseToken 解析Token
func ParseToken(token string) (string, bool) {
	claims := new(Claims)
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Key), nil
	})
	if !t.Valid || err != nil || claims == nil {
		return "", false
	}
	return claims.Username, true
}

```

---

### Token 差别与使用场景

- **Access Token**：项目中生成的即为 Access Token，用于 API 访问。短期有效（配置小时），携带用户身份。
- **Refresh Token**：项目未实现。用于长期会话，Access Token 过期时刷新。区别：Refresh Token 有效期更长，存储更安全。
- **Bearer Token**：项目支持 Bearer 格式（`Authorization: Bearer <token>`），也兼容 URL 参数（`?token=<token>`）。
- **安全性**：Token 不可撤销（除非更改密钥）。项目未实现黑名单机制，若需撤销，可添加 Redis 黑名单。
- **存储**：前端存储在 localStorage 或 sessionStorage，后续请求 Header 中携带。

---

## 中间件认证

认证通过 `middleware/jwt.Auth()` 中间件实现。

- 检查 `Authorization` Header（Bearer 格式）或 URL 参数 `token`。
- 解析 Token，验证有效性。
- 将用户名存储到 Gin 上下文 (`c.Set("userName", userName)`)。
- 失败时返回错误并中止请求。

路由中，未认证接口（如注册、登录）无需中间件；需认证接口（如 AI 聊天）应用中间件。

```go

// 读取jwt
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := new(controller.Response)

		var token string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// 兼容 URL 参数传 token
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		log.Println("token is ", token)
		userName, ok := myjwt.ParseToken(token)
		if !ok {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		c.Set("userName", userName)
		c.Next()
	}
}

```

---
## 验证码机制

- 生成：6 位随机数字。
- 存储：Redis，键为邮箱，值验证码，过期 2 分钟。
- 发送：通过 `common/email` 模块发送邮件。
- 验证：注册时检查 Redis 值。

此机制防止滥用注册，确保邮箱有效性。

---

## 总结

用户模块提供完整的注册登录流程，结合 JWT 实现安全认证。优点：无状态、易扩展；缺点：Token 不可撤销。项目可扩展 Refresh Token 或黑名单以增强安全性。

---

# AI模块 ：

## 概述

AI模块是GopherAI项目的核心组件，专注于集成多种AI模型和服务，提供智能聊天对话功能。该模块采用模块化设计，**支持OpenAI GPT系列、Ollama本地模型等多种AI后端**，实现统一的接口调用和管理。**核心功能包括多会话管理、消息历史维护、同步/流式响应输出，以及异步消息持久化。**

**模块使用工厂模式创建和管理AI助手实例**，每个用户和会话对应独立的AIHelper对象，确保会话隔离和上下文连续性。单例管理器（GlobalManager）维护全局AI助手映射，支持高效的实例获取和生命周期管理。消息历史在内存中维护，通过RabbitMQ异步存储到数据库，避免阻塞主线程。

**支持两种响应模式**：同步模式一次性返回完整AI回复，适用于短对话；流式模式使用SSE（Server-Sent Events）实时推送内容片段，提升用户体验。模块集成JWT认证，确保只有登录用户可访问AI功能。

**架构设计注重高并发和可扩展性**，AIHelper使用读写锁保护消息历史，支持多线程安全访问。配置驱动的模型选择和API密钥管理，便于部署和维护。该模块不仅服务于聊天功能，还为项目提供了AI集成的标准框架，可扩展到其他AI应用场景。


## 工作流程图：

![[Pasted image 20251208213235.png]]

- **前端**：AIChat.vue发送问题，支持同步/流式模式选择。
- **路由层**：AI.go路由匹配/chat/*路径，分发到聊天控制器。
- **中间件层**：jwt.go验证用户Token，确保已登录。
- **控制器层**：session.go控制器（ChatSend/ChatStreamSend等）解析参数，处理会话ID和模型类型。
- **服务层**：session.go服务（CreateSessionAndSendMessage等）管理会话创建、AI回复生成。
- **数据访问层**：session.go DAO创建会话记录，message.go可选消息存储。
- **通用组件**：AIHelper生成回复、AIHelperManager管理实例、AIModelFactory创建模型、RabbitMQ异步存储消息、Redis缓存状态。
- **外部服务**：AI API（如OpenAI/Ollama）提供模型推理。

流程支持新会话创建（生成UUID）和现有会话继续，消息通过RabbitMQ异步持久化，避免阻塞。流式模式用SSE推送实时内容。

---

## AI助手架构

AI助手模块基于AIHelper结构体实现，采用面向对象的封装设计，为每个会话提供独立的AI交互环境。结构体包含消息历史、模型绑定和会话管理，支持动态配置和扩展。

**核心组件：**

● **AIHelper结构体**：核心类封装单个会话的AI交互逻辑，结构体定义如下：
```go

type AIHelper struct {
    model     AIModel                    // AI模型接口，支持不同模型实现
    messages  []*model.Message           // 消息历史列表，存储用户和AI的对话记录
    mu        sync.RWMutex               // 读写锁，保护消息历史并发访问
    SessionID string                     // 会话唯一标识，用于绑定消息和上下文
    saveFunc  func(*model.Message) (*model.Message, error)  // 消息存储回调函数，默认异步发布到RabbitMQ
}

```

- **消息历史**：内存中维护[]*model.Message切片，按时间顺序存储对话消息。每个Message包含SessionID、Content、UserName、IsUser（区分用户/AI消息）和时间戳。历史用于构建AI上下文，确保连续对话的连贯性。
- **模型绑定**：通过AIModel接口动态绑定AI模型，支持运行时切换（如从GPT切换到Ollama）。接口定义GenerateResponse和StreamResponse方法，实现同步和流式生成。
- **会话ID**：字符串类型唯一标识会话，用于消息关联和实例隔离。每个AIHelper实例绑定一个SessionID，确保多会话独立。
- **存储函数**：函数指针类型，默认实现通过RabbitMQ异步发布消息到队列。支持自定义回调，如同步写入数据库，便于测试和扩展。消息发布包含SessionID、Content、UserName、IsUser参数。

### 工作机制：

1. **实例创建**：NewAIHelper构造函数初始化结构体，设置默认saveFunc为RabbitMQ发布。消息列表为空切片，SessionID从参数传入。
2. **消息添加**：AddMessage方法添加新消息到历史，自动调用saveFunc持久化。若Save参数为false，仅内存存储。使用锁保护并发安全。
3. **响应生成**：GenerateResponse/StreamResponse方法构建消息上下文，调用模型接口生成回复。用户消息先添加历史，AI回复后存储。流式模式通过回调实时输出。
4. **历史获取**：GetMessages返回历史副本，避免外部修改。使用读锁确保线程安全。
5. **自定义存储**：SetSaveFunc允许替换存储逻辑，如单元测试中的内存存储。

该架构通过组合模式和策略模式实现灵活性，支持多模型扩展和异步存储，适用于高并发的聊天应用。

---

## 模型接口

模块定义 `AIModel` 接口，支持多种 AI 模型实现。当前支持 OpenAI 和 Ollama。

### 接口定义

```go

type AIModel interface {
    GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
    StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error)
    GetModelType() string
}

```

- **GenerateResponse**：同步生成回复。
- **StreamResponse**：流式生成回复，通过回调函数实时输出。
- **GetModelType**：返回模型类型。

### OpenAI 实现

使用 `github.com/cloudwego/eino-ext/components/model/openai` 库。

- 配置：从环境变量读取 `OPENAI_API_KEY`、`OPENAI_MODEL_NAME`、`OPENAI_BASE_URL`。
- 流式处理：聚合内容并调用回调函数。

### Ollama 实现

使用 `github.com/cloudwego/eino-ext/components/model/ollama` 库。

- 配置：传入 `baseURL` 和 `modelName`。
- 类似 OpenAI，支持本地部署。

## 工厂模式

**AIModelFactory采用工厂模式实现AI模型的创建和管理，支持动态注册和实例化多种AI模型（如OpenAI和Ollama）。工厂使用map存储创建者函数，确保扩展性和解耦。**

**核心实现：**

- **注册创建者**：使用map[string]ModelCreator存储模型创建函数，键为模型类型字符串（如"1"表示OpenAI，"2"表示Ollama）。ModelCreator定义为func(ctx context.Context, config map[string]interface{}) (AIModel, error)。
- **创建模型**：CreateAIModel方法根据modelType从map获取创建者，调用函数实例化模型。传入context和配置参数，返回AIModel接口实例。
- **一键创建助手**：CreateAIHelper方法结合工厂和AIHelper创建，直接返回配置好的助手实例。内部调用CreateAIModel获取模型，然后NewAIHelper创建助手。
- **扩展性**：RegisterModel方法允许运行时注册新模型类型，动态扩展支持的AI服务。全局单例工厂通过sync.Once确保线程安全初始化。

**代码示例：**

```go

// 工厂结构体定义
type AIModelFactory struct {
    creators map[string]ModelCreator
    mu       sync.RWMutex
}

// 注册创建者
func (f *AIModelFactory) RegisterModel(modelType string, creator ModelCreator) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.creators[modelType] = creator
}

// 创建模型
func (f *AIModelFactory) CreateAIModel(ctx context.Context, modelType string, config map[string]interface{}) (AIModel, error) {
    f.mu.RLock()
    creator, exists := f.creators[modelType]
    f.mu.RUnlock()
    if !exists {
        return nil, fmt.Errorf("unknown model type: %s", modelType)
    }
    return creator(ctx, config)
}

// 一键创建助手
func (f *AIModelFactory) CreateAIHelper(ctx context.Context, modelType string, SessionID string, config map[string]interface{}) (*AIHelper, error) {
    model, err := f.CreateAIModel(ctx, modelType, config)
    if err != nil {
        return nil, err
    }
    return NewAIHelper(model, SessionID), nil
}

```

## 管理器

**AIHelperManager采用单例模式管理用户-会话-AIHelper的映射关系，实现实例缓存和生命周期控制。**

**核心实现：**

● **数据结构**：使用map[string]map[string]*AIHelper，外层map键为用户名，内层map键为会话ID，值为AIHelper指针。支持多用户多会话隔离。

● **获取或创建**：GetOrCreateAIHelper方法检查是否存在助手，若无则通过工厂创建并存储。传入用户名、会话ID、模型类型和配置，确保实例唯一性。

● **移除和查询**：RemoveAIHelper删除指定助手，GetAIHelper获取现有实例，GetUserSessions返回用户的所有会话ID列表。

● **全局单例**：GetGlobalManager使用sync.Once返回单例实例，提供统一管理入口。

**代码示例：**

```go

// 管理器结构体定义
type AIHelperManager struct {
    helpers map[string]map[string]*AIHelper  // 用户 -> 会话 -> 助手
    factory *AIModelFactory
    mu      sync.RWMutex
}

// 获取或创建助手
func (m *AIHelperManager) GetOrCreateAIHelper(userName, sessionID, modelType string, config map[string]interface{}) (*AIHelper, error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if m.helpers[userName] == nil {
        m.helpers[userName] = make(map[string]*AIHelper)
    }
    
    if helper, exists := m.helpers[userName][sessionID]; exists {
        return helper, nil
    }
    
    helper, err := m.factory.CreateAIHelper(context.Background(), modelType, sessionID, config)
    if err != nil {
        return nil, err
    }
    
    m.helpers[userName][sessionID] = helper
    return helper, nil
}

// 获取用户会话列表
func (m *AIHelperManager) GetUserSessions(userName string) []string {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    var sessions []string
    if userSessions, exists := m.helpers[userName]; exists {
        for sessionID := range userSessions {
            sessions = append(sessions, sessionID)
        }
    }
    return sessions
}

// 移除助手
func (m *AIHelperManager) RemoveAIHelper(userName, sessionID string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if userSessions, exists := m.helpers[userName]; exists {
        delete(userSessions, sessionID)
        if len(userSessions) == 0 {
            delete(m.helpers, userName)
        }
    }
}

```


## 会话管理

会话基于 `model.Session` 模型，存储在 MySQL 中。

- **Session 模型**：
```go

type Session struct {
    ID        string    `gorm:"primaryKey;type:varchar(36)"`
    UserName  string    `gorm:"index;not null"`
    Title     string    `gorm:"type:varchar(100)"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

```

- **创建会话**：使用 UUID 生成唯一 ID，标题为用户首个问题。
- **会话列表**：`GetUserSessionsByUserName` 返回用户所有会话 ID 和标题。

## 消息处理

消息基于 `model.Message` 模型，支持用户和 AI 消息区分。

- **Message 模型**：
```go

type Message struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    SessionID string    `gorm:"index;not null;type:varchar(36)"`
    UserName  string    `gorm:"type:varchar(20)"`
    Content   string    `gorm:"type:text"`
    IsUser    bool      `gorm:"not null"`
    CreatedAt time.Time
}

```

- **存储策略**：默认异步推送到 RabbitMQ 队列，由消费者处理持久化到 MySQL。
- **历史查询**：`GetMessagesBySessionID` 获取会话消息，按时间排序。

## API接口请求流程举例与MVC架构

**采用经典MVC（Model-View-Controller）架构设计**

实现前后端分离。
- 前端Vue.js作为View层处理用户界面
- 后端Go实现Controller、Service和Model层。

**以下以POST /chat/send-new-session接口为例，展示完整MVC流程：创建新会话并同步发送消息。**

**MVC架构概述：**

- **Model（模型层）**：定义数据结构和业务实体，如model/session.go的Session结构体、model/message.go的Message结构体。负责数据表示和持久化逻辑。
- **View（视图层）**：前端Vue.js组件，如AIChat.vue，负责渲染UI和处理用户输入，通过API调用与Controller交互。
- **Controller（控制器层）**：如controller/session/session.go，接收HTTP请求，调用Service执行业务逻辑，返回响应。

**完整流程举例（POST /chat/send-new-session）：**

1. **前端View发起请求**：用户在AIChat.vue中输入问题，点击发送。前端构造JSON请求体，调用axios.post('/chat/send-new-session', requestData)，requestData包含question、modelType。
2. **Controller接收并处理**：Gin路由分发到session.CreateSessionAndSendMessage控制器方法。方法执行：

	- 绑定请求参数到CreateSessionAndSendMessageRequest结构体。
	- 从JWT中间件获取userName。
	- 调用service.CreateSessionAndSendMessage(userName, req.UserQuestion, req.ModelType)。
```go

func CreateSessionAndSendMessage(c *gin.Context) {
	req := new(CreateSessionAndSendMessageRequest)
	res := new(CreateSessionAndSendMessageResponse)
	userName := c.GetString("userName") // From JWT middleware
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	//内部会创建会话并发送消息，并会将AI回答、当前会话返回
	session_id, aiInformation, code_ := session.CreateSessionAndSendMessage(userName, req.UserQuestion, req.ModelType)

	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.AiInformation = aiInformation
	res.SessionID = session_id
	c.JSON(http.StatusOK, res)
}

```

3. **Service执行业务逻辑**：session.CreateSessionAndSendMessage服务方法执行：

	- 创建新Session实体，调用dao.CreateSession持久化到MySQL。
	- 获取全局AIHelperManager，调用GetOrCreateAIHelper创建AIHelper实例。
	- 调用helper.GenerateResponse生成AI回复（同步模式）。
	- 返回sessionID和aiResponse。

```go

func CreateSessionAndSendMessage(userName string, userQuestion string, modelType string) (string, string, code.Code) {
	//1：创建一个新的会话
	newSession := &model.Session{
		ID:       uuid.New().String(),
		UserName: userName,
		Title:    userQuestion, // 可以根据需求设置标题，这边暂时用用户第一次的问题作为标题
	}
	createdSession, err := session.CreateSession(newSession)
	if err != nil {
		log.Println("CreateSessionAndSendMessage CreateSession error:", err)
		return "", "", code.CodeServerBusy
	}

	//2：获取AIHelper并通过其管理消息
	manager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"apiKey": "your-api-key", // TODO: 从配置中获取,目前这个没什么用
                                  // 跟本地模型对接的，只封装了本地模型对应接口
	}
	helper, err := manager.GetOrCreateAIHelper(userName, createdSession.ID, modelType, config)
	if err != nil {
		log.Println("CreateSessionAndSendMessage GetOrCreateAIHelper error:", err)
		return "", "", code.AIModelFail
	}

	//3：生成AI回复
	aiResponse, err_ := helper.GenerateResponse(userName, ctx, userQuestion)
	if err_ != nil {
		log.Println("CreateSessionAndSendMessage GenerateResponse error:", err_)
		return "", "", code.AIModelFail
	}

	return createdSession.ID, aiResponse.Content, code.CodeSuccess
}

```

4. **DAO数据访问**：dao.CreateSession使用GORM插入Session记录到数据库。AIHelper内部通过RabbitMQ异步存储消息到Message表。
```go

func CreateSession(session *model.Session) (*model.Session, error) {
	err := mysql.DB.Create(session).Error
	return session, err
}

```

5. **Controller返回响应**：构造CreateSessionAndSendMessageResponse，包含aiInformation和sessionID，返回JSON响应给前端。
```go

//内部会创建会话并发送消息，并会将AI回答、当前会话返回
session_id, aiInformation, code_ := session.CreateSessionAndSendMessage(userName, req.UserQuestion, req.ModelType)

if code_ != code.CodeSuccess {
	c.JSON(http.StatusOK, res.CodeOf(code_))
	return
}

res.Success()
res.AiInformation = aiInformation
res.SessionID = session_id
c.JSON(http.StatusOK, res)

```
6. **前端View更新**：接收响应，更新UI显示AI回复，添加新会话到侧边栏。

### 主要接口

AI模块提供RESTful API，支持同步和流式聊天，AI所有接口需JWT认证。

● GET /chat/sessions：获取用户会话列表。Controller调用service.GetUserSessionsByUserName，DAO查询Session表，返回[]SessionInfo。

● POST /chat/send-new-session：创建新会话并发送消息（同步）。如上流程，返回AI回复和sessionID。

● POST /chat/send：向现有会话发送消息（同步）。类似流程，但复用现有sessionID。

● POST /chat/history：获取会话历史。Controller调用service.GetChatHistory，AIHelper返回内存消息历史。

● POST /chat/send-stream-new-session：创建新会话并流式发送消息。Controller设置SSE头，调用service.CreateStreamSessionAndSendMessage，流式推送数据。

● POST /chat/send-stream：向现有会话流式发送消息。复用sessionID，流式输出。

### 请求/响应示例

**发送消息请求**（POST /chat/send-new-session）：

```json

{
  "question": "你好，请介绍一下Go语言",
  "modelType": "1"
}

```

● **响应**：

```json

{
  "code": 200,
  "message": "success",
  "data": {
    "Information": "Go语言是由Google开发的开源编程语言，以简洁、高效著称...",
    "sessionId": "550e8400-e29b-41d4-a716-446655440000"
  }
}

```

## **流式响应**

**AI模块支持流式响应，使用Server-Sent Events (SSE)协议实现实时内容推送，提升用户交互体验。流式模式适用于长文本生成，允许前端逐步显示AI回复。**

### SSE 配置：

- **Content-Type设置**：响应头设置为"text/event-stream"，告知浏览器这是SSE流。
- **缓存控制**：设置"Cache-Control: no-cache"禁用缓存，"Connection: keep-alive"保持连接，"Access-Control-Allow-Origin: *"支持跨域。
- **缓冲禁用**：设置"X-Accel-Buffering: no"防止代理服务器缓冲，确保数据实时到达前端。
```go

// 设置SSE头
c.Header("Content-Type", "text/event-stream")
c.Header("Cache-Control", "no-cache")
c.Header("Connection", "keep-alive")
c.Header("Access-Control-Allow-Origin", "*")
c.Header("X-Accel-Buffering", "no") // 禁止代理缓存

```

### 回调机制：

- **StreamCallback定义**：类型为func(msg string)，在流式生成过程中实时调用。每次接收模型输出片段时，执行回调发送数据。
- **数据格式**：回调内部构造SSE格式字符串"data: " + msg + "\n\n"，通过http.ResponseWriter.Write写入响应流。每个数据块以双换行符结束。
- **刷新机制**：调用http.Flusher.Flush()立即发送数据到客户端，确保低延迟。


### 结束信号：

- **完成标记**：流结束时发送"data: [DONE]\n\n"，前端据此停止监听并关闭连接。
- **错误处理**：若生成失败，通过SSE发送错误事件，如"event: error\ndata: {"message": "Failed to generate"}\n\n"。

### 实现细节：

- **StreamResponse方法**：在AIHelper中实现，接收StreamCallback参数。首先添加用户消息到历史，锁定读取消息构建上下文。调用model.StreamResponse(ctx, messages, cb)，传入回调函数。
- **回调实现**：内部定义cb函数，接收msg字符串，写入响应流并刷新。支持并发安全，通过Flusher确保顺序发送。
- **控制器集成**：ChatStreamSend控制器设置SSE头，调用session.ChatStreamSend传递ResponseWriter。函数内部获取助手，执行StreamMessageToExistingSession。
- **前端处理**：前端EventSource监听"data"事件，累积内容显示。收到[DONE]时处理完成逻辑。

该实现确保流式响应高效且可靠，支持中断恢复和错误处理，适用于实时AI对话场景。

## 总结

AI模块通过工厂模式和单例管理器实现灵活的模型集成，支持多会话并发。异步消息存储确保性能，流式响应提升用户体验。模块易扩展，可添加新模型或存储策略。

---

# 图像识别模块：

## 概述

图像识别模块是GopherAI项目中的核心功能组件，专注于提供高效、准确的图像分类服务。

该模块基于ONNX Runtime集成MobileNetV2预训练模型，实现对上传图像的实时分类识别，支持多种图像格式和输入方式。

模块采用轻量级架构，注重性能优化和资源管理，确保在Go Web应用中无缝集成。

**模块支持单次图像识别模式**，用户可通过API上传图像文件或传递图像缓冲区数据。

内部实现完整的图像预处理流程，包括格式解码、高质量缩放和数据归一化，将图像转换为模型所需的NCHW张量格式。

推理过程使用ONNX Runtime执行MobileNetV2模型，输出1000个ImageNet类别的概率分布，返回最高概率的类别名称。

架构设计强调模块化与扩展性，ImageRecognizer结构体封装所有识别逻辑，支持自定义模型路径、标签文件和输入尺寸。

预分配输入输出张量减少内存分配开销，全局单例管理ONNX环境提升初始化效率。同时提供完善的资源清理机制，确保长时间运行下的内存安全。

**该模块不仅服务于GopherAI的图像识别功能，还可作为独立组件复用于其他Go项目，体现项目在AI集成方面的技术深度和实用价值。**


# 工作流程图：

![[Pasted image 20251208214351.png]]

- **前端**：ImageRecognition.vue上传图片文件，通过FormData发送。
- **路由层**：Image.go路由匹配/image/recognize路径。
- **中间件层**：jwt.go验证用户身份。
- **控制器层**：image.go控制器（RecognizeImage）解析multipart文件。
- **服务层**：image.go服务调用识别器，处理文件读取。
- **通用组件**：ImageRecognizer封装ONNX推理，加载MobileNetV2模型和ImageNet标签。
- **外部服务**：ONNX Runtime执行本地模型推理，返回类别概率。

流程包括图像解码、缩放归一化、推理和结果映射。无DAO层，直接返回识别结果，适合轻量级处理。

# 图像识别架构

图像识别模块基于ImageRecognizer结构体实现，采用ONNX Runtime进行高效模型推理，支持MobileNetV2预训练模型的图像分类。架构设计注重性能优化和资源管理，确保在Go应用中无缝集成。

**核心组件：**

● **ImageRecognizer结构体**：封装完整的图像识别逻辑，包含ONNX会话、输入输出张量、标签列表和尺寸参数。结构体定义如下，支持自定义输入尺寸和模型路径：

```go

type ImageRecognizer struct {
    session      *ort.Session[float32]  // ONNX推理会话
    inputName    string                 // 输入张量名称（默认"data"）
    outputName   string                 // 输出张量名称（默认"mobilenetv20_output_flatten0_reshape0"）
    inputH       int                    // 输入图像高度（默认224）
    inputW       int                    // 输入图像宽度（默认224）
    labels       []string               // ImageNet类别标签列表
    inputTensor  *ort.Tensor[float32]   // 预分配输入张量（1x3xHxW）
    outputTensor *ort.Tensor[float32]   // 预分配输出张量（1x1000）
}

```

● **ONNX会话管理**：使用github.com/yalue/onnxruntime_go库创建和管理推理会话。会话在初始化时加载MobileNetV2 ONNX模型，支持CPU/GPU推理。全局单例模式确保ONNX环境仅初始化一次，避免重复开销。

● **张量预分配**：输入张量形状为[1, 3, H, W]（批次1，通道3，高度H，宽度W），输出张量为[1, 1000]（1000个ImageNet类别概率）。预分配张量减少每次推理的内存分配开销，提高性能。

● **标签加载**：从本地文件（如imagenet_classes.txt）加载1000个ImageNet类别标签，按行读取存储为字符串切片。支持自定义标签文件路径，便于模型替换。

● **输入尺寸配置**：默认224x224像素，支持构造函数中自定义。输入尺寸影响模型精度和性能，需与训练时保持一致。

## 工作流程：

1. **初始化**：调用NewImageRecognizer创建实例，传入模型路径、标签路径和尺寸。函数内部初始化ONNX环境，创建输入输出张量，加载模型到会话，读取标签文件。失败时返回错误，确保资源正确释放。
```go

func NewImageRecognizer(modelPath, labelPath string, inputH, inputW int) (*ImageRecognizer, error) {
	if inputH <= 0 || inputW <= 0 {
		inputH, inputW = 224, 224
	}

	// 初始化 ONNX 环境（全局一次）
	initOnce.Do(func() {
		initErr = ort.InitializeEnvironment()
	})
	if initErr != nil {
		return nil, fmt.Errorf("onnxruntime initialize error: %w", initErr)
	}

	// 预先创建输入输出 Tensor
	inputShape := ort.NewShape(1, 3, int64(inputH), int64(inputW))
	inData := make([]float32, inputShape.FlattenedSize())
	inTensor, err := ort.NewTensor(inputShape, inData)
	if err != nil {
		return nil, fmt.Errorf("create input tensor failed: %w", err)
	}

	outShape := ort.NewShape(1, 1000)
	outTensor, err := ort.NewEmptyTensor[float32](outShape)
	if err != nil {
		inTensor.Destroy()
		return nil, fmt.Errorf("create output tensor failed: %w", err)
	}

	// 创建 Session
	session, err := ort.NewSession[float32](
		modelPath,
		[]string{defaultInputName},
		[]string{defaultOutputName},
		[]*ort.Tensor[float32]{inTensor},
		[]*ort.Tensor[float32]{outTensor},
	)
	if err != nil {
		inTensor.Destroy()
		outTensor.Destroy()
		return nil, fmt.Errorf("create onnx session failed: %w", err)
	}

	// 读取 label 文件
	labels, err := loadLabels(labelPath)
	if err != nil {
		session.Destroy()
		inTensor.Destroy()
		outTensor.Destroy()
		return nil, err
	}

	return &ImageRecognizer{
		session:      session,
		inputName:    defaultInputName,
		outputName:   defaultOutputName,
		inputH:       inputH,
		inputW:       inputW,
		labels:       labels,
		inputTensor:  inTensor,
		outputTensor: outTensor,
	}, nil
}

```


2. **图像预处理**：PredictFromBuffer或PredictFromImage方法接收图像数据。首先解码为image.Image（支持JPEG/PNG/GIF），然后缩放到指定尺寸（使用CatmullRom插值）。转换为NCHW格式float32数组：归一化像素值（0-255到0-1），排列为[R通道, G通道, B通道]。
```go

func (r *ImageRecognizer) PredictFromBuffer(buf []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("failed to decode image from buffer: %w", err)
	}
	return r.PredictFromImage(img)
}
func (r *ImageRecognizer) PredictFromImage(img image.Image) (string, error) {
    	resizedImg := image.NewRGBA(image.Rect(0, 0, r.inputW, r.inputH))


	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	h, w := r.inputH, r.inputW
	ch := 3 // R, G, B
	data := make([]float32, h*w*ch)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := resizedImg.At(x, y)

			r, g, b, _ := c.RGBA()


			rf := float32(r>>8) / 255.0
			gf := float32(g>>8) / 255.0
			bf := float32(b>>8) / 255.0

			// NCHW format
			data[y*w+x] = rf
			data[h*w+y*w+x] = gf
			data[2*h*w+y*w+x] = bf
		}
	}
    //.....
}

```

3. **推理执行**：将预处理数据复制到inputTensor，调用session.Run()执行ONNX推理。输出张量包含1000个类别的概率值。
```go

inData := r.inputTensor.GetData()
copy(inData, data)

```

4. **结果解析**：遍历输出数组找到最大概率索引，映射到labels切片返回类别名称。若索引无效返回"Unknown"。
```go

if err := r.session.Run(); err != nil {
	return "", fmt.Errorf("onnx run error: %w", err)
}

outData := r.outputTensor.GetData()
if len(outData) == 0 {
	return "", errors.New("empty output from model")
}

maxIdx := 0
maxVal := outData[0]
for i := 1; i < len(outData); i++ {
	if outData[i] > maxVal {
		maxVal = outData[i]
		maxIdx = i
	}
}

if maxIdx >= 0 && maxIdx < len(r.labels) {
	return r.labels[maxIdx], nil
}
return "Unknown", nil

```

5. **资源清理**：Close方法销毁会话和张量，释放ONNX资源。defer确保异常情况下资源释放。

## 性能优化：

- **预分配张量**：避免每次推理动态分配内存。
- **单例环境**：ONNX环境全局初始化一次。
- **并发安全**：结构体不含共享状态，支持多实例并发使用。
- **GPU支持**：ONNX Runtime可配置GPU加速，提升推理速度。

该架构模块化设计，便于扩展到其他ONNX模型，支持自定义预处理和后处理逻辑。

# 模型接口

模块定义结构体方法，支持多种输入方式。 接口定义（实际为结构体方法）

- **PredictFromFile**：从文件路径识别图像。
- **PredictFromBuffer**：从字节缓冲区识别图像。
- **PredictFromImage**：核心识别方法，处理 image.Image 对象。
- **Close**：释放 ONNX 资源。

# ONNX 实现

**ONNX 实现**

图像识别模块使用github.com/yalue/onnxruntime_go库实现ONNX模型推理，提供跨平台、高性能的AI模型部署能力。该库封装了ONNX Runtime C API，支持CPU和GPU推理，适合在Go应用中集成预训练模型。

**配置参数：**

● **modelPath**：ONNX模型文件路径，如"/root/models/mobilenetv2/mobilenetv2-7.onnx"。模型文件需预先下载或训练，确保与推理代码兼容。

● **labelPath**：ImageNet类别标签文件路径，如"/root/imagenet_classes.txt"。文件包含1000行，每行一个类别名称，用于将推理结果映射为人可读的标签。

● **inputH 和 inputW**：输入图像高度和宽度，默认224像素。MobileNetV2模型标准输入尺寸，影响识别精度和性能。支持自定义，但需与模型训练时保持一致。

● **默认张量名称**：输入张量名为"data"，输出张量名为"mobilenetv20_output_flatten0_reshape0"。这些名称硬编码在代码中，对应MobileNetV2模型的输入输出节点。

## 预处理流程：

图像预处理是推理前的关键步骤，确保输入数据符合模型要求。PredictFromImage方法执行以下操作：

1. **尺寸缩放**：使用image.NewRGBA创建目标尺寸画布，调用draw.CatmullRom.Scale进行高质量双三次插值缩放。CatmullRom算法提供平滑缩放，保持图像细节。
2. **像素转换**：遍历缩放后的图像每个像素，提取c.RGBA()的R、G、B值（16位）。右移8位转换为8位值，除以255.0归一化到0-1浮点范围。
3. **NCHW格式排列**：数据按通道优先（NCHW）存储：先所有R通道（y_w + x位置），然后G通道，最后B通道。总大小h_w*3，存储为[]float32数组。


## 推理执行：

1. **数据加载**：将预处理float32数组复制到预分配的inputTensor.GetData()切片。
2. **会话运行**：调用session.Run()执行推理，无需额外参数。ONNX Runtime内部处理前向传播，输出概率分布到outputTensor。
3. **结果解析**：获取outputTensor.GetData()，遍历1000个概率值找到最大值索引。索引映射到labels切片返回类别字符串。若索引越界返回"Unknown"。

该实现优化了性能，通过预分配张量和全局环境初始化减少开销，支持并发使用，适用于生产环境下的图像分类任务。

## 接口定义

- **PredictFromFile(imagePath string) (string, error)**：从本地文件路径识别图像。首先使用os.Open打开文件，defer关闭文件句柄。然后调用image.Decode解码为image.Image对象，最后调用PredictFromImage执行识别。该方法适用于从磁盘读取图像文件，内部处理文件I/O和错误检查，确保文件存在且可读。
- **PredictFromBuffer(buf []byte) (string, error)**：从字节缓冲区识别图像。使用bytes.NewReader包装缓冲区，调用image.Decode解码为image.Image。该方法适用于网络上传或内存中的图像数据，无需文件系统操作，提高处理效率。支持JPEG、PNG、GIF等常见格式。
- **PredictFromImage(img image.Image) (string, error)**：核心识别方法，直接处理image.Image对象。首先将图像缩放到指定尺寸（默认224x224），使用draw.CatmullRom.Scale进行高质量插值。然后将像素数据转换为NCHW格式float32数组：遍历每个像素，提取RGBA值，归一化到0-1范围，按通道排列（R、G、B）。数据复制到预分配的inputTensor，调用session.Run()执行ONNX推理。解析outputTensor找到最大概率索引，映射到labels返回类别名称。该方法是其他方法的底层实现，确保一致的预处理和推理逻辑。
- **Close()**：资源释放方法。销毁ONNX会话、输入输出张量，设置指针为nil防止重复释放。应在识别器使用完毕后调用，防止内存泄漏。方法设计为幂等，即使多次调用也不会出错。

这些方法共同提供了完整的图像识别接口，支持文件、缓冲区和内存图像输入，同时确保资源的安全管理和高效推理。接口设计简洁，错误处理完善，便于集成到Web服务中。


# API接口

图像识别模块提供 RESTful API，支持图像上传和分类识别。 主要接口

- POST /recognize：上传图像文件进行识别。 请求/响应示例
- 上传图像请求：multipart/form-data，字段 "image" 为图像文件。
- 响应： { "class_name": "类别名称", "code": 200, "msg": "success" }


# 总结

图像识别模块通过 ONNX Runtime 实现高效图像分类，支持多种输入格式。单次识别模式确保资源隔离，易于集成。模块可扩展，支持更换模型或添加预处理步骤。


#  待扩展的模块：

- RAG 文档问答（Knowledge Base）
- 多模型调用（EINO本身只支持OpenAI格式，阿里百炼符合要求，豆包不符合要求）
- MCP（Model Context Protocol）工具化框架
- TTS ASR服务


---






