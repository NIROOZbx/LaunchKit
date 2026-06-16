# LaunchKit — Agent Context

## What This Project Is

LaunchKit is a **two-sided B2B + B2C multi-chain airdrop campaign platform**.

- **B2B (Projects)**: Crypto protocols, DeFi, NFT collections — create token airdrop campaigns, define eligibility rules, configure tasks (on-chain + social), run sybil detection, and distribute tokens via Merkle-tree proofs. No custom infrastructure needed.
- **B2C (Users)**: Active wallet holders — discover campaigns, complete tasks (on-chain verification + social OAuth), earn points, and claim airdrops with one click.
- **Community Managers**: Project team members — analytics dashboard, participant management, sybil reports, CSV exports.

**Killer feature**: On-chain behavior verification (wallet age, tx count, token/NFT holdings, protocol interactions, staking, DAO voting) + sybil detection engine.

Authentication is wallet-based using SIWE (Sign-In with Ethereum).

---

## Monorepo Structure

```
launchkit/
  gateway/           ← API gateway, routing, middleware, auth validation
  core/              ← Business logic: campaigns, tasks, sybil, eligibility, allocations
  chain/             ← (planned) All blockchain I/O: RPC calls, contract interactions
  verifier/          ← (planned) Task verification: Twitter API, Discord bot, webhooks
  notification/      ← (planned) Notifications for users
  shared/            ← Shared config, JWT, cache, logger, encryptor, serializer, responses
  db/                ← Database schema (goose migrations), sqlc queries + generated Go code
  contracts/         ← Solidity smart contracts (do not modify unless asked)
  deployments/       ← Docker Compose orchestration (Postgres 18, Redis 8, pgAdmin)
  docker-compose.yml
  go.work            ← Go workspace file tying all modules together
```

Each service is an independent Go module with its own `go.mod`.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.26 microservices |
| HTTP Framework | Fiber v3 |
| Auth | SIWE (Sign-In with Ethereum) via `siwe-go` + JWT sessions |
| DB | PostgreSQL 18 + sqlc (no raw SQL strings ever) |
| DB Driver | pgx/v5 |
| Migrations | goose |
| Inter-service | gRPC (`google.golang.org/grpc`) for sync + Kafka (`franz-go`) for async |
| Blockchain | `go-ethereum` for RPC calls, contract interactions |
| Cache | Redis 8 (`go-redis/v9`) |
| Config | Viper (YAML + env overlay) |
| Logging | zerolog + lumberjack |
| Serialization | sonic |
| Encryption | AES-256-GCM |
| Smart Contracts | Solidity 0.8, OpenZeppelin, Hardhat |
| Dev tooling | air (live reload), Docker Compose, Taskfile |
| Chains | Multi-EVM (MVP: Ethereum, Base. Testnets: Sepolia, Base Sepolia via Alchemy) |

---

## Architecture — Five Services

| # | Service | Port | Priority | Role |
|---|---------|------|----------|------|
| 1 | **Gateway** | :8000 | P1 | Public HTTP entry point. Owns SIWE auth, rate limiting (Redis, 100/min per wallet), routing, CORS. Validates JWT, injects `X-Wallet-Address` into gRPC metadata. |
| 2 | **Core** | :8001 | P2 | Business logic: campaign CRUD, task management, task completions, reward computation, eligibility rules, allocation management, sybil detection, campaign status lifecycle. |
| 3 | **Chain** | :8002 | P1 | All blockchain I/O. EVM abstraction via ChainProvider. Merkle tree generation, MerkleDistributor contract deployment, on-chain claim verification, event listening. Single service owns all RPC — no other service touches the chain. |
| 4 | **Verifier** | :8003 | P3 | Async task verification. Consumes Kafka events: Twitter API v2 (follow checks), Discord bot (guild membership), webhook callbacks. Updates task_completions status. |
| 5 | **Notification** | :8004 | P3 | Email/push notifications for claim reminders, campaign updates, etc. |

