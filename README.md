# go-link

A production-grade URL shortener built in Go — designed for scale, observability, and extensibility. Built with Echo, PostgreSQL, Redis, Kafka, Elasticsearch, and ClickHouse.

---

## Overview

go-link is a full-stack URL shortening platform that goes beyond the basics. It features structured short code generation, multi-tenant account isolation, event-driven analytics, full-text search, and a billing-aware rate limiting system. The architecture follows a microservices pattern where each service owns its domain and communicates through well-defined interfaces and message queues.

---

## Architecture

```
                        ┌─────────────────────────────────────────────────────┐
                        │              Load Balancer                          │
                        └───────────────────┬─────────────────────────────────┘
                                            │
                        ┌───────────────────▼─────────────────────────────────┐
                        │          API Gateway (GraphQL + REST)               │
                        │     Auth Middleware · Rate Limit Middleware          │
                        └──┬──────────┬──────────┬──────────┬─────────────────┘
                           │          │          │          │
              ┌────────────▼─┐  ┌─────▼──┐  ┌───▼────┐  ┌─▼──────────┐
              │ Auth Service │  │Shortener│  │Redirect│  │  Dashboard │
              │  Users DB    │  │ Service │  │Service │  │  Service   │
              └──────────────┘  └────┬────┘  └───┬────┘  └─────┬──────┘
                                     │            │             │
                              ┌──────▼────────────▼─┐          │
                              │      URL DB          │          │
                              │   (PostgreSQL)       │◄─────────┘
                              └──────────────────────┘
                                     │            │
                              ┌──────▼──┐   ┌────▼──────────────┐
                              │  Redis  │   │  Message Queue     │
                              │ L1 Cache│   │  (Kafka/BullMQ)    │
                              └─────────┘   └──┬──────────────┬──┘
                                               │              │
                                    ┌──────────▼──┐  ┌────────▼───────┐
                                    │  Analytics  │  │  Notification  │
                                    │  Service    │  │  Service       │
                                    └──────┬──────┘  └────────────────┘
                                           │
                                    ┌──────▼──────┐
                                    │  ClickHouse │
                                    │  / BigQuery │
                                    └─────────────┘
```

---

## Services

### Core Services

| Service | Description | Phase |
|---|---|---|
| **Shortener Service** | Creates and manages short URLs. Handles base62 code generation, collision retry, soft delete, and metadata management. | 1 |
| **Redirector Service** | High-throughput redirect resolution. Checks Redis L1 cache first, falls back to URL DB. Fires `url_redirected` event asynchronously. | 1 |
| **Auth Service** | User registration, login, and token lifecycle. Issues Paseto tokens stored in http-only cookies. Talks directly to Users DB. | 2 |
| **API Gateway** | Single entry point. Runs Auth Middleware and Rate Limit Middleware on every request. Exposes GraphQL + REST. | 2 |
| **Analytics Service** | Consumes `url_redirected` events. Extracts IP, geo, user-agent, and timestamp. Stores into ClickHouse for heavy aggregation. | 5 |
| **Dashboard Service** | Aggregates data for the user-facing dashboard. Exposes analytics views and proxies search to Elasticsearch. | 5 |
| **Billing Service** | Manages Stripe subscriptions, plan gating, and webhook handling. Embeds plan into Paseto token payload. | 7 |

### Background Workers

| Worker | Description | Phase |
|---|---|---|
| **Notification Service** | Consumes `url_created` and `url_redirected` events. Fans out to created worker and redirected worker for email/SMS/push delivery. | 4 |
| **Search Worker** | Indexes URL metadata into Elasticsearch on create or update. Batch re-indexes on click count change. | 6 |
| **Expiry Worker** | Cron job that scans for soon-to-expire URLs, sends notifications before expiry, and marks expired URLs inactive. | 5 |

---

## Tech Stack

| Layer | Technology | Reason |
|---|---|---|
| Language | Go (with Echo) | Performance, explicit error handling, excellent concurrency primitives |
| Primary DB | PostgreSQL 16 | ACID transactions, JSONB, partial indexes, UUIDv7 support via `pg_uuidv7` |
| Cache | Redis | L1 redirect cache, token bucket rate limiting per user+route |
| Queue | Kafka / BullMQ | Async event delivery for notifications, analytics, and search indexing |
| Search | Elasticsearch | Full-text search across URL metadata, tags, and UTM parameters |
| Analytics | ClickHouse | Append-only click events at high volume — OLAP workload, not OLTP |
| Auth | Paseto tokens | Stateless, http-only cookies, plan embedded in payload |
| Migrations | Goose | Versioned SQL migration files, rollback support |
| HTTP client | pgx/v5 | Postgres-native driver, pgxpool, precise error codes (23505, etc.) |

