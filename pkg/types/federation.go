package types

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteRegistrar = func(g *gin.RouterGroup, db *sql.DB, log *zap.Logger)
