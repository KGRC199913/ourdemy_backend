package routes

import (
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func AdminRoutes(route *gin.Engine) {
	adminroutesGroup := route.Group("/admin", middlewares.AdminAuthenticate)
	{
		adminroutesGroup.GET("/promote", func(c *gin.Context) {
			apprs, err := models.GetAllApprovingLecturers()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, apprs)
		})
		adminroutesGroup.POST("/promote/:id", func(c *gin.Context) {
			var appr models.Approve
			id := c.Param("id")
			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			if err := appr.FindByLecId(oid); err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			var user models.User

			if err := user.FindById(appr.LecId); err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			if err := user.UpdateLecturerStatus(true); err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			if err := appr.Remove(); err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
		})
	}
}
