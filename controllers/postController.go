package controllers

import (
	"backend/cloudinary"
	"backend/database"
	"backend/model"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var PostCollection *mongo.Collection = database.OpenCollection(database.Client, "post")

func FileUpload() gin.HandlerFunc {
	return func(c *gin.Context) {

		var post model.Post
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cloudinaryURl := make(chan string)
		go cloudinary.CloudinaryUpload(post.File, cloudinaryURl)
		post.File = <-cloudinaryURl

		post.Likes = []model.Like{}
		post.Comments = []model.Comment{}
		post.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		resultInsertionNumber, insertErr := PostCollection.InsertOne(ctx, post)
		if insertErr != nil {
			msg := fmt.Sprintf("Post item was not added")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func GetAllPost() gin.HandlerFunc {
	return func(c *gin.Context) {

		var posts []model.Post
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"created_at", -1}})
		cursor, err := PostCollection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			log.Fatal(err)
		}

		for cursor.Next(ctx) {
			var post model.Post
			if err = cursor.Decode(&post); err != nil {
				log.Fatal(err)
			}

			posts = append(posts, post)
		}

		c.JSON(http.StatusOK, posts)

	}
}

type LikeRequest struct {
	UserId   primitive.ObjectID `json:"_userid"`
	Username string             `json:"username"`
	PostId   primitive.ObjectID `json:"_postid"`
}

func LikePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LikeRequest
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var newLike model.Like
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newLike.UserId = req.UserId
		newLike.Username = req.Username
		newLike.Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		filter := bson.M{"_id": req.PostId}
		update := bson.M{"$push": bson.M{"likes": newLike}}

		_, err := PostCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
			return
		}

	}
}

func UnlikePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LikeRequest
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		fmt.Println("unliked")

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"_id": req.PostId}
		update := bson.M{"$pull": bson.M{"likes": bson.M{"userid": req.UserId}}}

		_, err := PostCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
			return
		}

	}
}
