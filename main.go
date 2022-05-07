package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"	
)

var MONGO_USERS *mongo.Collection

func main() {
	mogno_uri := os.Getenv("AUTH_SERVICE_MONGO_URI")
	db_name := os.Getenv("AUTH_SERVICE_MONGO_DB_NAME")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mogno_uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	MONGO_USERS = client.Database(db_name).Collection("Users")

	r := gin.Default()
	r.POST("/api/register", Register)
	r.POST("/api/login", Login)
	r.POST("/api/get_accsess_token", GetAccsessToken)
	r.Run() // listen and serve on 0.0.0.0:8080
}
