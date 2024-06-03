package router

import (
	"finalapi/controllers/photocontroller"
	"finalapi/controllers/usercontroller"
	"finalapi/middlewares"

	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	//Photos Endpoint
	photos := router.Group("/photos", middlewares.JWTMiddleware())
	{
		photos.GET("/", photocontroller.Index)
		photos.GET("/:id", photocontroller.Show)
		photos.POST("/", photocontroller.Create)
		photos.PUT("/:id", photocontroller.Update)
		photos.DELETE("/:id", photocontroller.Delete)
	}

	users := router.Group("/users")
	{
		users.POST("/register", usercontroller.Register)
		users.GET("/login", usercontroller.Login)
		users.PUT("/:userId", middlewares.JWTMiddleware(), usercontroller.UpdateUser)
		users.DELETE("/:userId", middlewares.JWTMiddleware(), usercontroller.DeleteUser)
		users.GET("/logout", middlewares.JWTMiddleware(), usercontroller.Logout)
	}

	return router
}
