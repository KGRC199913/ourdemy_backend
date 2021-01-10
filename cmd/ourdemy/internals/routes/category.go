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

			if res == nil {
				res = []catItem{}
			}

			c.JSON(http.StatusOK, res)
		})

		categoryRoutesGroup.GET("/sub/:catName", func(c *gin.Context) {
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

		adminCatsGroup := categoryRoutesGroup.Group("/admin", middlewares.AdminAuthenticate)
		{
			adminCatsGroup.POST("/create", func(c *gin.Context) {
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
			adminCatsGroup.POST("/update", func(c *gin.Context) {
				type updateCatInfo struct {
					Name string             `json:"name" binding:"required"`
					Id   primitive.ObjectID `json:"id" binding:"required"`
				}

				var info updateCatInfo
				if err := c.ShouldBind(&info); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}

				var cat models.Category
				if err := cat.FindCategoryById(info.Id); err != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"error": "cat not found",
					})

					return
				}

				if err := cat.UpdateName(info.Name); err != nil {
					c.JSON(http.StatusInternalServerError, err)
					fmt.Println(err)
					return
				}

				c.JSON(http.StatusOK, cat)
			})

			adminCatsGroup.DELETE("/delete/:cid", func(c *gin.Context) {
				cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "cat id invalid",
					})
					return
				}

				var cat models.Category
				if err := cat.FindCategoryById(cid); err != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"error": err.Error(),
					})
					return
				}

				courses, _ := models.FindAllCourseByCatId(cid)
				if courses != nil {
					c.JSON(http.StatusForbidden, gin.H{
						"error": "category not empty, please remove all course related to this cat first",
					})

					return
				}

				err = cat.Remove()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"Id": cat.Id,
				})
			})
			subcatRoutesGroup := adminCatsGroup.Group("/sub")
			{
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

					c.JSON(http.StatusOK, gin.H{
						"scid":    subCat.Id,
						"name":    subCat.Name,
						"catName": subcatCreateInfo.ParentName,
					})
				})
				subcatRoutesGroup.POST("/update", func(c *gin.Context) {
					type updateSubcatInfo struct {
						Name string             `json:"name" binding:"required"`
						Id   primitive.ObjectID `json:"id" binding:"required"`
					}

					var info updateSubcatInfo
					if err := c.ShouldBind(&info); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
						})
						return
					}

					var cat models.SubCategory
					if err := cat.FindSubCategoryById(info.Id); err != nil {
						c.JSON(http.StatusNotFound, gin.H{
							"error": "cat not found",
						})

						return
					}

					if err := cat.UpdateName(info.Name); err != nil {
						c.JSON(http.StatusInternalServerError, err)
						fmt.Println(err)
						return
					}

					c.JSON(http.StatusOK, cat)
				})
				subcatRoutesGroup.DELETE("/delete/:scid", middlewares.AdminAuthenticate, func(c *gin.Context) {
					scid, err := primitive.ObjectIDFromHex(c.Param("scid"))
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": "cat id invalid",
						})
						return
					}

					var sbcat models.SubCategory
					if err := sbcat.FindSubCategoryById(scid); err != nil {
						c.JSON(http.StatusNotFound, gin.H{
							"error": err.Error(),
						})
						return
					}

					courses, _ := models.FindAllCourseBySubcatId(scid)
					if courses != nil {
						c.JSON(http.StatusForbidden, gin.H{
							"error": "subcategory not empty, please remove all course related to this subcat first",
						})

						return
					}

					err = sbcat.Remove()
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					c.JSON(http.StatusOK, gin.H{
						"Id":  sbcat.Id,
						"Cid": sbcat.ParentCategoryId,
					})
				})
			}
		}

	}
}