### Communication

```
Client (Browser/Mobile)
     │ HTTPS (REST JSON)
     ▼
  Gateway (:8000, Fiber)
     │
     ├── gRPC ──► Core (:8001)   — sync: campaigns, tasks, completions, sybil, eligibility
     ├── gRPC ──► Chain (:8002)  — sync: proofs, claim status, tree generation
     │
     └── Kafka ──► Verifier (:8003)     — async: task verification events
                   └── Kafka ◄── Core   — can also produce verification events
     └── Kafka ──► Notification (:8004) — async: notification events
```

- All external traffic enters **only** via Gateway
- Gateway validates JWT locally (shared/jwt) — no gRPC call on every request
- Gateway passes user context (`wallet_address`) via gRPC metadata
- No service directly accesses another service's database
- **Chain Service is the only service that talks to blockchain RPC nodes**
- **Verifier** and **Notification** consume from Kafka topics (async)

---

## Architecture Principles

- **Repository pattern** — all DB access goes through a repository interface; never call the DB directly from a handler or service layer
- **Domain-centric service design** — each service owns its domain completely; no cross-service DB access
- **Interface-driven** — define interfaces in the domain layer; inject implementations
- **Gateway is the only public entry point** — all external HTTP traffic goes through `gateway/`; other services expose gRPC only
- **Auth is stateless at services** — the gateway validates JWT and passes user context (`wallet_address`) downstream via gRPC metadata; individual services trust the gateway
- **Single chain touchpoint** — only the Chain Service communicates with EVM RPC nodes
- **Async verification** — task verification flows through Kafka; gateway/core emit events, Verifier consumes and updates

---

## Go Conventions

### Error handling
```go
// Always wrap errors with context
return fmt.Errorf("campaignService.Create: %w", err)

// Never swallow errors
// WRONG: result, _ := repo.Find(id)
// RIGHT: result, err := repo.Find(id); if err != nil { ... }
```

### Package structure per service
```
<service>/
  cmd/
    main.go              ← entry point, bootstrap, wiring
  internal/
    handler/             ← gRPC handlers (all services except gateway, which uses Fiber handlers)
    middleware/          ← auth, logging, recovery
    service/             ← business logic
    repository/          ← sqlc-generated + repo implementations
    domain/              ← types, interfaces (no external dependencies)
    dto/                 ← request/response structs (gateway only)
    client/              ← downstream service clients (gateway only)
  db/
    migrations/          ← .sql files with goose markers
    queries/             ← .sql files for sqlc
    sqlc.yaml
```

### Database
- **Always use sqlc** — never write raw SQL strings in Go code
- All queries go in `db/queries/*.sql`, generated via `sqlc generate`
- Migrations use **goose** with numeric prefixes: `000001_create_users.sql` — each file contains both `-- +goose Up` and `-- +goose Down` sections
- Every table has `created_at`, `updated_at` timestamps
- Use UUIDs as primary keys (`gen_random_uuid()`)
- Wallet addresses are stored as `VARCHAR(42)` (0x + 40 hex chars)

### HTTP / gRPC
- Gateway exposes REST (Fiber); all other services expose gRPC only
- gRPC framework: `google.golang.org/grpc` (not connect-go)
- Async communication: Kafka via `franz-go`
- REST responses follow this envelope:
```json
{
  "success": true,
  "message": "user found",
  "data": {}
}
```

### Naming
- Files: `snake_case.go`
- Types/interfaces: `PascalCase`
- Functions/methods: `camelCase`
- DB columns: `snake_case`
- Constants: `PascalCase` (not ALL_CAPS)

---

## Auth — SIWE Specifics

