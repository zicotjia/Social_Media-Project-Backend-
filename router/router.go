package router

import (
	controller "backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.SignUp())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.GET("/users/getusers", controller.GetAllUsers())
	incomingRoutes.PATCH("/users/follow", controller.FollowUser())
	incomingRoutes.PATCH("/users/unfollow", controller.UnFollowUser())
	incomingRoutes.PATCH("/users/edit", controller.EditUser())
	incomingRoutes.PATCH("/users/editpic", controller.EditProfilePic())
	incomingRoutes.DELETE("/users/delete", controller.DeleteUser())
}

func PostRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/post/addpost", controller.FileUpload())
	incomingRoutes.GET("/post/getpost", controller.GetAllPost())
	incomingRoutes.PATCH("/post/edit", controller.EditPost())
	incomingRoutes.DELETE("/post/delete", controller.DeletePost())
}

func LikeRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.PATCH("/post/like", controller.LikePost())
	incomingRoutes.PATCH("/post/unlike", controller.UnlikePost())
}
