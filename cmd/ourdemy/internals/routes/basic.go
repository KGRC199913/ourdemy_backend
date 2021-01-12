package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BasicRoutes(route *gin.Engine) {
	route.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome to ourdemy api",
		})
	})
	route.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"authors": gin.H{
				"1": "Pham Hoang Anh Tuan",
				"2": "Nguyen Duong Tri",
				"3": "Cao Dinh Vi",
			},
			"license": "MIT License",
		})
	})
}
