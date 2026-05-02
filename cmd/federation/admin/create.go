package admin

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/pkg/federation/db"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

const (
	secretBytes = 32 // 256-bit raw secret
	saltBytes   = 16 // 128-bit salt

	// Argon2id parameters. Time and memory are tuned for
	// interactive authentication (~100ms on modest hardware).
	argonTime    = 2
	argonMemory  = 64 * 1024 // 64 MiB
	argonThreads = 4
	argonKeyLen  = 32
)

// generateSecret returns a cryptographically random URL-safe base64 string.
func generateSecret() (string, error) {
	b := make([]byte, secretBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate secret: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateSalt returns a cryptographically random salt.
func generateSalt() ([]byte, error) {
	b := make([]byte, saltBytes)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}
	return b, nil
}

// hashSecret derives an Argon2id hash of secret using salt.
// The result is base64-encoded for storage.
func hashSecret(secret string, salt []byte) string {
	h := argon2.IDKey(
		[]byte(secret),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)
	return base64.RawStdEncoding.EncodeToString(h)
}

var create = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new user with administrative access to the ProjDocs cluster",
	RunE: func(cmd *cobra.Command, args []string) error {

		secret, err := generateSecret()
		if err != nil {
			return err
		}

		salt, err := generateSalt()
		if err != nil {
			return err
		}

		hash := hashSecret(secret, salt)
		clientID := uuid.New().String()
		saltEncoded := base64.RawStdEncoding.EncodeToString(salt)

		cfgDir, cfgDirErr := config.GetConfigDir()
		if cfgDirErr != nil {
			return fmt.Errorf("create config: %w", cfgDirErr)
		}

		db, dbErr := db.Get(db.DefaultConfig(filepath.Join(cfgDir, "federation.db")), zap.NewNop())
		if dbErr != nil {
			return fmt.Errorf("create config: %w", dbErr)
		}

		_, err = db.ExecContext(cmd.Context(),
			`INSERT INTO api_keys (client_id, secret_hash, salt) VALUES (?, ?, ?)`,
			clientID, hash, saltEncoded,
		)
		if err != nil {
			return fmt.Errorf("insert api key: %w", err)
		}

		color.Blue("    Client ID: %s", clientID)
		color.Blue("Client Secret: %s", secret)

		return nil
	},
}
