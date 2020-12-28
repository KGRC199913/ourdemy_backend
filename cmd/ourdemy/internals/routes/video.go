package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sony/sonyflake"
	"io"
	"net/http"
	"os"
	"strconv"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{})

func videoRoutes(route *gin.Engine) {
	videoRoutesGroup := route.Group("/vid")
	{
		videoRoutesGroup.POST("/test/upload", func(c *gin.Context) {
			file, _ := c.FormFile("file")
			f, _ := file.Open()
			fileUid, err := flake.NextID()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "video upload failed",
				})
			}
			target, _ := os.Create("/vid/" + strconv.FormatUint(fileUid, 10))
			_, _ = io.Copy(target, f)
		})
	}
}
