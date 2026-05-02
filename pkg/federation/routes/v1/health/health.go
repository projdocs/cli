package health

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projdocs/cli/pkg"
	"github.com/projdocs/cli/pkg/federation/response"
	"github.com/projdocs/cli/pkg/types"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	db     *sql.DB
)

type HealthResponse struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

func Get(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(&HealthResponse{
		Version:   pkg.Version,
		Timestamp: time.Now().Format(time.RFC3339),
	}))
}

var RegisterRoutes types.RouteRegistrar = func(g *gin.RouterGroup, d *sql.DB, l *zap.Logger) {
	db = d
	logger = l.Named("health")
	health := g.Group("/health")
	{
		health.GET("", Get)
	}
}
