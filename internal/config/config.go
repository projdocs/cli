package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/projdocs/cli/internal/config/keys"
)

type Config struct {
	File     File
	Supabase Supabase
}

type SupabaseKeysPair struct {
	Symmetric  string
	Asymmetric string
}

type SupabaseKeys struct {
	Anon        SupabaseKeysPair
	Service     SupabaseKeysPair
	Publishable string
	Secret      string
	JWTKeys     string
	JWTJWKS     string
	JWTSecret   string
}

type Postgres struct {
	Password string
}

type Studio struct {
	Password string
}

type Storage struct {
	MinioPassword string
	S3AccessKeyID string
	S3SecretKey   string
}

type Supabase struct {
	Keys     SupabaseKeys
	Postgres Postgres
	Studio   Studio
	Storage  Storage
}

func FromFile(file *File) (*Config, error) {

	pgPass, pgPassErr := generateHex(16)
	if pgPassErr != nil {
		return nil, fmt.Errorf("failed to generate Postgres password: %w", pgPassErr)
	}

	dashboardPass, dashboardPassErr := generateHex(16)
	if dashboardPassErr != nil {
		return nil, fmt.Errorf("failed to generate dashboard password: %w", dashboardPassErr)
	}

	minioPass, minioPassErr := generateHex(16)
	if minioPassErr != nil {
		return nil, fmt.Errorf("failed to generate dashboard password: %w", minioPassErr)
	}

	s3AccessKey, s3AccessKeyErr := generateHex(16)
	if s3AccessKeyErr != nil {
		return nil, fmt.Errorf("failed to generate S3 access key: %w", s3AccessKeyErr)
	}

	s3Secret, s3SecretErr := generateHex(32)
	if s3SecretErr != nil {
		return nil, fmt.Errorf("failed to generate S3 secret key: %w", s3SecretErr)
	}

	jwtSecret, err := keys.GenerateJWTSecret()
	if err != nil {
		return nil, fmt.Errorf("unable to generate JWT secret: %w", err)
	}

	pemBytes, err := keys.GenerateECKey()
	if err != nil {
		return nil, fmt.Errorf("unable to generate EC key: %w", err)
	}

	if derivedKeys, err := keys.GenerateDerivedKeys(pemBytes, jwtSecret); err != nil {
		return nil, fmt.Errorf("failed to generate cryptographic keys: %w", err)
	} else {
		return &Config{
			File: *file,
			Supabase: Supabase{
				Postgres: Postgres{
					Password: pgPass,
				},
				Studio: Studio{
					Password: dashboardPass,
				},
				Storage: Storage{
					MinioPassword: minioPass,
					S3AccessKeyID: s3AccessKey,
					S3SecretKey:   s3Secret,
				},
				Keys: SupabaseKeys{
					JWTJWKS:     derivedKeys.JWTJWKS,
					JWTKeys:     derivedKeys.JWTKeys,
					JWTSecret:   jwtSecret,
					Publishable: derivedKeys.PublishableKey,
					Secret:      derivedKeys.SecretKey,
					Anon: SupabaseKeysPair{
						Symmetric:  derivedKeys.AnonKey,
						Asymmetric: derivedKeys.AnonKeyAsymmetric,
					},
					Service: SupabaseKeysPair{
						Symmetric:  derivedKeys.ServiceRoleKey,
						Asymmetric: derivedKeys.ServiceRoleKeyAsymmetric,
					},
				},
			},
		}, nil
	}
}

func generateHex(n int) (string, error) {
	raw := make([]byte, n)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate hex: %w", err)
	}
	return hex.EncodeToString(raw), nil
}
