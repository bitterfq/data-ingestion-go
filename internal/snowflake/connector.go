package snowflake

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/snowflakedb/gosnowflake"
	_ "github.com/snowflakedb/gosnowflake"
)

// getPath resolves a path relative to this package's directory if not absolute.
func getPath(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	_, b, _, _ := runtime.Caller(0) // b = full path of this source file
	base := filepath.Dir(b)         // directory of this file (internal/snowflake)
	return filepath.Join(base, rel)
}

func loadPrivateKey(path string, passphrase []byte) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in key")
	}

	var der []byte
	if x509.IsEncryptedPEMBlock(block) {
		der, err = x509.DecryptPEMBlock(block, passphrase)
		if err != nil {
			return nil, fmt.Errorf("decrypt key: %w", err)
		}
	} else {
		der = block.Bytes
	}

	parsed, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS#8: %w", err)
	}

	priv, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return priv, nil
}

func NewClient(env string) (*sql.DB, error) {
	if env == "" {
		env = ".env"
	}

	err := godotenv.Load(env)

	if err != nil {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	cfg := gosnowflake.Config{
		Account:       os.Getenv("SNOWFLAKE_ACCOUNT"),
		User:          os.Getenv("SNOWFLAKE_USER"),
		Database:      os.Getenv("SNOWFLAKE_DATABASE"),
		Schema:        os.Getenv("SNOWFLAKE_SCHEMA"),
		Warehouse:     os.Getenv("SNOWFLAKE_WAREHOUSE"),
		Role:          os.Getenv("SNOWFLAKE_ROLE"),
		Authenticator: gosnowflake.AuthTypeJwt,
	}

	keyPath := getPath(os.Getenv("SNOWFLAKE_PRIVATE_KEY_FILE"))
	privateKey, err := loadPrivateKey(keyPath,
		[]byte(os.Getenv("SNOWFLAKE_PRIVATE_KEY_PASSPHRASE")))
	if err != nil {
		return nil, err
	}
	cfg.PrivateKey = privateKey

	dsn, _ := gosnowflake.DSN(&cfg)
	db, _ := sql.Open("snowflake", dsn)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
