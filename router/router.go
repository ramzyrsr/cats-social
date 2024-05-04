package router

import (
	"cats-social/controllers"
	"cats-social/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Cats Social",
		})
	})

	database.ConnectDatabase()

	router.POST("/v1/user/register", controllers.RegisterUser)
	router.POST("/v1/user/login", controllers.Login)

	authorized := router.Group("/")

	{
		authorized.POST("/v1/cat", controllers.CreateCat)
		authorized.DELETE("/v1/cat/:id", controllers.DeleteCats)
		authorized.POST("/v1/cat/match/reject", controllers.RejectCat)
	}

	return router
}
