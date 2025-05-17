package albumController

import (
	"database/sql"
	"models"
	"strconv"

	_ "github.com/lib/pq"
)

// query for all albums
func GetAllAlbums(db *sql.DB) ([]models.Album, error) {
	// slice of albums to hold results
	var albums []models.Album

	// prepare query
	stmt, err := db.Prepare("SELECT * FROM albums")
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // close statement when done

	// execute query
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close() // close rows when done

	// loop through rows, using Scan to assign values to struct fields
	for rows.Next() {
		var album models.Album
		err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	// check for errors from iteratings over rows
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// response with slice of albums
	return albums, nil
}

// query for album by ID
func GetAlbumById(db *sql.DB, id string) (models.Album, error) {
	// album to hold result
	var albumFound models.Album

	// convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return albumFound, err
	}

	// prepare query
	stmt, err := db.Prepare("SELECT * FROM albums WHERE id = $1")
	if err != nil {
		return albumFound, err
	}
	// close statement when done
	defer stmt.Close()

	// query
	row := stmt.QueryRow(idInt)

	// scan result into albumFound
	err = row.Scan(&albumFound.Id, &albumFound.Title, &albumFound.Artist, &albumFound.Price)
	if err != nil {
		return albumFound, err
	}

	// return album found
	return albumFound, nil
}

// exec to add an album
func AddAlbum(db *sql.DB, albumToAdd models.Album) (int64, error) {
	// prepare query with RETURNING id
	stmt, err := db.Prepare("INSERT INTO albums(title, artist, price) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}
	// close statement when done
	defer stmt.Close()

	// use QueryRow to get the returned id
	var id int64
	err = stmt.QueryRow(albumToAdd.Title, albumToAdd.Artist, albumToAdd.Price).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// exec to update an album
func UpdateAlbum(db *sql.DB, albumToUpdate models.Album) (int64, error) {
	// prepare query
	stmt, err := db.Prepare("UPDATE albums SET title = $1, artist = $2, price = $3 WHERE id = $4")
	if err != nil {
		return 0, err
	}
	// close statement when done
	defer stmt.Close()

	// exec query
	result, err := stmt.Exec(albumToUpdate.Title, albumToUpdate.Artist, albumToUpdate.Price, albumToUpdate.Id)
	if err != nil {
		return 0, err
	}

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }

	return rowsAffected, nil
}

// exec to delete an album
func DeleteAlbum(db *sql.DB, id string) (int64, error) {
	// prepare query
	stmt, err := db.Prepare("DELETE FROM albums WHERE id = $1")
	if err != nil {
		return 0, err
	}
	// close statement when done
	defer stmt.Close()

	// exec query
	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }

	return rowsAffected, nil
}