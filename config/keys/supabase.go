package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

const projectRef = "supabase-self-hosted"

// ── Types ─────────────────────────────────────────────────────────────────────

type JWK struct {
	Kty    string   `json:"kty"`
	Kid    string   `json:"kid,omitempty"`
	Use    string   `json:"use,omitempty"`
	KeyOps []string `json:"key_ops,omitempty"`
	Alg    string   `json:"alg"`
	Ext    bool     `json:"ext,omitempty"`
	Crv    string   `json:"crv,omitempty"`
	X      string   `json:"x,omitempty"`
	Y      string   `json:"y,omitempty"`
	D      string   `json:"d,omitempty"`
	K      string   `json:"k,omitempty"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type GeneratedKeys struct {
	PublishableKey           string
	SecretKey                string
	AnonKeyAsymmetric        string
	AnonKey                  string
	ServiceRoleKeyAsymmetric string
	ServiceRoleKey           string
	JWTKeys                  string // JSON array — for Auth to sign
	JWTJWKS                  string // JSON object — for PostgREST, Realtime, Storage to verify
}

// ── Public entry point ────────────────────────────────────────────────────────
func GenerateDerivedKeys(pemBytes []byte, jwtSecret string) (*GeneratedKeys, error) {
	// Parse EC private key
	privateKey, err := parseECPrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("parse EC private key: %w", err)
	}

	kid, err := randomUUID()
	if err != nil {
		return nil, fmt.Errorf("generate kid: %w", err)
	}

	// EC JWK components
	ecPrivateJWK := ecPrivateKeyToJWK(privateKey, kid)
	ecPublicJWK := ecPublicKeyToJWK(privateKey, kid)

	// Symmetric key
	octKey := JWK{
		Kty: "oct",
		K:   base64.RawURLEncoding.EncodeToString([]byte(jwtSecret)),
		Alg: "HS256",
	}

	// JWKS for Auth (signs — includes private key)
	jwksKeypair := JWKS{Keys: []JWK{ecPrivateJWK, octKey}}

	// JWKS for PostgREST/Realtime/Storage (verifies — public key only)
	jwksPublic := JWKS{Keys: []JWK{ecPublicJWK, octKey}}

	// JWTs
	iat := time.Now().Unix()
	exp := iat + 5*365*24*3600

	anonJWT, err := signES256(privateKey, kid, map[string]any{
		"role": "anon",
		"iss":  "supabase",
		"iat":  iat,
		"exp":  exp,
	})
	if err != nil {
		return nil, fmt.Errorf("sign anon JWT: %w", err)
	}

	serviceJWT, err := signES256(privateKey, kid, map[string]any{
		"role": "service_role",
		"iss":  "supabase",
		"iat":  iat,
		"exp":  exp,
	})
	if err != nil {
		return nil, fmt.Errorf("sign service JWT: %w", err)
	}

	// Opaque keys
	publishableKey, err := generateOpaqueKey("sb_publishable_")
	if err != nil {
		return nil, fmt.Errorf("generate publishable key: %w", err)
	}
	secretKey, err := generateOpaqueKey("sb_secret_")
	if err != nil {
		return nil, fmt.Errorf("generate secret key: %w", err)
	}

	// Serialise JWKS
	jwtKeysJSON, err := json.Marshal(jwksKeypair.Keys)
	if err != nil {
		return nil, fmt.Errorf("marshal JWT keys: %w", err)
	}
	jwksPublicJSON, err := json.Marshal(jwksPublic)
	if err != nil {
		return nil, fmt.Errorf("marshal JWKS: %w", err)
	}

	anonKey, err := generateLegacyJWT(jwtSecret, "anon")
	if err != nil {
		return nil, fmt.Errorf("generate anon key: %w", err)
	}

	serviceRoleKey, err := generateLegacyJWT(jwtSecret, "service_role")
	if err != nil {
		return nil, fmt.Errorf("generate service role key: %w", err)
	}

	return &GeneratedKeys{
		PublishableKey:           publishableKey,
		SecretKey:                secretKey,
		AnonKey:                  anonKey,
		AnonKeyAsymmetric:        anonJWT,
		ServiceRoleKey:           serviceRoleKey,
		ServiceRoleKeyAsymmetric: serviceJWT,
		JWTKeys:                  string(jwtKeysJSON),
		JWTJWKS:                  string(jwksPublicJSON),
	}, nil
}

func generateLegacyJWT(jwtSecret string, role string) (string, error) {
	header := map[string]any{
		"alg": "HS256",
		"typ": "JWT",
	}

	iat := time.Now().Unix()
	exp := iat + 5*365*24*3600

	payload := map[string]any{
		"role": role,
		"iss":  "supabase",
		"iat":  iat,
		"exp":  exp,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	b64Header := base64.RawURLEncoding.EncodeToString(headerJSON)
	b64Payload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := b64Header + "." + b64Payload

	mac := hmac.New(sha256.New, []byte(jwtSecret))
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + sig, nil
}

// ── EC key helpers ────────────────────────────────────────────────────────────
func GenerateJWTSecret() (string, error) {
	raw := make([]byte, 30)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate JWT secret: %w", err)
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}

func GenerateECKey() ([]byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate EC key: %w", err)
	}
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("marshal EC key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	}), nil
}

func parseECPrivateKey(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8
		parsed, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("parse EC private key: %w (also tried PKCS8: %v)", err, err2)
		}
		ec, ok := parsed.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("PEM key is not an EC key")
		}
		return ec, nil
	}
	return key, nil
}

func ecCoord(b *big.Int, curve elliptic.Curve) string {
	byteLen := (curve.Params().BitSize + 7) / 8
	buf := make([]byte, byteLen)
	b.FillBytes(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func ecPrivateKeyToJWK(key *ecdsa.PrivateKey, kid string) JWK {
	return JWK{
		Kty:    "EC",
		Kid:    kid,
		Use:    "sig",
		KeyOps: []string{"sign", "verify"},
		Alg:    "ES256",
		Ext:    true,
		Crv:    key.Curve.Params().Name,
		X:      ecCoord(key.PublicKey.X, key.Curve),
		Y:      ecCoord(key.PublicKey.Y, key.Curve),
		D:      ecCoord(key.D, key.Curve),
	}
}

func ecPublicKeyToJWK(key *ecdsa.PrivateKey, kid string) JWK {
	return JWK{
		Kty:    "EC",
		Kid:    kid,
		Use:    "sig",
		KeyOps: []string{"verify"},
		Alg:    "ES256",
		Ext:    true,
		Crv:    key.Curve.Params().Name,
		X:      ecCoord(key.PublicKey.X, key.Curve),
		Y:      ecCoord(key.PublicKey.Y, key.Curve),
	}
}

// ── JWT signing ───────────────────────────────────────────────────────────────

func signES256(key *ecdsa.PrivateKey, kid string, payload map[string]any) (string, error) {
	header := map[string]any{"alg": "ES256", "typ": "JWT", "kid": kid}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	b64Header := base64.RawURLEncoding.EncodeToString(headerJSON)
	b64Payload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := b64Header + "." + b64Payload

	digest := sha256.Sum256([]byte(signingInput))
	r, s, err := ecdsa.Sign(rand.Reader, key, digest[:])
	if err != nil {
		return "", fmt.Errorf("ECDSA sign: %w", err)
	}

	// IEEE P1363 format: r || s, each padded to curve byte length
	byteLen := (key.Curve.Params().BitSize + 7) / 8
	sig := make([]byte, 2*byteLen)
	r.FillBytes(sig[:byteLen])
	s.FillBytes(sig[byteLen:])

	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

// ── Opaque key generation ─────────────────────────────────────────────────────

func generateOpaqueKey(prefix string) (string, error) {
	raw := make([]byte, 17)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("random bytes: %w", err)
	}
	random := base64.RawURLEncoding.EncodeToString(raw)[:22]
	intermediate := prefix + random
	sum := sha256.Sum256([]byte(projectRef + "|" + intermediate))
	checksum := base64.RawURLEncoding.EncodeToString(sum[:])[:8]
	return intermediate + "_" + checksum, nil
}

// ── UUID ──────────────────────────────────────────────────────────────────────

func randomUUID() (string, error) {
	var uuid [16]byte
	if _, err := rand.Read(uuid[:]); err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:],
	), nil
}
