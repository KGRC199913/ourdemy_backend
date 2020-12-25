package routes

import (
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LecturerRoutes(route *gin.Engine) {
	lecRoutesGroup := route.Group("/lecturers", middlewares.Authenticate, middlewares.LecturerAuthenticate)
	{
		lecRoutesGroup.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "kay",
			})
		})
		//lecRoutesGroup.POST("/create", func(c *gin.Context) {
		//	var lecturer models.Lecturer
		//	if err := c.ShouldBind(&lecturer); err != nil {
		//		c.JSON(http.StatusBadRequest, gin.H{
		//			"error": err.Error(),
		//		})
		//		return
		//	}
		//
		//	if err := lecturer.Save(); err != nil {
		//		c.JSON(http.StatusBadRequest, gin.H{
		//			"error": err.Error(),
		//		})
		//		return
		//	}
		//
		//	c.JSON(http.StatusOK, lecturer)
		//})
	}

}
