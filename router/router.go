package router

import (
	controller "backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.SignUp())
	incomingRoutes.POST("/users/login", controller.Login())
}

func PostRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/post/addpost", controller.FileUpload())
	incomingRoutes.GET("/post/getpost", controller.GetAllPost())
}

func LikeRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.PATCH("/post/like", controller.LikePost())
}
