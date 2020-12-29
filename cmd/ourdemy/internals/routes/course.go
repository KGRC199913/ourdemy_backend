package routes

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"github.com/oliamb/cutter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
)

func CourseRoutes(route *gin.Engine) {

	courseRoutesGroup := route.Group("/course")
	{
		courseRoutesGroup.GET("/simple/:cid", func(c *gin.Context) {
			cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("course id invalid"),
				})
				return
			}

			var course models.Course
			if err := course.FindById(cid); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			simple, err := course.ConvertToSimpleCourse()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, simple)
		})
		courseRoutesGroup.GET("/full/:cid", func(c *gin.Context) {
			cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("course id invalid"),
				})
				return
			}

			var course models.Course
			if err := course.FindById(cid); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			full, err := course.ConvertToFullCourse()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, full)
		})
		authCourseRoutesGroup := courseRoutesGroup.Group("/", middlewares.Authenticate)
		{
			authCourseRoutesGroup.POST("/buy/:cid", func(c *gin.Context) {
				cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course id invalid"),
					})
					return
				}

				uid, ok := c.Get("id")
				if !ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "id missing? wtf",
					})
					return
				}

				if err := models.AddUserToCourseInfo(uid.(primitive.ObjectID), cid); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"message": "successfully registered",
				})
			})
			lecturerCourseRoutesGroup := authCourseRoutesGroup.Group("/", middlewares.LecturerAuthenticate)
			{
				lecturerCourseRoutesGroup.POST("/create", func(c *gin.Context) {
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
				lecturerCourseRoutesGroup.POST("/markDone/:cid", func(c *gin.Context) {
					cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": errors.New("course id invalid"),
						})
						return
					}

					var course models.Course
					if err := course.FindById(cid); err != nil {
						c.JSON(http.StatusNotFound, gin.H{
							"error": err.Error(),
						})
						return
					}

					if course.IsDone {
						c.JSON(http.StatusConflict, gin.H{
							"error": "already mark done",
						})
						return
					}

					if err := course.UpdateCourseStatus(true); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					c.JSON(http.StatusOK, gin.H{
						"message": "marked as done",
					})
				})
				lecturerCourseRoutesGroup.POST("/markUndone/:cid", func(c *gin.Context) {
					cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": errors.New("course id invalid"),
						})
						return
					}

					var course models.Course
					if err := course.FindById(cid); err != nil {
						c.JSON(http.StatusNotFound, gin.H{
							"error": err.Error(),
						})
						return
					}

					if !course.IsDone {
						c.JSON(http.StatusConflict, gin.H{
							"error": "already is undone",
						})
						return
					}

					if err := course.UpdateCourseStatus(false); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					c.JSON(http.StatusOK, gin.H{
						"message": "marked as undone",
					})
				})
				lecturerCourseRoutesGroup.POST("/chapter", func(c *gin.Context) {
					var chapter models.CourseChapter
					if err := c.ShouldBind(&chapter); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
						})
						return
					}

					if err := chapter.Save(); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					c.JSON(http.StatusOK, chapter)
				})
				lecturerCourseRoutesGroup.POST("/chapter/:ccid", func(c *gin.Context) {
					var chapter models.CourseChapter
					if err := c.ShouldBind(&chapter); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
						})
						return
					}

					if err := chapter.Save(); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					c.JSON(http.StatusOK, chapter)
				})
				lecturerCourseRoutesGroup.DELETE("/chapter/:ccid", func(c *gin.Context) {
					ccid, err := primitive.ObjectIDFromHex(c.Param("ccid"))
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": errors.New("course chapter id invalid"),
						})
						return
					}

					var chapter models.CourseChapter
					if err := chapter.FindById(ccid); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					if err := chapter.Remove(); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					c.JSON(http.StatusOK, chapter)
				})

			}

		}
	}
}
