package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	DBUser     = "root"
	DBPassword = "root"
	DBName     = "root"
	DBHost     = "0.0.0.0"
	DBPort     = "5432"
	DBDriver   = "postgres"
	SSLMode    = "disable"
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
