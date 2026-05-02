package admin

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/pkg/federation/db"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rm = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm"},
	Short:   "Remove a user with administrative access to the ProjDocs cluster",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		cfgDir, cfgDirErr := config.GetConfigDir()
		if cfgDirErr != nil {
			return fmt.Errorf("create config: %w", cfgDirErr)
		}

		db, dbErr := db.Get(db.DefaultConfig(filepath.Join(cfgDir, "federation.db")), zap.NewNop())
		if dbErr != nil {
			return fmt.Errorf("create config: %w", dbErr)
		}

		arg := args[0]

		var err error

		if _, uuidErr := uuid.Parse(arg); uuidErr == nil {
			err = revoke(cmd.Context(), db, "client_id = ?", arg)
		} else if id, intErr := strconv.ParseInt(arg, 10, 64); intErr == nil {
			err = revoke(cmd.Context(), db, "id = ?", id)
		} else {
			return fmt.Errorf("%q is neither a valid integer ID nor a UUID", arg)
		}

		if err != nil {
			return fmt.Errorf("delete: %w", err)
		}

		color.Blue("removed user!")
		return nil
	},
}

func revoke(ctx context.Context, db *sql.DB, where string, arg any) error {
	res, err := db.ExecContext(ctx, `
		UPDATE api_keys
		SET    revoked_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
		WHERE  `+where+`
		AND    revoked_at IS NULL
	`, arg)
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("api key %q not found or already revoked", arg)
	}

	return nil
}
