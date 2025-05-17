package main

import (
	albumController "controllers"
	"database/sql"
	"fmt"
	"log"
	"models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// define database connection string
	connStr := "user=s password=s dbname=testing host=localhost port=5432 sslmode=disable"

	// get a database handler
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// check database connection
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("connected to database")

	// --- setting up router ---
	router := gin.Default()

	// route to get all albums
	router.GET("/albums/all", func(c *gin.Context) {
		albums, err := albumController.GetAllAlbums(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, albums)
	})

	// route to get a specific album by id
	router.GET("/albums/:id", func(c *gin.Context) {
		album, err := albumController.GetAlbumById(db, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, album)
	})

	// route to create a new album
	router.POST("/albums", func(c *gin.Context) {
		var newAlbum models.Album

		err := c.ShouldBindJSON(&newAlbum)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		newAlbumId, err := albumController.AddAlbum(db, newAlbum)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, newAlbumId)
	})

	// route to update an album
	router.PUT("/albums/edit", func(c *gin.Context) {
		var updatedAlbum models.Album

		err := c.ShouldBindJSON(&updatedAlbum)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		rowsAffected, err := albumController.UpdateAlbum(db, updatedAlbum)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"rows Updated": rowsAffected})
	})

	// route to delete an album
	router.DELETE("/albums/delete/:id", func(c *gin.Context) {
		rowsAffected, err := albumController.DeleteAlbum(db, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"rows Deleted": rowsAffected})
	})

	router.Run("localhost:8080")
}
