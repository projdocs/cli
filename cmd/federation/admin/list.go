package admin

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aquasecurity/table"
	"github.com/fatih/color"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/pkg/federation/db"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var list = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List users with administrative access to the ProjDocs cluster",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfgDir, cfgDirErr := config.GetConfigDir()
		if cfgDirErr != nil {
			return fmt.Errorf("create config: %w", cfgDirErr)
		}

		db, dbErr := db.Get(db.DefaultConfig(filepath.Join(cfgDir, "federation.db")), zap.NewNop())
		if dbErr != nil {
			return fmt.Errorf("create config: %w", dbErr)
		}

		if rows, err := db.QueryContext(cmd.Context(), `
		SELECT
			id,
			client_id,
			created_at,
			revoked_at
		FROM api_keys
		ORDER BY id ASC
	`); err != nil {
			return fmt.Errorf("query api keys: %w", err)
		} else {
			defer rows.Close()

			t := table.New(os.Stdout)
			t.SetHeaders("ID", "Client ID", "Created At", "Status")
			t.SetHeaderStyle(table.StyleBold)
			t.SetLineStyle(table.StyleBlue)
			t.SetDividers(table.UnicodeRoundedDividers)

			var count int
			for rows.Next() {
				var (
					id           int64
					clientID     string
					createdAtStr string
					revokedAtStr *string
				)
				if err := rows.Scan(&id, &clientID, &createdAtStr, &revokedAtStr); err != nil {
					return fmt.Errorf("scan api key: %w", err)
				}

				createdAt, err := time.Parse(time.RFC3339, createdAtStr)
				if err != nil {
					return fmt.Errorf("parse created_at %q: %w", createdAtStr, err)
				}

				status := "✓ active"
				if revokedAtStr != nil {
					revokedAt, err := time.Parse(time.RFC3339, *revokedAtStr)
					if err != nil {
						return fmt.Errorf("parse revoked_at %q: %w", *revokedAtStr, err)
					}
					status = "✗ revoked " + revokedAt.Format(time.DateTime)
				}

				t.AddRow(
					fmt.Sprintf("%d", id),
					clientID,
					createdAt.Format(time.DateTime),
					status,
				)
				count++
			}
			if err := rows.Err(); err != nil {
				return fmt.Errorf("rows error: %w", err)
			}

			if count == 0 {
				color.Yellow("no api keys found")
			} else {
				t.Render()
			}
		}

		return nil
	},
}
