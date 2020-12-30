package routes

import (
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		categoryRoutesGroup.GET("/all", func(c *gin.Context) {
			type subCatItem struct {
				SubCatId   primitive.ObjectID `json:"scid"`
				SubCatName string             `json:"subcat_name"`
			}

			type catItem struct {
				CatId   primitive.ObjectID `json:"cid"`
				CatName string             `json:"cat_name"`
				Subcat  []subCatItem       `json:"subcats"`
			}

			cats, err := models.GetAllCategory()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				fmt.Println(err)
				return
			}

			var res []catItem
			for _, cat := range cats {
				subcats, err := models.FindSubcatsByCatName(cat.Name)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
					fmt.Println(err)
					return
				}

				catItem := catItem{
					CatId:   cat.Id,
					CatName: cat.Name,
					Subcat:  []subCatItem{},
				}

				for _, subcat := range subcats {
					catItem.Subcat = append(catItem.Subcat, subCatItem{
						SubCatId:   subcat.Id,
						SubCatName: subcat.Name,
					})
				}

				res = append(res, catItem)
			}

			c.JSON(http.StatusOK, res)
		})

		categoryRoutesGroup.POST("/create", middlewares.AdminAuthenticate, func(c *gin.Context) {
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
			subcatRoutesGroup.GET("/:catName", func(c *gin.Context) {
				catName := c.Param("catName")
				if catName == "" {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "wrong params",
					})
					return
				}

				subcats, err := models.FindSubcatsByCatName(catName)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"errors": "not found",
					})
					return
				}

				c.JSON(http.StatusOK, subcats)
			})
			subcatRoutesGroup.POST("/create", middlewares.AdminAuthenticate, func(c *gin.Context) {
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
					return
				}

				c.JSON(http.StatusOK, subCat)
			})
		}
	}
}
