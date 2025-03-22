package dbData

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	godotenv.Load()

	useDbPostgres := os.Getenv("USE_DB_POSTGRES")

	if useDbPostgres != "" {
		connectPostgres()
	} else {
		connectMysql()
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("DB connected!")
}

func connectPostgres() {
	connectRecordingsDb()

	if pingErr := db.Ping(); pingErr != nil {
		log.Println(pingErr)

		connPostgresDb := fmt.Sprintf(
			"postgres://%v:%v@%v/postgres?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
		)

		var err error
		db, err = sql.Open("postgres", connPostgresDb)
		if err != nil {
			log.Fatal(err)
		}

		if postgresDbPingErr := db.Ping(); postgresDbPingErr != nil {
			log.Fatal(postgresDbPingErr)
		}

		if createDbErr := createDb(); createDbErr != nil {
			log.Fatal(createDbErr)
		}

		if closeDbErr := db.Close(); closeDbErr != nil {
			log.Fatal(closeDbErr)
		}

		connectRecordingsDb()
	}
}

func connectRecordingsDb() {
	connRecordingsDb := fmt.Sprintf(
		"postgres://%v:%v@%v/recordings?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
	)

	var err error
	db, err = sql.Open("postgres", connRecordingsDb)
	if err != nil {
		log.Fatal(err)
	}
}

func connectMysql() {
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOST"),
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}
