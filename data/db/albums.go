package dbData

import (
	"database/sql"
	"fmt"
	"jolly-roger/web-service-exercise/defs"
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
	var album defs.Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumByID %d: no such album", id)
		}
		return album, fmt.Errorf("albumByID %d: %v", id, err)
	}

	return album, nil
}

func AddAlbum(alb defs.Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}
