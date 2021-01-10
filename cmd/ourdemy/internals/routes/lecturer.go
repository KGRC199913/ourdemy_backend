package routes

import (
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func LecturerRoutes(route *gin.Engine) {
	lecRoutesGroup := route.Group("/lecturers", middlewares.Authenticate)
	{
		lecRoutesGroup.POST("/promote", func(c *gin.Context) {
			id, ok := c.Get("id")
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Id not found",
				})
				return
			}

			appr := models.Approve{
				LecId: id.(primitive.ObjectID),
			}
			if err := appr.Save(); err != nil {
				c.JSON(http.StatusNotAcceptable, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "register success",
			})
		})
	}

}
