# LaunchKit — Gateway Service

Public HTTP entry point for all client traffic. Owns SIWE authentication, rate limiting, and request routing to downstream services via gRPC.

## Port

`8000` (config: `gateway.http_addr`)

## Dependencies

| Dependency | Purpose |
|------------|---------|
| Redis | SIWE nonces, rate limiting cache, session cache |
| Core (gRPC) | Campaigns, tasks, users, projects, completions |
| Chain (gRPC) | Merkle proofs, on-chain claim status |

The gateway does **not** connect to PostgreSQL — it only uses Redis directly. All business data is accessed through downstream services via gRPC.

## Config

See `shared/config/config.yaml` — `gateway:` section.

| Field | Default | Description |
|-------|---------|-------------|
| `http_addr` | `:8000` | Listen address |
| `read_timeout` | `15s` | Request read timeout |
| `write_timeout` | `15s` | Response write timeout |
| `idle_timeout` | `60s` | Idle connection timeout |

---

## Directory Structure

To maintain clean code and separate concerns without over-engineering, the Gateway is organized as follows:

```
gateway/
  cmd/
    main.go              # Entry point / bootstrapper
    cmd.go               # Server runner and shutdown controls
  internal/
    app/
      app.go             # Dependency injection and wiring
    dtos/
      request.go         # Request validation payloads
      response.go        # Response payload definitions
    handler/
      auth.go            # Wallet authentication handlers (SIWE)
      campaign.go        # Public campaign discovery handlers
      project.go         # Project / Org management handlers (B2B)
    client/
      core_client.go     # Client wrapper connecting to Core service (gRPC)
      chain_client.go    # Client wrapper connecting to Chain service (gRPC)
    router/
      router.go          # Routes HTTP paths to handler methods
```

---

## API Endpoints

### 1. Authentication (Gateway Local)

*Note: Access token refresh is handled transparently on the backend (e.g., via cookies or middleware checks) instead of exposing a dedicated client-facing refresh endpoint.*

| Method | Path | Description | Downstream Target |
|--------|------|-------------|-------------------|
| `GET` | `/api/v1/auth/nonce` | Request SIWE nonce (stores in Redis) | Local (Redis) |
| `POST` | `/api/v1/auth/verify` | Verify SIWE signature, issue JWT pair | Local (Redis) |
| `POST` | `/api/v1/auth/logout` | Invalidate session / delete tokens | Local (Redis) |

### 2. Campaigns & Progress (B2C User JWT Required)

| Method | Path | Description | Downstream Target |
|--------|------|-------------|-------------------|
| `GET` | `/api/v1/campaigns` | List all active/published campaigns (public) | Core (gRPC) |
| `GET` | `/api/v1/campaigns/:id` | Get details and tasks of a campaign (public) | Core (gRPC) |
| `GET` | `/api/v1/campaigns/:id/tasks` | Get tasks annotated with user completion status | Core (gRPC) |
| `POST` | `/api/v1/campaigns/:id/tasks/:taskId/submit` | Submit task completion proof for verification | Core/Verifier |
| `GET` | `/api/v1/campaigns/:id/claims/proof` | Generate or fetch Merkle proof for claims | Chain (gRPC) |
| `GET` | `/api/v1/campaigns/:id/claims/status` | Retrieve claim status from blockchain | Chain (gRPC) |
| `GET` | `/api/v1/users/me/sybil-score` | Get current user's Sybil risk score | Core (gRPC) |

### 3. Campaign & Project Management (B2B Admin JWT Required)

| Method | Path | Description | Downstream Target |
|--------|------|-------------|-------------------|
| `POST` | `/api/v1/projects` | Register a new project / organization | Core (gRPC) |
| `GET` | `/api/v1/projects/:projectId` | Get project organization settings | Core (gRPC) |
| `PUT` | `/api/v1/projects/:projectId` | Update treasury wallet, details, and chain config | Core (gRPC) |
| `POST` | `/api/v1/projects/:projectId/campaigns` | Create a new campaign draft | Core (gRPC) |
| `PUT` | `/api/v1/projects/:projectId/campaigns/:id` | Edit campaign details and vesting schedules | Core (gRPC) |
| `POST` | `/api/v1/projects/:projectId/campaigns/:id/tasks` | Create a new task (e.g. Hold Token) | Core (gRPC) |
| `PUT` | `/api/v1/projects/:projectId/campaigns/:id/tasks/:taskId` | Edit an existing task's parameters | Core (gRPC) |
| `DELETE` | `/api/v1/projects/:projectId/campaigns/:id/tasks/:taskId` | Remove a task from a campaign | Core (gRPC) |
| `POST` | `/api/v1/projects/:projectId/campaigns/:id/publish` | Transition campaign from draft to active | Core (gRPC) |
| `POST` | `/api/v1/projects/:projectId/campaigns/:id/sybil/run` | Trigger a batch Sybil check scoring scan | Core (gRPC) |
| `GET` | `/api/v1/projects/:projectId/campaigns/:id/analytics` | View campaign completion & view metrics | Core (gRPC) |
| `POST` | `/api/v1/projects/:projectId/campaigns/:id/distributor/deploy` | Calculate Merkle root and deploy distributor contract | Chain (gRPC) |

---

## Response Envelope

All responses follow the standard format:

```json
{
  "success": true,
  "message": "campaigns found",
  "data": {}
}
```

Error responses:

```json
{
  "success": false,
  "error": "campaign not found",
  "data": null
}
```

---

## Development

```bash
# From repo root — start infrastructure
task up

# Run migrations
task migrate-up

# Start gateway standalone
go run ./gateway/cmd

# Or use air for live reload
cd gateway && air
```

---

## Architecture

```
Client (Browser/Mobile)
     │ HTTPS
     ▼
  Gateway (:8000)
     │
     ├── gRPC ──► Core (:8001)
     ├── gRPC ──► Chain (:8002)
     │
     ├── Redis (nonces, rate limiting)
     └── Kafka (async) ──► Verifier
                         ──► Notification
```

* Gateway validates JWT locally via `shared/jwt` — no gRPC call per request
* User context (`wallet_address`, `role`) passed to services via gRPC metadata
* No direct database connections in gateway — it authenticates, validates, and routes
