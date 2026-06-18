# SIWE Authentication — Implementation Tasks

## Flow Overview

```
Frontend                     Gateway                      Redis               Core (gRPC)
   │                            │                          │                     │
   │  GET /api/v1/auth/nonce    │                          │                     │
   │  { wallet_address }        │                          │                     │
   │───────────────────────────►│                          │                     │
   │                            │  1. Generate random nonce │                     │
   │                            │  2. Store nonce (5min TTL)│                     │
   │                            │ ────────────────────────►│                     │
   │  3. Return { nonce,        │                          │                     │
   │       siwe_message }       │                          │                     │
   │◄───────────────────────────│                          │                     │
   │                            │                          │                     │
   │  4. User signs message     │                          │                     │
   │     in MetaMask            │                          │                     │
   │                            │                          │                     │
   │  POST /api/v1/auth/verify  │                          │                     │
   │  { message, signature }    │                          │                     │
   │───────────────────────────►│                          │                     │
   │                            │  5. Parse SIWE message   │                     │
   │                            │  6. Extract nonce        │                     │
   │                            │  7. Consume nonce (DEL)  │                     │
   │                            │ ────────────────────────►│                     │
   │                            │  8. Verify signature     │                     │
   │                            │     (ecrecover)          │                     │
   │                            │  9. GetOrCreateUser      │                     │
   │                            │ ──────────────────────────────────────────────►│
   │                            │ 10. Return user data     │                     │
   │                            │◄──────────────────────────────────────────────│
   │                            │ 11. Generate JWT pair    │                     │
   │                            │ 12. Set HTTP-only cookies│                     │
   │ 13. Return { user,         │                          │                     │
   │       tokens }             │                          │                     │
   │◄───────────────────────────│                          │                     │
```

---

## ✅ Completed Structure

| File | Status |
|---|---|
| `shared/cache/cache.go` | ✅ `Cache` interface + `redisCache` impl |
| `gateway/internal/domain/auth_store.go` | ✅ `AuthStore` interface |
| `gateway/internal/store/redis_auth_store.go` | ✅ `RedisAuthStore` impl |
| `gateway/internal/handler/auth.go` | ✅ Handler stub (3 methods) |
| `gateway/internal/router/router.go` | ✅ Auth routes wired |
| `gateway/internal/app/app.go` | ✅ Dependencies wired |
| `gateway/jwt/jwt.go` | ✅ JWT generate/parse/cookies |

---

## ❌ Task 1 — Add SIWE Dependencies

**Files:** `gateway/go.mod`

```bash
cd gateway
go get github.com/spruceid/siwe-go
go get github.com/ethereum/go-ethereum
go mod tidy
```

**What this brings:**
- `siwe-go` — EIP-4361 SIWE message parsing (`siwe.ParseMessage`) and verification (`message.Verify`)
- `go-ethereum` — `crypto.SigToPub` / `crypto.VerifySignature` for ecrecover (used internally by siwe-go, also needed for address utilities)

---

## ❌ Task 2 — Define DTOs

### `gateway/internal/dtos/request.go`

```go
package dtos

type NonceRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required,hexadecimal,len=42"`
}

