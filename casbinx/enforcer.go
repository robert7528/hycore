// Package casbinx provides a Casbin enforcer backed by GORM (PostgreSQL).
package casbinx

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewEnforcer creates a Casbin enforcer using the GORM adapter.
// tableName is the Casbin policy table (e.g. "hyadmin_casbin_rules").
// modelPath is the Casbin model conf file path.
func NewEnforcer(db *gorm.DB, modelPath, tableName string) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", tableName)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, err
	}
	e.EnableAutoSave(true)
	return e, nil
}
