# kratos_single

一个基于 Go + Kratos 的微服务项目，整合常用功能与真实业务场景，采用清晰ddd分层（service / biz / data），一个服务可同时支持 HTTP + gRPC，兼顾开发效率与后续微服务扩展能力。

---

## 项目定位

本项目适合：

* 想用 Kratos，微服务渐进式升级
* 后台管理系统 / CMS / 电商 / 内容平台
* 需要 API + 管理后台 + 用户系统 的项目
* 未来可平滑升级为微服务架构

---

## 环境要求

* go 1.22
* kratos v2.9.2

---

## 🏗️ 架构设计

```text
API(Request)
   ↓
Service（接口层）
   ↓
Biz（业务逻辑层）
   ↓
Data（数据访问层）
   ↓
MySQL / Redis
```

---

## 📂 目录结构

```text
kratos_single/
├── api/                 # proto 定义
├── cmd/server/          # 启动入口
├── configs/             # 配置文件
├── internal/
│   ├── service/         # 主要业务处理层（核心）
│   ├── biz/             # 业务规则层
│   ├── data/            # 数据访问层
│   ├── server/          # http / grpc 服务注册
│   ├── middleware/      # JWT / Auth / 日志中间件
│   └── pkg/
│       ├── auth/            # JWT / 密码 / context 用户信息
│       ├── utils/           # 常用工具函数
│       └── i18n/            # 多语言
└── README.md
```

---

## 已集成功能（常用后台功能）

### 用户系统

* 注册
* 登录
* JWT Token
* 修改密码
* 修改资料
* 用户状态启用/禁用
* 登录IP记录
* 登录时间记录

### 权限系统

* 用户角色（user/admin）
* JWT 中间件
* 登录接口白名单
* 后台接口鉴权

### 内容管理（Article / Ad / Upload）

* 新增文章
* 修改文章
* 删除文章
* 文章列表
* 前台文章详情（缓存）
* 后台文章详情（直查DB）
* 新增广告
* 修改广告
* 删除广告
* 广告列表
* 前台广告详情（缓存）
* 后台广告详情（直查DB）
* 上传文件

### 通用能力

* MySQL（GORM）
* Redis
* Wire 依赖注入
* 中间件
* 日志
* 多语言 i18n
* 分页查询
* 错误码规范
* 文件上传 OSS
* 多端登录控制
* 定时任务
* 支持 air 热更新

---

## JWT 鉴权

Header：

```text
Authorization: Bearer xxxxxx
```

自动解析：

```go
userID := auth.GetUser(ctx)
role   := auth.GetRole(ctx)
```

---

## 支持未来微服务拆分

当前单体：

```text
user + article + ad + admin
```

未来可拆：

```text
user-service
article-service
ad-service
admin-service
gateway
```

业务代码复用率高。

---

## 技术栈

* Go
* Kratos
* GORM
* MySQL
* Redis
* JWT
* Wire
* Protobuf

---

## 启动项目

```bash
go mod tidy
make config
make api
kratos run
```

---

## 适合业务场景

* 企业后台管理系统
* 内容发布平台
* 电商 API
* 用户中心
* 管理后台 + APP API 共用后端

---

## 后续规划

* RBAC 权限系统
* 操作日志
* 支付系统
* 微服务版升级

---

## License

[MIT License](https://github.com/joanbabyfet/kratos_single/blob/main/LICENSE)