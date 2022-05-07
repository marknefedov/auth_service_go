package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c *gin.Context) {
	var requestBody RegisterRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request body"})
		return
	}

	var user User
	// #TODO: This whole thing should be in transaction
	err := MONGO_USERS.FindOne(context.TODO(), bson.D{{"email", requestBody.Email}}).Decode(&user)
	hash, _ := HashPassword(requestBody.Password)
	if err == mongo.ErrNoDocuments {
		user := User{
			Uuid:     uuid.New(),
			FullName: requestBody.FullName,
			Email:    requestBody.Email,
			Password: hash,
			Sessions: []Session{},
		}
		session := Session{
			SessionName: "Registration session",
			OpaqueToken: xid.New().String(),
			LastUsed:    time.Now(),
		}
		user.Sessions = append(user.Sessions, session)

		_, err := MONGO_USERS.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, SessionResponse{
			SessionUUID: session.OpaqueToken,
		})
		return
	} else if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusConflict, gin.H{"error": "Email already taken"})
}

func Login(c *gin.Context) {
	var requestBody LoginRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request body"})
		return
	}

	var user User
	err := MONGO_USERS.FindOne(context.TODO(), bson.D{{"email", requestBody.Email}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password invalid"})
		return
	} else if err != nil {
		log.Fatal(err)
	}
	if CheckPasswordHash(requestBody.Password, user.Password) {
		user_agent := c.Request.Header.Get("User-Agent")
		session := Session{
			SessionName: user_agent,
			OpaqueToken: xid.New().String(),
			LastUsed:    time.Now(),
		}
		_, err := MONGO_USERS.UpdateOne(context.TODO(), bson.D{{"email", requestBody.Email}}, bson.D{{"$push", bson.D{{"sessions", session}}}})
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, SessionResponse{
			SessionUUID: session.OpaqueToken,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password invalid"})
		return
	}
}

func GetAccsessToken(c *gin.Context) {
	var requestBody AccsessTokenRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user User
	err := MONGO_USERS.FindOne(context.TODO(), bson.D{{"sessions.opaque_token", requestBody.SessionUUID}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalid"})
		return
	} else if err != nil {
		log.Fatal(err)
	}

	_, err = MONGO_USERS.UpdateOne(context.TODO(), bson.D{{"sessions.opaque_token", requestBody.SessionUUID}}, bson.D{{"$set", bson.D{{"sessions.$.last_used", time.Now()}}}})
	if err != nil {
		log.Fatal(err)
	}

	token := CreateJWTES384(user.Uuid.String())

	c.JSON(http.StatusOK, AccsessTokenResponse{
		AccsessToken: token,
	})

}
