package database

// TenantDBConfig stores per-tenant database connection settings in the admin DB.
// Mode "database" = separate DSN; mode "schema" = same PostgreSQL, different schema.
type TenantDBConfig struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	TenantCode  string `gorm:"uniqueIndex;not null" json:"tenant_code"`
	Mode        string `gorm:"not null;default:database" json:"mode"` // "database" | "schema"
	DSN         string `gorm:"column:primary_dsn" json:"dsn"`
	Schema      string `json:"schema"`
	ReplicaDSNs string `gorm:"type:text" json:"replica_dsns"` // JSON array of replica DSNs
}

func (TenantDBConfig) TableName() string { return "hyadmin_tenant_db_configs" }