---

## Database Schema

### `urls` table

```sql
CREATE TABLE urls (
  id             UUID        PRIMARY KEY DEFAULT uuid_generate_v7(),
  code           TEXT        NOT NULL UNIQUE,
  url            TEXT        NOT NULL,
  account_number VARCHAR(10),
  expires_at     TIMESTAMPTZ,
  click_count    BIGINT      NOT NULL DEFAULT 0,
  state          TEXT        NOT NULL DEFAULT 'active'
                 CHECK (state IN ('active','inactive','expired','banned','pending','archived')),
  metadata       JSONB,
  is_custom      BOOLEAN     NOT NULL DEFAULT false,
  version        BIGINT      NOT NULL DEFAULT 1,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at     TIMESTAMPTZ
);
```

**Indexes**

```sql
CREATE UNIQUE INDEX idx_urls_code           ON urls(code);
CREATE INDEX        idx_urls_account_number ON urls(account_number);
CREATE INDEX        idx_urls_expires_at     ON urls(expires_at)
  WHERE expires_at IS NOT NULL AND state = 'active';
CREATE INDEX        idx_urls_metadata_gin   ON urls USING GIN(metadata);
CREATE INDEX        idx_urls_deleted_at     ON urls(deleted_at)
  WHERE deleted_at IS NULL;
```

---

## Short Code Design

System-generated codes are **8 characters** with a structured anatomy that makes them identifiably ours:

```
1  B7  kR9zQm
│  ││  └──── random payload  (5 chars, crypto/rand base62)
│  └┘─────── account fingerprint  (2 chars, SHA-256 of account_number)
└─────────── version byte  (1 char, currently '1')
```

- **Version byte** at position 0 allows algorithm changes without migrating existing codes. Version `1` codes use today's algorithm; version `2` will use the next.
- **Account fingerprint** at positions 1–2 is a stable 2-char base62 fragment derived from a SHA-256 hash of the account number. Not reversible, but consistent — all codes from the same account share the same fingerprint.
- **Random payload** at positions 3–7 is cryptographically random (`crypto/rand`). 62⁵ ≈ 916 million combinations per fingerprint.

Custom aliases bypass this structure entirely and are stored as-is with `is_custom = true`. Custom codes are validated against a regex to prevent them from accidentally matching the system format.

---

## URL States

| State | HTTP status returned | Triggered by |
|---|---|---|
| `active` | `302 Found` → destination | URL created and scan passed |
| `pending` | `202 Accepted` | New URL awaiting safety scan |
| `inactive` | `404 Not Found` | Owner toggled off |
| `expired` | `410 Gone` | Expiry Worker — `now()` passed `expires_at` |
| `banned` | `451 Unavailable` | Admin or abuse detection pipeline |
| `archived` | `404 Not Found` | Owner soft-deleted from dashboard |

State transitions are enforced in application code via an explicit transition map. `banned` and `archived` are terminal for users — only admins can lift a ban, and deleted URLs cannot be resurrected.

---

## API Endpoints

### Auth Service

| Method | Path | Description |
|---|---|---|
| `POST` | `/sign-up` | Register a new user |
| `POST` | `/sign-in` | Login, returns http-only Paseto cookie |
| `POST` | `/sign-out` | Clear auth cookie |
| `POST` | `/refresh` | Rotate access token |
| `GET` | `/me` | Get current user profile |

### Shortener Service

| Method | Path | Description |
|---|---|---|
| `POST` | `/shorten` | Create a new short URL |
| `GET` | `/urls` | List all URLs for the authenticated user |
| `GET` | `/urls/:id` | Get single URL details |
| `PUT` | `/urls/:id` | Update metadata or expiry (versioned — requires `version` field) |
| `PATCH` | `/urls/:id/toggle` | Activate or deactivate a URL |
| `DELETE` | `/urls/:id` | Soft-delete a URL (sets `archived` state + `deleted_at`) |

### Redirector Service

