package settings

import (
	//	"github.com/go-pg/pg/v10"
	"os"
)

const DB_NAME string = "unicornkitchen"

type DBOptions struct {
	User     string
	Password string
	Database string
	Address  string // Full address, including the port.
}

func GetDatabaseOptions() *DBOptions {
	res := &DBOptions{
		User:     envVar("DB_USER", "postgres"),
		Password: envVar("DB_PASS", "abc123"), // Or a secrets manager...
		Database: envVar("DB_NAME", DB_NAME),
		Address:  envVar("DB_ADDRESS", "localhost"),
	}

	return res
}

func envVar(varName string, fallback string) string {
	res := os.Getenv(varName)
	if res == "" {
		res = fallback
	}
	return res
}

// res := pg.Connect(&pg.Options{
// 	// NOTE: This is dev data.  In real life, one would use 'secrets' from their container or server provider.
// 	User:     "postgres",
// 	Password: "abc123",
// 	Database: DB_NAME,
// })
