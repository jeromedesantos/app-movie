package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jeromedesantos/app-movie/server/database"
	"github.com/jeromedesantos/app-movie/server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var validate = validator.New()

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err				error
			ctx       		context.Context
			cancel    		context.CancelFunc
			movies			[]models.Movie
			cursor			*mongo.Cursor
			movieCollection	*mongo.Collection
		)

		ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieCollection = database.OpenCollection("movies", client)

		cursor, err = movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch movies.",
			})
			return
		}
		
		defer cursor.Close(ctx)

		err = cursor.All(ctx, &movies)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to decode movies.",
			})
			return
		}

		c.JSON(http.StatusOK, movies)
	}

}

func GetMovie(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
		var (
			err 			error
			ctx 			context.Context
			cancel			context.CancelFunc
			movie			models.Movie
			movieCollection	*mongo.Collection
			movieID			string
		)

		ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()
		
		movieID = c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Movie ID is required.",
			})
			return
		}

		movieCollection = database.OpenCollection("movies", client)

		err =  movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Movie not found.",
			})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
		var (
			err 			error
			ctx 			context.Context
			cancel			context.CancelFunc
			movie			models.Movie
			movieCollection	*mongo.Collection
			movieID			string
		)

		ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		err = c.ShouldBindJSON(&movie)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input.",
			})
			return
		}
		
		err = validate.Struct(&movie)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Validation failed.",
				"details": err.Error(),
			})
			return
		}

	}
}

/*
CATATAN PENGGUNAAN CONTEXT & BEST PRACTICES

1. Kenapa Membungkus *gin.Context dengan context.WithTimeout?
   - Timeout Protection: Mencegah query database menggantung (hang) selamanya jika
     terjadi network down/latensi tinggi, dengan membatalkan operasi saat batas waktu habis.
   - Resource Cleanup: 'defer cancel()' memastikan alokasi timer internal context
     segera dibersihkan dari memori begitu handler selesai.

2. Apakah Penggunaan Context Wajib?
   - Wajib di MongoDB Driver: Semua operasi I/O (Find, FindOne, Insert, dll.)
     mewajibkan context.Context sebagai argumen pertama untuk kontrol pembatalan operasi.

3. Best Practices:
   - Durasi Masuk Akal: Gunakan timeout 5–10 detik untuk operasi baca standar.
   - Return Saat Error: Selalu sertakan 'return' setelah response error (c.JSON)
     agar tidak mengeksekusi kode berikutnya yang bisa memicu nil pointer panic.
   - Defer Close Cursor: Panggil 'defer cursor.Close(ctx)' tepat setelah Find berhasil
     guna mencegah kebocoran memori pada pool koneksi.

*/