- Authentication is **wallet-based only** — no email/password, no Google OAuth
- SIWE flow: frontend requests a nonce → user signs message (MetaMask/Coinbase/WalletConnect) → gateway verifies signature → issues JWT pair
- Wallet address is the primary identity — auto-register on first login
- Gateway owns the entire auth flow (nonce generation, SIWE verification, JWT issue/refresh)
- JWT session: 24h access token + refresh token
- Nonces are single-use, expire in 5 minutes — store in Redis, delete on use
- Never store private keys anywhere in this codebase
- Downstream services receive only `X-Wallet-Address` via gRPC metadata (no raw JWT)
- JWT payload includes: `wallet_address`, `role`

---

## Multi-chain Rules

- Chain identity is a **string name** in the database (e.g. `'ethereum'`, `'base'`, `'arbitrum'`), not an integer chain ID
  - `projects.blockchain` — `VARCHAR(50)`, CHECK constrained to `'ethereum'`, `'base'`, `'arbitrum'`
  - `campaigns.chain` — `VARCHAR(20)`, no CHECK constraint (extensible)
- MVP chains: Ethereum, Base (testnets: Sepolia, Base Sepolia)
- RPC URLs and contract addresses are per-chain — read from config, never hardcoded
- Chain Service should expose a `ChainProvider` interface abstracting EVM RPC interactions

---

## Task Verification System

### On-Chain Tasks (RPC-based, tamper-resistant)

| Task Type | What it checks | Verification |
|-----------|---------------|-------------|
| Hold Token | Min ERC-20 balance | `balanceOf` RPC call |
| Hold NFT | At least 1 NFT | ERC-721 `balanceOf` |
| Wallet Age | Min days since first tx | Binary search for first tx |
| Min ETH Balance | Min native ETH | `eth_getBalance` |
| Min Transaction Count | Min outbound txs | `eth_getTransactionCount` |
| Protocol Interaction | Tx to specific contract | Tx history scan |
| Staked in Contract | Active stake in protocol | Contract-specific method call |
| Voted in DAO | Voted in Snapshot space | Snapshot GraphQL API |

### Social Tasks (OAuth + APIs)

| Task Type | What it checks | Verification |
|-----------|---------------|-------------|
| Connect Twitter | Connected account | OAuth token presence |
| Follow Twitter | Follows specified handle | Twitter API v2 lookup |
| Connect Discord | Connected account | OAuth token presence |
| Join Discord | Member of server | Discord bot API guild member check |
| Refer a Friend | Referred new user | Internal referral tracking |

### Flow
1. User submits completion → Core records `task_completions` with status `pending`
2. Core emits Kafka event → Verifier consumes and processes
3. Verifier calls external API (Twitter, Discord) or marks as `verified`/`rejected`
4. Alternative: Core can verify simple tasks synchronously (on-chain RPC calls go through Chain Service)

---

## Sybil Detection Engine

Assigns 0–100 risk score per wallet. Default threshold (score > 50) excludes from eligibility.

| Signal | Max Penalty | Logic |
|--------|-------------|-------|
| Wallet Age | 40 | 40 if <30d, 20 if <90d |
| Transaction Count | 25 | 25 if <5 txs, 10 if <20 |
| Task Completion Speed | 20 | 20 if <30s, 10 if <120s (bots are instant) |
| Gitcoin Passport Score | 15 | 15 if <10 stamps, 8 if <20 |
| Funding Source Analysis | 30 | 30 if funds shared with 20+ participants |
| IP Address Clustering | 50 | 50 if 5+ wallets share same IP |

- Score is computed post-campaign (batch, goroutines, <60s for 1000 wallets)
- Threshold adjustable per campaign in eligibility rules
- Users can see their score and flags (transparency)

---

## Merkle Distribution

### Flow
1. **Compute eligibility** — combine completed tasks + reward formula + sybil filter
2. **Build Merkle tree** — sorted leaves: `keccak256(abi.encodePacked(walletAddress, amount))`
3. **Deploy MerkleDistributor** — Chain Service deploys contract with root hash
4. **Claim window** — User submits Merkle proof; contract verifies and transfers tokens
5. **Recovery** — Project can recover unclaimed tokens after window expires

