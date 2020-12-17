package routes

import (
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestRoutes(route *gin.Engine) {
	testRoutesGroup := route.Group("/test")
	{
		testRoutesGroup.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "OK!")
		})
		testRoutesGroup.GET("/create", func(c *gin.Context) {
			testUser := models.NewUser("ABC", "ABC@Meow.com", "123")
			//err := testUser.Create()
			//
			err := testUser.Save()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusOK, testUser)
			}
		})
		testRoutesGroup.GET("/get", func(c *gin.Context) {
			user := &models.User{}

			err := user.FindByUsername("ABC")
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, "Err")
			} else {
				c.JSON(http.StatusOK, user)
			}
		})

		testRoutesGroup.GET("/testAuth", middlewares.Authenticate, func(c *gin.Context) {
			c.JSON(http.StatusOK, "Authed")
		})
	}
}
