package memoryData

import (
	"fmt"
	"jolly-roger/web-service-exercise/defs"
)

var albums = []defs.Album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func GetAlbums() []defs.Album {
	return albums
}

func AddAlbum(newAlbum defs.Album) {
	albums = append(albums, newAlbum)
}

func GetAlbumByID(id int64) (defs.Album, error) {
	for _, a := range albums {
		if a.ID == id {
			return a, nil
		}
	}
	return defs.Album{}, fmt.Errorf("album with id %v not found", id)
}
