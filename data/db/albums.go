package dbData

import (
	"database/sql"
	"fmt"
	"jolly-roger/web-service-exercise/defs"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetAlbums() ([]defs.Album, error) {
	var albums []defs.Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("no albums in db: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var album defs.Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("can not read album from db response: %v", err)
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can not read albums from response: %v", err)
	}

	return albums, nil
}

func GetAlbumByID(id int64) (defs.Album, error) {
	godotenv.Load()

	var (
		useDbPostgres = os.Getenv("USE_DB_POSTGRES")
		album         defs.Album
		row           *sql.Row
	)

	if useDbPostgres != "" {
		row = db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	} else {
		row = db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	}

	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumByID %d: no such album", id)
		}

		return album, fmt.Errorf("albumByID %d: %v", id, err)
	}

	return album, nil
}

func AddAlbum(alb defs.Album) (int64, error) {
	godotenv.Load()

	var (
		useDbPostgres = os.Getenv("USE_DB_POSTGRES")
		id            int64
	)

	if useDbPostgres != "" {
		row := db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) returning id", alb.Title, alb.Artist, alb.Price)

		if err := row.Scan(&id); err != nil {
			return id, fmt.Errorf("addAlbum: %v", err)
		}
	} else {
		var (
			err    error
			result sql.Result
		)

		result, err = db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)

		if err != nil {
			return 0, fmt.Errorf("addAlbum: %v", err)
		}

		id, err = result.LastInsertId()
		if err != nil {

			log.Panicln(err)

			return 0, fmt.Errorf("addAlbum: %v", err)
		}
	}

	return id, nil
}
