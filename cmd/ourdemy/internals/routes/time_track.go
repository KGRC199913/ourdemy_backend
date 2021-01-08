package routes

import (
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func TimeMark(route *gin.Engine) {
	timeMarkRoutesGroup := route.Group("/time")
	{
		timeMarkRoutesGroup.GET("/:vid", middlewares.Authenticate, func(c *gin.Context) {
			vid, err := primitive.ObjectIDFromHex(c.Param("vid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			var timeMark models.TimeMark
			if err := timeMark.FindByVideoId(vid); err != nil {
				c.JSON(http.StatusOK, 0)
			}

			c.JSON(http.StatusOK, timeMark.CurTime)
		})

		timeMarkRoutesGroup.POST("/:vid", middlewares.Authenticate, func(c *gin.Context) {
			type upsertTimeMark struct {
				CurTime int `json:"cur_time"`
			}
			var curTimeMark upsertTimeMark
			if err := c.ShouldBind(&curTimeMark); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			vid, err := primitive.ObjectIDFromHex(c.Param("vid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			curUserId, _ := c.Get("id")

			var timeMark models.TimeMark
			if err := timeMark.FindByVideoId(vid); err != nil {
				timeMark.UserId = curUserId.(primitive.ObjectID)
				timeMark.VideoId = vid
				timeMark.CurTime = int64(curTimeMark.CurTime)

				if err := timeMark.Save(); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}
			}

			if err := timeMark.UpdateTimeMark(int64(curTimeMark.CurTime)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, timeMark)
		})
	}
}
