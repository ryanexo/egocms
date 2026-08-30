package provider

import (
	"cms/internal/httpx"
	"cms/internal/infra/casbin"
	"cms/internal/infra/db"
	"cms/internal/infra/logger"
	"cms/internal/infra/store/gorm/gquery"
	"cms/internal/pkg/hashid"

	"github.com/google/wire"

	"gorm.io/gorm"
)

var InfraProvider = wire.NewSet(
	logger.New,
	db.NewDB,
	casbin.NewRoleCasbin,
	casbin.NewMenuCasbin,
	hashid.New,
	httpx.New,
	NewQuery,
)

func NewQuery(db *gorm.DB) *gquery.Query {
	return gquery.Use(db)
}
