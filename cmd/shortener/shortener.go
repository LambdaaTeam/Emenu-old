package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"

	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	possibleChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortSize     = 8
)

type ShortenedUrl struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Url   string             `json:"url"`
	Short string             `json:"short"`
}

func randomString(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = possibleChars[rand.Intn(len(possibleChars))]
	}

	return string(b)
}

func generateShortUrl(ctx context.Context) string {
	for {
		short := randomString(shortSize)

		var result ShortenedUrl

		err := database.DB.Collection("urls").FindOne(ctx, bson.M{"short": short}).Decode(&result)

		if err != nil {
			return short
		}
	}
}

func init() {
	godotenv.Load(".env")

	dbName := "shortener-dev"

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
		dbName = "shortener"
	}

	database.DB = database.Connect(dbName)

	_, err := database.DB.Collection("urls").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{
			"short": 1,
		},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/", func(c *gin.Context) {
		ctx := c.Request.Context()

		var url struct {
			Url string `json:"url"`
		}

		if err := c.ShouldBindJSON(&url); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shortenedUrl := ShortenedUrl{
			ID:    primitive.NewObjectID(),
			Url:   url.Url,
			Short: generateShortUrl(ctx),
		}

		_, err := database.DB.Collection("urls").InsertOne(ctx, shortenedUrl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, shortenedUrl)
	})

	r.GET("/:short", func(c *gin.Context) {
		ctx := c.Request.Context()

		short := c.Param("short")

		var shortenedUrl ShortenedUrl

		err := database.DB.Collection("urls").FindOne(ctx, bson.M{"short": short}).Decode(&shortenedUrl)

		if err != nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}

		c.Redirect(http.StatusMovedPermanently, shortenedUrl.Url)
	})

	r.Run()
}
