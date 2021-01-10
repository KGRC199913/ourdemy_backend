package routes

import (
	"github.com/gin-gonic/gin"
)

func HomeRoutes(route *gin.Engine) {
	homeRoutesGroup := route.Group("/home")
	{
		homeRoutesGroup.GET("/highlights", func(c *gin.Context) {

		})
	}
}