type VerifyRequest struct {
	Message   string `json:"message" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}
```

### `gateway/internal/dtos/response.go`

```go
package dtos

type NonceResponse struct {
	Nonce    string `json:"nonce"`
	Message  string `json:"message"`  // Full SIWE message for the user to sign
}

type UserResponse struct {
	ID            string `json:"id"`
	WalletAddress string `json:"wallet_address"`
	EnsName       string `json:"ens_name,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	Role          string `json:"role"`
	IsOnboarded   bool   `json:"is_onboarded"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}
```

---

## ❌ Task 3 — Implement Auth Handler

**File:** `gateway/internal/handler/auth.go`

Replace stubs with real implementations.

### Logic for each handler:

#### `GetNonce(c fiber.Ctx) error`
1. Parse query param / body for `wallet_address` (use `NonceRequest`)
2. Validate it's a valid 42-char hex address (`0x` + 40 hex chars)
3. Generate secure random nonce: `crypto/rand.Read(32)` → hex encode → 64-char string
4. Build SIWE message:

```
${frontendURL} wants you to sign in with your Ethereum account:
${walletAddress}

I accept the LaunchKit Terms of Service

URI: ${frontendURL}
Version: 1
Chain ID: 1
Nonce: ${nonce}
Issued At: ${nowISO}
```

5. Store nonce in Redis via `authStore.SaveNonce(ctx, nonce, 5*time.Minute)`
6. Return `NonceResponse{Nonce: nonce, Message: siweMessage}`

**Key points:**
- Domain must match `cfg.FrontendURL` (validated on verify)
- Nonce is for one-time use
- 5-minute expiry
- The SIWE message is constructed here OR on the frontend — decide which approach. If the frontend builds it, only return the nonce. If the backend builds it, return the full message.

**Recommended approach:** Backend constructs the SIWE message so the domain binding is enforced server-side.

#### `Verify(c fiber.Ctx) error`
1. Parse body into `VerifyRequest{Message, Signature}`
2. Call `siwe.ParseMessage(request.Message)` → get `*siwe.Message`
3. Validate:
   - Domain matches `cfg.JwtConfig.FrontendURL`
   - Nonce exists in Redis (consume it via `authStore.ConsumeNonce`)
   - Message is not expired (check `IssuedAt` + reasonable clock drift)
4. Call `message.Verify(request.Signature, nil, nil, nil)` — this does ecrecover internally
5. Extract `walletAddress := message.GetAddress().Hex()`
6. **Call Core gRPC** `GetOrCreateUser(walletAddress)` to get/create user
7. Build `jwt.TokenPayload` from returned user data
8. Call `jwt.GenerateTokenPair(payload, jwtConfig)` → get `TokenPair`
9. Call `jwt.SetTokenCookies(c, pair, ...)` to set HTTP-only cookies
10. Return `AuthResponse{User: userData, AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}`

**Error cases:**
- Invalid message format → 400
- Nonce not found / already consumed → 401 ("nonce expired or invalid")
- Signature verification fails → 401 ("invalid signature")
- Domain mismatch → 401 ("domain mismatch")
- Message expired → 401 ("message expired")

#### `Logout(c fiber.Ctx) error`
1. Call `jwt.ClearTokenCookies(c)`
2. Return success

---

## ❌ Task 4 — Create Auth Middleware

**New file:** `gateway/internal/middleware/auth.go`

### What it does:
1. Reads `access_token` from cookie (or `Authorization: Bearer <token>` header)
2. Calls `jwt.ParseAccessToken(token, accessKey)`
3. If valid → set `wallet_address`, `user_id`, `role`, `is_onboarded` in `c.Locals()` → `c.Next()`
4. If expired and refresh_token exists → attempt silent refresh:
   - Parse refresh token
   - Generate new token pair
   - Set new cookies
   - Continue with new claims
5. If invalid/expired/no token → return 401

### Helper exports:
```go
func GetWalletAddress(c fiber.Ctx) string
func GetUserID(c fiber.Ctx) string
func GetRole(c fiber.Ctx) string
func GetProjectID(c fiber.Ctx) string
func IsOnboarded(c fiber.Ctx) bool
```

### Routes needing auth:
All routes in `api/v1` except `/auth/*` will require this middleware.

---

## ❌ Task 5 — Implement Core gRPC Service

### 5a — Define Proto

**New file:** `shared/proto/auth/v1/auth.proto`

```protobuf
syntax = "proto3";
package auth.v1;
option go_package = "github.com/Launchkit-org/LaunchKit/shared/proto/auth/v1;authv1";

service AuthService {
  rpc GetOrCreateUser(GetOrCreateUserRequest) returns (GetOrCreateUserResponse);
}

message GetOrCreateUserRequest {
  string wallet_address = 1;
  string ens_name = 2;
}

message GetOrCreateUserResponse {
  string user_id = 1;
  string wallet_address = 2;
  string ens_name = 3;
  string display_name = 4;
  string avatar_url = 5;
  string role = 6;
  bool is_onboarded = 7;
}
```

### 5b — Generate Go code

```bash
protoc --go_out=. --go-grpc_out=. shared/proto/auth/v1/auth.proto
```

### 5c — Core gRPC Server (`core/cmd/main.go`)

Set up:
- gRPC listener on `cfg.Core.HTTPAddr`
- PostgreSQL connection with `db/sqlc`
- Register `AuthServiceServer` implementation
- Handler calls `queries.GetUserByWallet` → if not found, `queries.CreateUser` → return response

### 5d — Add dependencies to `core/go.mod`

```
google.golang.org/grpc
google.golang.org/protobuf
github.com/Launchkit-org/LaunchKit/shared
github.com/Launchkit-org/LaunchKit/db
github.com/jackc/pgx/v5
```

---

## ❌ Task 6 — Gateway gRPC Client

**New file:** `gateway/internal/client/core_client.go`

```go
package client

type CoreClient struct {
	conn *grpc.ClientConn
	client authv1.AuthServiceClient
}

func NewCoreClient(addr string) (*CoreClient, error)
func (c *CoreClient) GetOrCreateUser(ctx context.Context, walletAddress string) (*UserDTO, error)
func (c *CoreClient) Close() error
```

Wire this into `app.go` and pass to `AuthHandler`.

---

## ❌ Task 7 — Wire Middleware in Router

**File:** `gateway/internal/router/router.go`

1. Import `gateway/internal/middleware`
2. Create protected group:
   ```go
   protected := api.Group("", middleware.AuthRequired(jwtCfg))
   ```
3. Add campaign routes etc. to protected group (future)

---

## ❌ Task 8 — Fix Migration (if needed)

**File:** `db/migrations/000001_create_users.sql`

Check line 20: `user_type VARCHAR(20) NOT NULL DEFAULT 'b2c' CHECK (...),` — ensure the comma **is present** after the CHECK constraint (it should be `),` not `);`). Currently it looks correct — verify against what's actually in the file.

---

## ❌ Task 9 — Set Up `.env`

**File:** `.env` (copy from `.env.example`)

```env
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
POSTGRES_DB=launchkit
POSTGRES_USER=launchkit
POSTGRES_PASSWORD=launchkit

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

ACCESS_SECRET=your-256-bit-access-secret-min-32-chars-long!!
REFRESH_SECRET=your-256-bit-refresh-secret-min-32-chars-long!!
```

Also set `CONFIG_PATH` to absolute path of `shared/config/config.yaml`.

---

## Execution Order

| Step | Task | Depends On |
|------|------|------------|
| 1 | Set up `.env` | — |
| 2 | Fix migration if broken | — |
| 3 | Add SIWE deps to gateway | — |
| 4 | Create proto definition | — |
| 5 | Generate proto Go code | 4 |
| 6 | Implement Core gRPC server | 5 |
| 7 | Create gateway gRPC client | 5 |
| 8 | Define DTOs | — |
| 9 | Implement AuthHandler.GetNonce | 3, 8 |
| 10 | Implement AuthHandler.Verify | 3, 7, 8, 9 |
| 11 | Implement AuthHandler.Logout | 8 |
| 12 | Create auth middleware | — |
| 13 | Wire middleware in router | 12 |
| 14 | Test full flow | all |

---

## Key Dependencies

| Library | Purpose | Where |
|---------|---------|-------|
| `github.com/spruceid/siwe-go` | EIP-4361 message parsing + verification | `gateway/handler/auth.go` |
| `github.com/ethereum/go-ethereum` | `common.HexToAddress`, `crypto` utilities | `gateway/handler/auth.go` |
| `google.golang.org/grpc` | Inter-service gRPC | `core/`, `gateway/client/` |
| `google.golang.org/protobuf` | Proto message types | `core/`, `gateway/client/` |
| `github.com/jackc/pgx/v5` | PostgreSQL driver | `core/` |

## Testing the Flow

1. Start infrastructure: `docker compose -f deployments/docker-compose.yml up -d`
2. Run migrations: `goose -dir ./db/migrations postgres "$DATABASE_DSN" up`
3. Start Core: `go run ./core/cmd`
4. Start Gateway: `go run ./gateway/cmd`
5. Test nonce:
   ```
   curl http://localhost:8080/api/v1/auth/nonce?wallet_address=0x...
   ```
6. Test verify (sign message with wallet, then POST):
   ```
   curl -X POST http://localhost:8080/api/v1/auth/verify \
     -H "Content-Type: application/json" \
     -d '{"message":"...","signature":"0x..."}'
   ```
