package routes

import (
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeRoutes(route *gin.Engine) {
	homeRoutesGroup := route.Group("/home")
	{
		homeRoutesGroup.GET("/highlights", func(c *gin.Context) {
			courses, err := models.GetTop4HighlightCourse()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, courses)
		})

		homeRoutesGroup.GET("/mostWatch", func(c *gin.Context) {
			courses, err := models.GetTop10MostWatchCourse()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, courses)
		})

		homeRoutesGroup.GET("/mostReg", func(c *gin.Context) {
			cats, err := models.GetAllMostRegisterCategory()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, cats)
		})

		homeRoutesGroup.GET("/newest", func(c *gin.Context) {
			courses, err := models.GetTop10NewestCourse()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, courses)
		})
	}
}
