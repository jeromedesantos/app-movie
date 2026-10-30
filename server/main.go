package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	controller "github.com/jeromedesantos/app-movie/server/controllers"
	"github.com/jeromedesantos/app-movie/server/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	var (
		client  *mongo.Client
		router  *gin.Engine
		err     error
	)

	client = database.Connect()
	router = gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, MagicStreamMovies!")
	})
	
	router.GET("/movies", controller.GetMovies(client))
	router.GET("/movie/:imdb_id", controller.GetMovie(client))

	err = router.Run()
	if err != nil {
		fmt.Println("Failed to start server", err)
	}
}