| Method | Path | Description |
|---|---|---|
| `GET` | `/:code` | Resolve short code and redirect |
| `GET` | `/:code/preview` | Preview destination without redirecting |

### Analytics Service

| Method | Path | Description |
|---|---|---|
| `GET` | `/analytics/:url_id` | Click breakdown for a URL |
| `GET` | `/analytics/:url_id/geo` | Geographic distribution |
| `GET` | `/analytics/:url_id/devices` | Device and browser breakdown |
| `GET` | `/analytics/summary` | Account-wide click summary |

### Billing Service

| Method | Path | Description |
|---|---|---|
| `GET` | `/billing/plans` | List available subscription plans |
| `GET` | `/billing/subscription` | Get current user's subscription |
| `POST` | `/billing/checkout` | Create Stripe Checkout session |
| `POST` | `/billing/portal` | Open Stripe customer portal |
| `POST` | `/billing/webhook` | Stripe webhook receiver |

---

## Subscription Plans

| Feature | Free | Pro |
|---|---|---|
| URLs per day | 10 | Unlimited |
| Redirects per minute | 100 | 1,000 |
| Analytics history | 7 days | Full history |
| Custom aliases | — | ✓ |
| URL expiry control | — | ✓ |
| Safety scan | On every URL | Background only |

Plan is embedded in the Paseto token payload at login time. Rate Limit Middleware reads the plan from the token — no DB call on every request.

---

## Development Setup

### Prerequisites

- Go 1.22+
- Docker and Docker Compose
- `pg_uuidv7` Postgres extension (included in the compose setup)

### Running locally

```bash
git clone https://github.com/your-username/go-link
cd go-link

# Copy environment variables
cp .env.example .env

# Start Postgres, Redis, and Kafka
docker compose up -d

# Run migrations
go run ./cmd/migrate up

# Start the API
go run ./cmd/api
```

### Running tests

```bash
# Unit tests
go test ./...

# Integration tests (requires Docker)
go test ./... -tags integration

# Short code generation stress test
go test ./internal/shortcode/... -run TestGenerateUniqueness -count 100000
```

### Project structure

```
go-link/
├── cmd/
│   ├── api/          # API Gateway entrypoint
│   └── migrate/      # Goose migration runner
├── internal/
│   ├── auth/         # Auth Service
│   ├── shortener/    # Shortener Service
│   ├── redirector/   # Redirector Service
│   ├── analytics/    # Analytics Service
│   ├── dashboard/    # Dashboard Service
│   ├── billing/      # Billing Service
│   ├── notification/ # Notification Service + workers
│   ├── search/       # Search Worker
│   ├── expiry/       # Expiry Worker
│   ├── shortcode/    # Base62 code generation package
│   └── domain/       # Shared types, state machine, errors
├── db/
│   └── migrations/   # Goose SQL migration files
├── docker-compose.yml
└── .env.example
```

---

## Locking Strategy

| Operation | Strategy | Reason |
|---|---|---|
| Insert short URL | Optimistic (INSERT + catch 23505) | Collisions are rare at 62⁷ space — no lock needed |
| Update metadata | Optimistic (`version` column check) | Read-then-write — version mismatch returns 409 |
| Increment click count | None (atomic UPDATE) | Single-row increment is atomic at Postgres row level |
| Soft delete | None (idempotent) | Setting `deleted_at` twice has no negative consequence |
| Billing quota decrement | Pessimistic (`SELECT FOR UPDATE`) | Must atomically check quota and insert in same transaction |

---

## Build Roadmap

| Phase | Scope | Status |
|---|---|---|
| Phase 1 | Shortener Service, URL DB, base62 generation, redirect | 🔧 In progress |
| Phase 2 | Auth Service, Users DB, API Gateway (GraphQL), Paseto tokens | ⏳ Planned |
| Phase 3 | Redis L1 cache, Rate Limit Middleware, cache invalidation | ⏳ Planned |
| Phase 4 | Kafka events, Notification Service, created/redirected workers | ⏳ Planned |
| Phase 5 | Analytics Service, ClickHouse, Dashboard Service, Expiry Worker | ⏳ Planned |
| Phase 6 | Search Worker, Elasticsearch cluster, full-text search | ⏳ Planned |
| Phase 7 | Load Balancer, Billing Service, Stripe integration, hardening | ⏳ Planned |

---

## License

MIT