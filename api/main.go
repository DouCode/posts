// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Mohamed Labouardy <mohamed@labouardy.com> https://labouardy.com
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"building-distributed-app-in-gin-chapter06/api/handlers"
	"building-distributed-app-in-gin-chapter06/api/middleware"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var authHandler *handlers.AuthHandler
var recipesHandler *handlers.RecipesHandler
var postController *handlers.PostController

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collectionPost := client.Database(os.Getenv("MONGO_DATABASE")).Collection("posts")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	postController = handlers.NewPostController(ctx, collectionPost, redisClient)

	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	router.POST("/api/auth/register", authHandler.RegisterHandler)
	router.POST("/api/auth/login", authHandler.SignInHandler)
	router.GET("/api/auth/info", authHandler.AuthMiddleware(), authHandler.Info)
	//router.POST("/refresh", authHandler.RefreshHandler)

	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
		authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	}

	postRoutes := router.Group("/")
	postRoutes.GET("/posts/:id", postController.Show)
	postRoutes.GET("/posts/page/list", postController.PageList)
	postRoutes.Use(authHandler.AuthMiddleware())
	postRoutes.POST("/posts", postController.Create)
	postRoutes.PUT("/posts/:id", postController.Update)
	postRoutes.DELETE("/posts/:id", postController.Delete)

	router.Run(":1016")
}
