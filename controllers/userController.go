package controllers

import (
	"backend/cloudinary"
	"backend/database"
	helper "backend/helpers"
	"backend/model"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user model.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := HashPassword(user.Password)
		user.Password = password

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		user.Follower = []primitive.ObjectID{}
		user.Following = []primitive.ObjectID{}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.First_name, user.Last_name, user.User_id)
		user.Token = token
		user.Refresh_token = refreshToken

		resultInsertionNumber, insertErr := UserCollection.InsertOne(ctx, user)

		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

//Login is the api used to get a single user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user model.User
		var foundUser model.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_id)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser)

	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var allUser []model.User
		defer cancel()

		cursor, err := UserCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}

		for cursor.Next(ctx) {
			var user model.User
			if err = cursor.Decode(&user); err != nil {
				log.Fatal(err)
			}
			allUser = append(allUser, user)
		}

		if err = cursor.Close(ctx); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUser)
	}
}

type FollowRequest struct {
	Follower  primitive.ObjectID `json:"follower"`
	Following primitive.ObjectID `json:"following"`
}

func FollowUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var req FollowRequest
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			log.Fatal(err)
		}

		fmt.Println(req)

		filter1 := bson.M{"_id": req.Follower}
		update1 := bson.M{"$push": bson.M{"following": req.Following}}

		filter2 := bson.M{"_id": req.Following}
		update2 := bson.M{"$push": bson.M{"follower": req.Follower}}

		_, err := UserCollection.UpdateOne(ctx, filter1, update1)
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err = UserCollection.UpdateOne(ctx, filter2, update2)
		if err != nil {
			log.Fatal(err)
			return
		}

	}
}

func UnFollowUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var req FollowRequest
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			log.Fatal(err)
		}

		filter1 := bson.M{"_id": req.Follower}
		update1 := bson.M{"$pull": bson.M{"following": req.Following}}

		filter2 := bson.M{"_id": req.Following}
		update2 := bson.M{"$pull": bson.M{"following": req.Follower}}

		_, err := UserCollection.UpdateOne(ctx, filter1, update1)
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err = UserCollection.UpdateOne(ctx, filter2, update2)
		if err != nil {
			log.Fatal(err)
			return
		}

	}
}

type ProfileChangeRequest struct {
	UserId primitive.ObjectID `json:"user_id"`
	Bio    string             `json:"bio"`
}

func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var newProfile ProfileChangeRequest
		var updatedUser model.User
		defer cancel()

		if err := c.BindJSON(&newProfile); err != nil {
			log.Fatal(err)
		}

		filter := bson.M{"_id": newProfile.UserId}
		update := bson.M{"$set": bson.M{"bio": newProfile.Bio}}
		_, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = UserCollection.FindOne(ctx, bson.M{"_id": newProfile.UserId}).Decode(&updatedUser)
		if err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, updatedUser)
	}
}

type ProfilePicChangeRequest struct {
	UserId        primitive.ObjectID `json:"user_id"`
	ProfilePicUrl string             `json:"profilepicurl"`
}

func EditProfilePic() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ProfilePicChangeRequest
		var updatedUser model.User
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cloudinaryURl := make(chan string)
		go cloudinary.CloudinaryUpload(req.ProfilePicUrl, cloudinaryURl)
		req.ProfilePicUrl = <-cloudinaryURl

		filter := bson.M{"_id": req.UserId}
		update := bson.M{"$set": bson.M{"profilepicurl": req.ProfilePicUrl}}

		_, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			msg := fmt.Sprintf("Post item was not added")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		err = UserCollection.FindOne(ctx, bson.M{"_id": req.UserId}).Decode(&updatedUser)
		if err != nil {
			log.Fatal(err)
			return
		}

		filter = bson.M{"user._id": req.UserId}
		update = bson.M{"$set": bson.M{"user": updatedUser}}
		_, err = PostCollection.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, updatedUser)
	}
}

type DeleteUserRequest struct {
	UserId primitive.ObjectID `json:"userid"`
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DeleteUserRequest
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&req); err != nil {
			log.Fatal(err)
		}

		filter := bson.M{"_id": req.UserId}

		_, err := UserCollection.DeleteOne(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}

		filter = bson.M{"user._id": req.UserId}

		_, err = PostCollection.DeleteMany(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
	}
}
