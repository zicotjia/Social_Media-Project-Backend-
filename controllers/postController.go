package controllers

import (
	"backend/cloudinary"
	"backend/database"
	"backend/model"
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")

func cloudinaryUpload(post model.Post, r chan string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	uploadParam, err := cloudinary.Cloudi.Upload.Upload(ctx, post.File, uploader.UploadParams{UseFilename: api.Bool(true)})
	if err != nil {
		log.Fatal(err)
	}
	r <- uploadParam.SecureURL
}

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
		go cloudinaryUpload(post, cloudinaryURl)
		post.File = <-cloudinaryURl

		post.Likes = []model.Like{}
		post.Comments = []model.Comment{}
		post.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		resultInsertionNumber, insertErr := postCollection.InsertOne(ctx, post)
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

		cursor, err := postCollection.Find(ctx, bson.M{})
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
	UserID primitive.ObjectID `json:"_userid"`
	PostID primitive.ObjectID `json:"_postid"`
}

func LikePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LikeRequest
		var post model.Post
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var newLike model.Like
		defer cancel()

		fmt.Println("liked")

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(req.PostID)
		fmt.Println(req.UserID)
		newLike.User = req.UserID
		newLike.Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		filter := bson.M{"_id": req.PostID}
		update := bson.M{"$push": bson.M{"likes": newLike}}

		err := postCollection.FindOne(ctx, bson.M{"_id": req.PostID}).Decode(&post)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(post)
		_, err = postCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
			return
		}

	}
}