### Key rules
- Merkle tree generation happens **only** in Chain Service — no other service touches this
- Proofs are generated per wallet on demand and cached in Redis
- The on-chain contract is the source of truth for claim status — always verify on-chain before marking claimed in DB
- No double-claim, immutable root, time-locked recovery

---

## Smart Contracts

### MerkleDistributor.sol (only contract, per chain)
- Constructor: `token`, `merkleRoot`, `claimWindowDays`
- `claim(campaignId, amount, proof[])` — verifies proof, transfers tokens, emits `Claimed`
- `recoverUnclaimed()` — project owner reclaims unclaimed tokens after window
- LaunchKit never holds custody of tokens — projects approve distributor contract

---

## Database Schema (10 tables)

| Table | Purpose |
|---|---|
| `users` | Wallet-based user profiles with ENS, social identities, `user_type` (b2c/b2b/admin) |
| `projects` | Company/org info, token metadata, treasury wallet, chain |
| `project_api_keys` | Hashed API keys for webhook verification |
| `project_members` | RBAC with invitation workflow (unique per user — single org per user) |
| `campaigns` | Full campaign lifecycle with JSONB configs (reward, eligibility, vesting) |
| `tasks` | Campaign tasks with type, verification type, config, points |
| `task_completions` | User submissions with proof, points, status (pending/verified/rejected) |
| `auth_nonces` | Wallet challenge-response authentication (single-use, 5min expiry) |
| `audit_logs` | Immutable action log per project |
| `campaign_analytics_snapshots` | Periodic campaign metrics |

---

## Current Status (June 2026)

| Component | Status | Notes |
|---|---|---|
| Frontend (UI) | ✅ Complete | Desktop-first glassmorphism design; do not modify |
| DB Schema | ✅ Done | 10 migrations; need to fix migration syntax error in 000001 (line 20: `;` → `,`) |
| Shared config | ✅ Done | Config struct missing `RateLimit` section (exists in YAML) |
| Gateway | 🔄 In progress | Stub with /health only; needs gRPC clients + handlers |
| Core | 🔄 In progress | Stub with /health only; needs business logic |
| Chain | 📋 Planned | |
| Verifier | 📋 Planned | |
| Notification | 📋 Planned | |
| Smart contracts | 📋 Planned | |
| Kafka | 📋 Planned | |
| Kubernetes deployment | 📋 Planned | |

---

## Environment & Local Dev

```bash
# Start all infrastructure
docker compose -f deployments/docker-compose.yml up -d

# Run a specific service (from repo root)
go run ./gateway/cmd
go run ./core/cmd

# Run migrations (from repo root)
goose -dir ./db/migrations postgres "$DATABASE_DSN" up

# Regenerate sqlc (from db/ directory)
sqlc generate

# Tidy all modules
task tidy
```

Environment variables loaded from `.env` (local only, never committed). See `.env.example` for required keys.

---

## What Agents Must Never Do

- **Never modify `frontend/`** unless the task explicitly says so
- **Never modify `contracts/`** unless the task explicitly says so
- **Never write raw SQL** — always use sqlc
- **Never hardcode chain names, contract addresses, or RPC URLs** — use config
- **Never add a Go dependency** without asking first
- **Never access another service's database** — inter-service communication is gRPC or Kafka only
- **Never store secrets in code** — use environment variables loaded via config struct
- **Never generate mock implementations** unless explicitly asked
- **Never call blockchain RPCs from Gateway, Core, Verifier, or Notification** — only Chain Service touches the chain
- **Never modify `shared/responses/response.go` response envelope format** without updating all handlers

---

## When in Doubt

- Follow the pattern already established in `gateway/` or `core/`
- The authoritative source is `docs/PROJECT_CONTEXT.md`
- Ask before making architectural decisions that affect multiple services
- Prefer explicit over clever — this codebase will be worked on by a team of 3
