package store

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DBUser     = os.Getenv("DB_USERNAME")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName     = os.Getenv("DB_NAME")
	DBHost     = os.Getenv("DB_HOST")
	DBPort     = os.Getenv("DB_PORT")
	DBDriver   = "postgres"
	SSLMode    = "disable" // require
)

var pgDB *sqlx.DB

func strConnection() string {
	str := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		DBHost,
		DBPort,
		DBUser,
		DBName,
		DBPassword,
		SSLMode,
	)
	return str
}

func MustNewPG() *sqlx.DB {
	if pgDB != nil {
		return pgDB
	}
	connection, err := sqlx.Open(DBDriver, strConnection())
	if err != nil {
		panic(err)
	}
	if err := connection.Ping(); err != nil {
		panic(err)
	}

	pgDB = connection
	return pgDB
}

// GetPGInstance returns the already initialized postgres instance
func GetPGInstance() *sqlx.DB {
	return pgDB
}
