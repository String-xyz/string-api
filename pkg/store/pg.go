package store

import (
	"fmt"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/env"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

var pgDB *sqlx.DB
var DBDriver = "postgres"

func strConnection() string {
	dbUser, _ := env.Get("DB_USERNAME")
	dbPassword, _ := env.Get("DB_PASSWORD")
	dbName, _ := env.Get("DB_NAME")
	dbHost, _ := env.Get("DB_HOST")
	dbPort, _ := env.Get("DB_PORT")

	var (
		DBUser     = dbUser
		DBPassword = dbPassword
		DBName     = dbName
		DBHost     = dbHost
		DBPort     = dbPort
	)

	var SSLMode string

	if libcommon.IsLocalEnv() {
		SSLMode = "disable"
	} else {
		SSLMode = "require"
	}

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
	sqltrace.Register(DBDriver, &pq.Driver{}, sqltrace.WithServiceName("string-api"))
	connection, err := sqlxtrace.Open(DBDriver, strConnection())
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
