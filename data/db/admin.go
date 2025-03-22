package dbData

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func createDb() error {
	if _, err := db.Exec("create database recordings"); err != nil {
		return fmt.Errorf("can not create database recordings: %v", err)
	}
	return nil
}

func SeedData() error {
	godotenv.Load()

	useDbPostgres := os.Getenv("USE_DB_POSTGRES")

	if useDbPostgres != "" {
		if err := seedDataPostgres(); err != nil {
			return err
		}
	} else {
		if err := seeDataMysql(); err != nil {
			return err
		}
	}

	return nil
}

func seedDataPostgres() error {
	if _, err := db.Exec("DROP TABLE IF EXISTS album"); err != nil {
		return fmt.Errorf("can not drop table album: %v", err)
	}

	if _, err := db.Exec("CREATE TABLE album (" +
		"id         serial primary key," +
		"title      VARCHAR(128) NOT NULL," +
		"artist     VARCHAR(255) NOT NULL," +
		"price      DECIMAL(5,2) NOT NULL" +
		")"); err != nil {
		return fmt.Errorf("can not create table album: %v", err)
	}

	if _, err := db.Exec("INSERT INTO album" +
		"(title, artist, price)" +
		"VALUES" +
		"('Blue Train', 'John Coltrane', 56.99)," +
		"('Giant Steps', 'John Coltrane', 63.99)," +
		"('Jeru', 'Gerry Mulligan', 17.99)," +
		"('Sarah Vaughan', 'Sarah Vaughan', 34.98)"); err != nil {
		return fmt.Errorf("can not insert data into table album: %v", err)
	}

	return nil
}

func seeDataMysql() error {
	if _, err := db.Exec("drop database if exists recordings"); err != nil {
		return fmt.Errorf("can not drop database recordings: %v", err)
	}

	if err := createDb(); err != nil {
		return err
	}

	if _, err := db.Exec("use recordings"); err != nil {
		return fmt.Errorf("can not use database recordings: %v", err)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS album"); err != nil {
		return fmt.Errorf("can not drop table album: %v", err)
	}

	if _, err := db.Exec("CREATE TABLE album (" +
		"id         INT AUTO_INCREMENT NOT NULL," +
		"title      VARCHAR(128) NOT NULL," +
		"artist     VARCHAR(255) NOT NULL," +
		"price      DECIMAL(5,2) NOT NULL," +
		"PRIMARY KEY (`id`))"); err != nil {
		return fmt.Errorf("can not create table album: %v", err)
	}

	if _, err := db.Exec("INSERT INTO album" +
		"(title, artist, price)" +
		"VALUES" +
		"('Blue Train', 'John Coltrane', 56.99)," +
		"('Giant Steps', 'John Coltrane', 63.99)," +
		"('Jeru', 'Gerry Mulligan', 17.99)," +
		"('Sarah Vaughan', 'Sarah Vaughan', 34.98)"); err != nil {
		return fmt.Errorf("can not insert data into table album: %v", err)
	}

	return nil
}
