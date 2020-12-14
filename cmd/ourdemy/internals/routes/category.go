package routes

import (
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(route *gin.Engine) {
	//categoryRoutesGroup := route.Group("/category")
	//{
	//	categoryRoutesGroup.GET("/", func(c *gin.Context) {
	//		cats, err := models.GetAllCategory()
	//		if err != nil {
	//			c.JSON(http.StatusInternalServerError, err)
	//			fmt.Println(err)
	//			return
	//		}
	//		c.JSON(http.StatusOK, cats)
	//	})
	//
	//	//TEST
	//	categoryRoutesGroup.GET("/create", func(c *gin.Context) {
	//		cat, _ := models.NewCategory("Cat_1")
	//		err := cat.Create()
	//		if err != nil {
	//			c.JSON(http.StatusInternalServerError, err)
	//			fmt.Println(err)
	//			return
	//		}
	//		c.JSON(http.StatusOK, "cat created")
	//	})
	//}

}
