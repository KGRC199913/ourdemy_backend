package routes

import (
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CourseRoutes(route *gin.Engine) {
	courseRoutesGroup := route.Group("/course")
	{
		courseRoutesGroup.GET("/test/create", func(c *gin.Context) {
			testCourse := models.NewCourse("123", "123", "123", "ABC GHI A ASA", "ABCD", "Neko Neko Cat Kitty", 25.0, 10)
			err := testCourse.Save()
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, testCourse)
			}
		})
	}
}
