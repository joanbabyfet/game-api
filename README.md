# Slot Provider API

基于 Go 开发的 Slot Provider API。项目负责 Provider 业务、钱包流程与注单管理，并通过自定义 TCP + Protobuf RPC 与 Skynet Slot Game Server 通信。

## 已实现能力

- Go（Gin）Provider API 与 Skynet Game Server 解耦
- 自定义 TCP + Protobuf RPC（Request / Response / Error Packet）
- JWT 玩家认证及 Operator 请求签名验证
- 单一钱包（Single Wallet）与转账钱包（Transfer Wallet）
- 普通 Spin、Free Spin、结算与回滚流程
- 转入、转出及转账状态查询
- request_id、Operator 订单号与转账订单幂等处理
- Game Order 记录钱包模式，确保补偿时沿用原始钱包流程
- Worker 自动补偿未完成注单，包含抢占锁、退避重试与最大重试次数
- Mock Operator，支持 Balance、Bet、Settle、Rollback 及内部幂等订单记录
- Controller / Service / Repository / Adapter 分层
- DTO、领域模型与 Protobuf 解耦
- 统一错误码和响应格式
- 金额以最小货币单位写入核心账务数据，API 使用带小数金额

> RTP、Reel Strip、Payline、Jackpot 与 Risk Control 等游戏逻辑由 Skynet Slot Game Server 负责，不在本仓库内实现。

## 钱包模式

Agent 的 `wallet_mode` 决定余额及 Spin 的处理方式：

| 值 | 模式 | Balance | Spin 账务 |
| --- | --- | --- | --- |
| `1` | 单一钱包 | 请求 Operator `/balance` | 调用 Operator `/bet`、`/settle`、`/rollback` |
| `2` | 转账钱包 | 查询 Provider 本地钱包 | 使用转入后的本地余额完成扣款、派奖与回滚 |

`game_order.wallet_mode` 会保存创建订单时的钱包模式。即使 Agent 后续修改配置，Worker 仍按订单原始模式补偿。

## 系统架构

```text
Operator / Casino
       │
       │ HTTP / JSON
       ▼
Slot Provider API (Go)
       ├──────────────► Operator Wallet API（单一钱包）
       │
       │ TCP + Protobuf RPC
       ▼
Skynet Slot Game Server
       │
       └──────────────► MySQL / Redis

Recovery Worker
       └──────────────► Game Order / Operator / Skynet
```

## 技术栈

- Go 1.23
- Gin
- GORM
- MySQL
- Redis
- Protobuf
- TCP Socket
- robfig/cron

## Provider API

所有 Provider 路由均使用 `/provider` 前缀。

| Method | API | 说明 |
| --- | --- | --- |
| POST | `/provider/player/login` | 玩家登录并签发游戏 Token |
| POST | `/provider/game_url` | 获取游戏启动地址 |
| POST | `/provider/game_list` | 获取游戏列表 |
| POST | `/provider/player/balance` | 按钱包模式查询余额 |
| POST | `/provider/spin` | 按钱包模式执行普通或 Free Spin |
| POST | `/provider/player/kick` | 踢出玩家 |
| POST | `/provider/get_order_log` | 查询注单记录 |
| POST | `/provider/wallet/transfer_in` | 转入 Provider 钱包 |
| POST | `/provider/wallet/transfer_out` | 转出 Provider 钱包 |
| POST | `/provider/wallet/transfer_status` | 查询转账状态 |
| POST | `/provider/ping` | 服务与 Skynet 连通性检查 |
| POST | `/provider/debug/sign` | 生成测试签名，仅供开发测试 |

## Operator Wallet API

单一钱包模式下，Provider 会调用 Operator 提供的四个接口：

| Method | API | 说明 |
| --- | --- | --- |
| POST | `/operator/balance` | 查询玩家余额 |
| POST | `/operator/bet` | 下注扣款 |
| POST | `/operator/settle` | 派奖结算 |
| POST | `/operator/rollback` | 取消下注 |

Operator 必须以订单号实现幂等：

- 相同订单重复 Bet、Settle 或 Rollback 不得重复修改余额。
- 相同订单号但玩家、游戏或金额不一致时应返回冲突错误。
- Rollback 找不到原下注时应成功返回当前余额，不执行加款。

这一约定使 Worker 能对结果不确定的 Pending 单一钱包注单直接执行安全取消。

## 注单补偿

`cmd/worker` 周期扫描可恢复订单，并使用数据库锁避免多个 Worker 重复处理。

| 注单状态 | 补偿动作 |
| --- | --- |
| `Pending` | 单一钱包调用幂等 Rollback，消除下注结果不确定性 |
| `BetSuccess` | 重新向 Skynet 获取游戏结果 |
| `WaitSettle` | 按订单钱包模式重新结算 |
| `WaitRollback` | 按订单钱包模式重新回滚 |

失败任务会记录 `retry_count`、`next_retry_time`、`locked_until` 与 `last_error`，并按配置进行退避重试。

## 项目目录

```text
cmd/
├── admin/          # 管理 API
├── mock/           # Mock Operator
├── provider/       # Provider API
└── worker/         # 注单补偿 Worker

configs/            # 公共及各进程配置
internal/
├── adapter/        # Skynet RPC 适配
├── bootstrap/      # MySQL、Redis、日志、RPC 初始化
├── client/         # Operator 与 Skynet Client
├── config/
├── controller/
├── cron/           # Worker 定时任务
├── dto/
├── middleware/
├── mock/
├── model/
├── repository/
├── router/
└── service/

migrations/         # 数据库增量迁移
pkg/                # 公共错误、签名、金额等工具
proto/              # Protobuf 定义及生成代码
```

## 本地启动

### 1. 准备依赖

需要可连接的 MySQL、Redis 与 Skynet 服务。公共连接配置位于 `configs/app.yaml`。

### 2. 启动服务

```bash
# Provider API，默认监听 :8080
go run ./cmd/provider

# Mock Operator，默认监听 :8082
go run ./cmd/mock

# 注单补偿 Worker
go run ./cmd/worker
```

## License

MIT
