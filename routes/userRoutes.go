package routes

import (
	"github.com/falasefemi2/golang-jwt-project/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controllers.Signup())
	incomingRoutes.POST("users/login", controllers.Login())
}

// package routes

// import (
//     "github.com/gin-gonic/gin"
// )

// func AuthRoutes(incomingRoutes *gin.Engine) {
//     incomingRoutes.POST("users/signup", controller.Signup())
//     incomingRoutes.POST("users/login", controller.Login())
// }
