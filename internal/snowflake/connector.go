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
	pass := os.Getenv("SNOWFLAKE_PASSWORD")
	acct := os.Getenv("SNOWFLAKE_ACCOUNT")
	dbname := os.Getenv("SNOWFLAKE_DATABASE")
	schema := os.Getenv("SNOWFLAKE_SCHEMA")
	wh := os.Getenv("SNOWFLAKE_WAREHOUSE")
	role := os.Getenv("SNOWFLAKE_ROLE")

	if user == "" || acct == "" || dbname == "" || schema == "" || wh == "" || role == "" {
		return nil, fmt.Errorf("missing required env vars")
	}

	dsn := fmt.Sprintf("%s:%s@%s/%s/%s?warehouse=%s",
		user, pass, acct, dbname, schema, wh,
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
