package main

import (
	"context"
	dbData "jolly-roger/web-service-exercise/data/db"
	memoryData "jolly-roger/web-service-exercise/data/memory"
	"jolly-roger/web-service-exercise/defs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var _, isLambda = os.LookupEnv("AWS_LAMBDA_RUNTIME_API")

var ginLambda *ginadapter.GinLambda

func init() {
	godotenv.Load()
	useDb := os.Getenv("USE_DB")

	if useDb != "" {
		dbData.Init()
	}

	log.Println("Is DB used: ", useDb)
	log.Printf("Gin cold start")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.GET("/admin/seed-data", getAdminSeedData)

	if isLambda {
		ginLambda = ginadapter.New(router)
	} else {
		router.Run(":8080")
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if isLambda {
		lambda.Start(Handler)
	}
}

func getAlbums(c *gin.Context) {
	useDb := os.Getenv("USE_DB")

	var albums []defs.Album
	var err error

	if useDb != "" {
		albums, err = dbData.GetAlbums()
	} else {
		albums = memoryData.GetAlbums()
	}

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no albumes"})
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	useDb := os.Getenv("USE_DB")

	var newAlbum defs.Album
	var err error
	var id int64

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	if useDb != "" {
		id, err = dbData.AddAlbum(newAlbum)
		newAlbum.ID = id
	} else {
		memoryData.AddAlbum(newAlbum)
	}

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not add albume"})
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	useDb := os.Getenv("USE_DB")

	id, idErr := strconv.ParseInt(c.Param("id"), 10, 64)
	if idErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "incorrect id"})
		return
	}

	var a defs.Album
	var err error

	if useDb != "" {
		a, err = dbData.GetAlbumByID(id)
	} else {
		a, err = memoryData.GetAlbumByID(id)
	}

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "albume not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

func getAdminSeedData(c *gin.Context) {
	err := dbData.SeedData()
	if err != nil {
		log.Printf("db error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "can not seed data"})
		return
	}

	c.IndentedJSON(http.StatusOK, "data is seeded")
}
