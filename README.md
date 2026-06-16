# LaunchKit

Multi-chain airdrop campaign platform for creating, managing, and distributing token airdrops via gamified task-based campaigns. Two-sided platform serving crypto projects (B2B) and active wallet holders (B2C).

## Services

| # | Service | Port | Protocol | Role |
|---|---------|------|----------|------|
| 1 | **Gateway** | `:8000` | REST (Fiber) | Public entry point, SIWE auth, rate limiting, routing |
| 2 | **Core** | `:8001` | gRPC | Campaigns, tasks, users, projects, sybil detection, eligibility |
| 3 | **Chain** | `:8002` | gRPC | Blockchain I/O, Merkle trees, on-chain claims, event listening |
| 4 | **Verifier** | `:8003` | Kafka | Async task verification (Twitter API, Discord bot, webhooks) |
| 5 | **Notification** | `:8004` | Kafka | Email/push notifications for claims and campaign updates |

## Inter-Service Communication

```
Client (Browser/Mobile)
     │ HTTPS (REST JSON)
     ▼
  Gateway (:8000)
     │
     ├── gRPC (sync) ────► Core (:8001)
     ├── gRPC (sync) ────► Chain (:8002)
     │
     └── Kafka (async) ──► Verifier (:8003)
                          Notification (:8004)
```

- All external traffic enters only via **Gateway**
- Gateway validates JWT locally — no gRPC call per request
- User context (`wallet_address`) passed via gRPC metadata
- **Chain** is the only service that talks to blockchain RPC nodes
- Async flows (task verification, notifications) use **Kafka**

## Architecture

```
launchkit/
├── gateway/           ← Public HTTP API (SIWE auth, routing)
├── core/              ← Business logic (campaigns, tasks, sybil)
├── chain/             ← Blockchain I/O (Merkle trees, RPC)
├── verifier/          ← Async task verification
├── notification/      ← Notifications
├── shared/            ← Config, JWT, cache, logger, encryptor, serializer
├── db/                ← Migrations (goose), sqlc queries + generated code
├── frontend/          ← React + TypeScript (do not modify)
├── contracts/         ← Solidity smart contracts
└── deployments/       ← Docker Compose (Postgres 18, Redis 8, pgAdmin)
```

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP Framework | Fiber v3 |
| Inter-service | gRPC (`google.golang.org/grpc`) |
| Async | Kafka (`franz-go`) |
| Database | PostgreSQL 18 (pgx/v5) |
| DB Codegen | sqlc |
| Migrations | goose |
| Cache | Redis 8 (go-redis/v9) |
| Config | Viper (YAML + env overlay) |
| Auth | SIWE (`siwe-go`) + JWT HS256 |
| Logging | zerolog + lumberjack |
| Serialization | sonic |
| Encryption | AES-256-GCM |
| Blockchain | go-ethereum |
| Smart Contracts | Solidity 0.8, OpenZeppelin, Hardhat |
| Frontend | React + TypeScript, Vite, wagmi/viem |
| Dev | air (live reload), Docker Compose, Taskfile |

## Database Schema

10 tables covering the full domain:

| Table | Purpose |
|---|---|
| `users` | Wallet-based user profiles with ENS and social identities |
| `projects` | Company/org info, token metadata, treasury wallet |
| `project_api_keys` | Hashed API keys for webhook verification |
| `project_members` | RBAC with invitation workflow |
| `campaigns` | Full campaign lifecycle with JSONB configs |
| `tasks` | Campaign tasks with type, verification type, config, points |
| `task_completions` | User submissions with proof, points, leaderboard |
| `auth_nonces` | Wallet challenge-response authentication |
| `audit_logs` | Immutable action log per project |
| `campaign_analytics_snapshots` | Periodic campaign metrics |

## Development

### Prerequisites

- Go 1.26+
- Docker & Docker Compose
- Taskfile ([task](https://taskfile.dev))

### Quick Start

```bash
# Start infrastructure and boot services
task start

# Or step by step:
task up         # docker compose up -d
task migrate-up  # apply database migrations
task logs       # tail all service logs
```

### Run a Service Standalone

```bash
# Gateway
go run ./gateway/cmd

# Core
go run ./core/cmd
```

### Available Tasks

| Task | Description |
|---|---|
| `up` | Start Docker services |
| `down` | Stop and remove all containers/volumes |
| `build` | Rebuild Docker images |
| `logs` | Tail container logs |
| `migrate-status` | Show migration status |
| `migrate-up` | Apply pending migrations |
| `migrate-down` | Roll back last migration |
| `migrate-redo` | Re-run latest migration |
| `migrate-create -- <name>` | Create a new migration |
| `gen-sqlc` | Regenerate sqlc Go code |
| `tidy [module]` | Run `go mod tidy` on a module or all |

### Environment

Copy `.env.example` to `.env` and configure:

| Variable | Purpose |
|---|---|
| `POSTGRES_*` | Database connection |
| `REDIS_ADDR`, `REDIS_PASSWORD` | Redis connection |
| `ACCESS_SECRET`, `REFRESH_SECRET` | JWT signing keys |
| `CONFIG_PATH` | Path to YAML config |

## Service READMEs

Each service has its own README with service-specific details:

- [gateway/README.md](./gateway/README.md)
- core/README.md
- chain/README.md
- verifier/README.md
- notification/README.md

## Status

| Component | Status | Notes |
|---|---|---|
| Frontend (UI) | ✅ Complete | Desktop-first glassmorphism design |
| DB Schema | ✅ Done | 10 migrations (goose) |
| Shared config | ✅ Done | Config, JWT, cache, logger, serializer |
| Gateway | 🔄 In progress | Auth, routing, rate limiting |
| Core | 🔄 In progress | Campaigns, tasks, sybil, eligibility |
| Chain | 📋 Planned | Blockchain I/O, Merkle trees |
| Verifier | 📋 Planned | Async task verification |
| Notification | 📋 Planned | Email/push notifications |
| Smart contracts | 📋 Planned | MerkleDistributor.sol |
| Kafka | 📋 Planned | Async event pipeline |
| Kubernetes | 📋 Planned | Production deployment |
