package settings

import (
	"github.com/go-pg/pg/v10"
	"os"
)

const DB_NAME string = "unicornkitchen"

func GetDatabaseOptions() *pg.Options {
	res := &pg.Options{
		User:     envVar("DB_USER", "postgres"),
		Password: envVar("DB_PASS", "abc123"), // Or a secrets manager...
		Database: envVar("DB_NAME", DB_NAME),
		Addr:     os.Getenv("DB_ADDRESS"), // Empty string is default..
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
