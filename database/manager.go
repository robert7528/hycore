package database

import (
	"encoding/json"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// DBManager manages per-tenant database connections with lazy loading and caching.
type DBManager struct {
	adminDB *gorm.DB
	cache   sync.Map // map[tenantCode]*gorm.DB
}

// NewManager creates a DBManager backed by the admin database.
func NewManager(db *gorm.DB) *DBManager {
	return &DBManager{adminDB: db}
}

// GetDB returns the *gorm.DB for the given tenant code (lazy-loaded, cached).
func (m *DBManager) GetDB(tenantCode string) (*gorm.DB, error) {
	if v, ok := m.cache.Load(tenantCode); ok {
		return v.(*gorm.DB), nil
	}

	var cfg TenantDBConfig
	if err := m.adminDB.Where("tenant_code = ?", tenantCode).First(&cfg).Error; err != nil {
		return nil, fmt.Errorf("tenant db config not found for %q: %w", tenantCode, err)
	}

	db, err := m.buildTenantDB(&cfg)
	if err != nil {
		return nil, err
	}

	actual, _ := m.cache.LoadOrStore(tenantCode, db)
	return actual.(*gorm.DB), nil
}

// InvalidateCache removes a tenant's cached DB connection.
// Call this when tenant DB config changes.
func (m *DBManager) InvalidateCache(tenantCode string) {
	m.cache.Delete(tenantCode)
}

func (m *DBManager) buildTenantDB(cfg *TenantDBConfig) (*gorm.DB, error) {
	dsn := cfg.DSN
	if cfg.Mode == "schema" {
		dsn = m.withSearchPath(dsn, cfg.Schema)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect tenant %q: %w", cfg.TenantCode, err)
	}

	replicas := m.parseReplicaDSNs(cfg.ReplicaDSNs, cfg.Mode, cfg.Schema)
	if len(replicas) > 0 {
		dialectors := make([]gorm.Dialector, len(replicas))
		for i, r := range replicas {
			dialectors[i] = postgres.Open(r)
		}
		if err := db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: dialectors,
			Policy:   dbresolver.RandomPolicy{},
		})); err != nil {
			return nil, fmt.Errorf("register replicas for %q: %w", cfg.TenantCode, err)
		}
	}

	return db, nil
}

func (m *DBManager) withSearchPath(dsn, schema string) string {
	if schema == "" {
		return dsn
	}
	return dsn + " search_path=" + schema
}

func (m *DBManager) parseReplicaDSNs(raw, mode, schema string) []string {
	if raw == "" {
		return nil
	}
	var dsns []string
	if err := json.Unmarshal([]byte(raw), &dsns); err != nil {
		return nil
	}
	if mode == "schema" && schema != "" {
		for i, d := range dsns {
			dsns[i] = m.withSearchPath(d, schema)
		}
	}
	return dsns
}
