package federation

import (
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/internal/logging"
	"github.com/projdocs/cli/pkg/federation"
	"github.com/projdocs/cli/pkg/federation/db"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	serveJsonOutput bool
	serveVerbose    bool
	serveMode       = logging.ModeConsole
	serveLevel      = zapcore.InfoLevel
)

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Serve the federation cluster host",
	PreRun: func(cmd *cobra.Command, args []string) {
		if serveJsonOutput {
			serveMode = logging.ModeJSON
		}
		if serveVerbose {
			serveLevel = zapcore.DebugLevel
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		logger := logging.NewZapLogger(serveMode, serveLevel)

		log := logger.Named("cmd")
		log.Info("starting federation cluster")

		cfgDir, cfgDirErr := config.GetConfigDir()
		if cfgDirErr != nil {
			log.Fatal("could not get config dir", zap.Error(cfgDirErr))
		}

		db, dbErr := db.Get(db.DefaultConfig(filepath.Join(cfgDir, "federation.db")), logger.Named("db"))
		if dbErr != nil {
			log.Fatal("could not get federation database", zap.Error(dbErr))
		}

		srv := federation.NewServer(federation.DefaultServerConfig(), db, logger.Named("server"))
		if err := srv.Serve(ctx); err != nil {
			log.Fatal("could not start server", zap.Error(err))
		}

		log.Info("federation cluster stopped")
		return nil
	},
}

func init() {
	serve.Flags().BoolVarP(&serveJsonOutput, "json", "j", false, "Output as JSON")
	serve.Flags().BoolVarP(&serveVerbose, "verbose", "v", false, "Verbose output")
}
