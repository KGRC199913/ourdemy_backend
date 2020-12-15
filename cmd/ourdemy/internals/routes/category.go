package routes

import (
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CategoryRoutes(route *gin.Engine) {
	categoryRoutesGroup := route.Group("/category")
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

		//TEST
		categoryRoutesGroup.GET("/create", func(c *gin.Context) {
			cat, _ := models.NewCategory("Cat_1")
			err := cat.Save()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, "cat created")
		})
		categoryRoutesGroup.GET("/sub/create", func(c *gin.Context) {
			cat, err := models.FindCategoryByName("Cat_1")
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				fmt.Println(err)
				return
			}
			subcat, _ := models.NewSubCategory("SubCat_1", cat.Id.Hex())
			err = subcat.Save()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, "cat created")
		})
	}
}
