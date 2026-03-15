# hycore

## Development Environment

- **Windows local**: Source code editing only. No Go runtime available.
- **GitHub**: `robert7528/hycore` (public)
- **CI/CD**: GitHub Actions → `go build` + `go vet` + auto-commit `go.sum`

## What is hycore

Go shared module (`github.com/robert7528/hycore`) for all HySP API modules.
Provides common infrastructure packages so that each module doesn't duplicate them.

## Package Structure

```
hycore/
├── config/       # Viper config (server, database, log, JWT, Tink)
├── logger/       # zap + lumberjack (file rotation + console)
├── database/     # GORM Connect, DBManager (multi-tenant pool), TenantDBConfig
├── middleware/    # tenant, tenant_db, auth (JWT), permission (Casbin), recovery
├── auth/         # Claims, Provider interface, Service (JWT sign/verify), Handler
├── crypto/       # Tink AES-GCM encryptor + NOP fallback
├── casbinx/      # Casbin enforcer with GORM adapter
├── migrator/     # SQL file migration runner
└── auditlog/     # AuditLog model + auto-recording middleware
```

## Usage

```go
// go.mod
require github.com/robert7528/hycore v0.1.0

// Local development (use replace directive, CI blocks this on main)
replace github.com/robert7528/hycore => ../hycore
```

## Key Design Decisions

- All packages are **exported** (not `internal/`), so consuming modules can import them.
- `config.Load()` returns `*Config` (no error). Calls `log.Fatalf` on unmarshal failure.
- `casbinx.NewEnforcer` takes `tableName` as parameter (each module can use its own table).
- `auth.NewClaims` is exported so module-specific auth providers can create claims.
- Table names (`hyadmin_tenant_db_configs`, `hyadmin_audit_logs`) are currently hardcoded.
  These tables are managed centrally by hyadmin.

## Versioning

- Tags follow semver: `v0.1.0`, `v0.2.0`, etc.
- Breaking changes bump minor version (pre-v1).
- Consuming modules pin to specific versions via `go.mod`.
