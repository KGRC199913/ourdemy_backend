package routes

import (
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CategoryRoutes(route *gin.Engine) {
	categoryRoutesGroup := route.Group("/category", middlewares.AdminAuthenticate)
	{
		categoryRoutesGroup.GET("/", func(c *gin.Context) {
			cats, err := models.GetAllCategory()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, cats)
		})

		categoryRoutesGroup.POST("/create", func(c *gin.Context) {
			var cat models.Category
			if err := c.ShouldBind(&cat); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := cat.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, err)
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, cat)
		})

		subcatRoutesGroup := categoryRoutesGroup.Group("/sub")
		{
			subcatRoutesGroup.POST("/create", func(c *gin.Context) {
				type subcatCreate struct {
					Name       string `json:"name" binding:"required"`
					ParentName string `json:"parent_name" binding:"required"`
				}

				var subcatCreateInfo subcatCreate
				if err := c.ShouldBind(&subcatCreateInfo); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}

				subCat, err := models.CreateSubCategory(subcatCreateInfo.Name, subcatCreateInfo.ParentName)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}

				if err := subCat.Save(); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
				}

				c.JSON(http.StatusOK, subCat)
			})
		}
	}
}
