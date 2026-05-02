package v1

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/projdocs/cli/pkg/federation/routes/v1/health"
	"github.com/projdocs/cli/pkg/types"
	"go.uber.org/zap"
)

var RegisterRoutes types.RouteRegistrar = func(g *gin.RouterGroup, db *sql.DB, log *zap.Logger) {

	v1 := g.Group("/v1")

	// register health routes
	health.RegisterRoutes(v1, db, log.Named("v1"))
}
