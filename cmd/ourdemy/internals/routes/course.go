package routes

import (
	"bytes"
	"encoding/base64"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"github.com/oliamb/cutter"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
)

func CourseRoutes(route *gin.Engine) {
	courseRoutesGroup := route.Group("/course", middlewares.Authenticate, middlewares.LecturerAuthenticate)
	{
		courseRoutesGroup.POST("/create", func(c *gin.Context) {
			var course models.Course
			if err := c.ShouldBind(&course); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			avaff, err := c.FormFile("ava")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			avaf, err := avaff.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			imgData, imgType, err := image.Decode(avaf)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			cropped, err := cutter.Crop(imgData, cutter.Config{
				Width:  300,
				Height: 300,
				Mode:   cutter.Centered,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			var buff bytes.Buffer
			switch imgType {
			case "png":
				if err := png.Encode(&buff, cropped); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
			case "jpeg":
				if err := jpeg.Encode(&buff, imgData, nil); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
			default:
				{
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "unknown img type",
					})
					return
				}
			}

			course.Ava = base64.StdEncoding.EncodeToString(buff.Bytes())
			if err := course.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, course)
		})
	}
}
