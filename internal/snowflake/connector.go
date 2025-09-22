package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/snowflakedb/gosnowflake"
)

func NewClient() (*sql.DB, error) {
	_ = godotenv.Load("../../.env") // load .env file if present

	user := os.Getenv("SNOWFLAKE_USER")
	acct := os.Getenv("SNOWFLAKE_ACCOUNT")
	dbname := os.Getenv("SNOWFLAKE_DATABASE")
	schema := os.Getenv("SNOWFLAKE_SCHEMA")
	wh := os.Getenv("SNOWFLAKE_WAREHOUSE")
	role := os.Getenv("SNOWFLAKE_ROLE")
	auth := os.Getenv("SNOWFLAKE_AUTH")

	if user == "" || acct == "" || dbname == "" || schema == "" || wh == "" || role == "" || auth == "" {
		return nil, fmt.Errorf("missing required env vars")
	}

	dsn := fmt.Sprintf("%s:@%s/%s/%s?warehouse=%s&role=%s&authenticator=%s",
		user, acct, dbname, schema, wh, role, auth,
	)

	db, err := sql.Open("snowflake", dsn)

	if err != nil {
		log.Fatal("open failed:", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
