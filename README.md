# Slot Provider API

基于 Go 开发的 Slot Provider API，负责对接第三方 Operator，并通过自定义 TCP + Protobuf RPC 与 Skynet Slot Game Server 通信。

## 项目特点

* Go（Gin）+ Skynet 双服务架构
* 自定义 TCP + Protobuf RPC（Request / Response / Error Packet）
* Provider API 与 Game Server 完全解耦
* JWT 玩家认证
* 单一钱包（Single Wallet）及 Mock Operator 测试框架
* Slot Game 完整游戏流程（Login / Balance / Spin / Free Spin / Rollback）
* RTP 统计、Jackpot Pool、Risk Control 风控机制
* Controller / Service / Repository / Adapter 分层架构
* DTO 与 Protobuf 解耦
* 统一错误码、统一响应格式、幂等处理
* 金额统一采用最小货币单位存储，支持多币种（Currency）
* 支持多 Provider、多 Agent、多游戏扩展（Slot / Fish / Casino）

---

## 技术栈

* Go
* Gin
* GORM
* MySQL
* Redis
* Protobuf
* TCP Socket

---

## 项目架构

```text
                +----------------------+
                |   Operator / Casino  |
                +----------+-----------+
                           |
                       HTTP / JSON
                           |
                           v
                +----------------------+
                |   Slot Provider API  |
                |        (Go)          |
                +----------+-----------+
                           |
                TCP + Protobuf RPC
                           |
                           v
                +----------------------+
                | Slot Game Server     |
                |   Skynet + Lua       |
                +----------+-----------+
                           |
                  Redis + MySQL
```

---

## Provider API

目前已完成：

| Method | API | Status |
|---------|-----|--------|
| POST | `/authenticate` | ✅ |
| POST | `/balance` | ✅ |
| POST | `/bet` | ✅ |
| POST | `/rollback` | ✅ |
| GET | `/history` | ✅ |
| GET | `/game/list` | ✅ |
| POST | `/player/kick` | ✅ |

---

## 项目目录

```text
cmd/

configs/

internal/
├── adapter/
├── bootstrap/
├── client/
│   └── skynet/
├── config/
├── controller/
├── dto/
├── middleware/
├── model/
├── repository/
├── router/
└── service/

pkg/

proto/
```

---

## RPC 架构

```text
Controller
        │
        ▼
Service
        │
        ▼
Adapter
        │
        ▼
TCP + Protobuf
        │
        ▼
Skynet
```

---

## 快速启动

```bash
go run ./cmd/provider/main.go
```

---

## 对应游戏服务器

本项目负责 Provider API，不包含 Slot 游戏逻辑。

Slot Game Server：

* Skynet
* Lua
* Reel Strip
* Payline
* Jackpot
* RTP
* Risk Control
* Wallet
* Rollback

---

## License

MIT