package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	v1 "github.com/projdocs/cli/pkg/federation/routes/v1"
	"go.uber.org/zap"
)

// RegisterAll mounts all routes.
func RegisterAll(g *gin.RouterGroup, db *sql.DB, log *zap.Logger) {
	v1.RegisterRoutes(g, db, log.Named("routes"))
}
