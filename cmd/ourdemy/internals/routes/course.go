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
	"strconv"
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

		courseRoutesGroup.POST("/search", func(c *gin.Context) {
			catId, err := primitive.ObjectIDFromHex(c.Query("catId"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("category id invalid"),
				})
			}

			limitStr := c.Query("limit")
			offsetStr := c.Query("offset")
			if limitStr == "" {
				limitStr = "10"
			}

			if offsetStr == "" {
				offsetStr = "0"
			}

			limit, err := strconv.ParseInt(limitStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			}

			offset, err := strconv.ParseInt(offsetStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			}

			res, err := models.FindByCatId(catId, limit, offset)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
			}
			c.JSON(http.StatusOK, res)
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
					type courseCreateInfo struct {
						LecId        string  `json:"lid" bson:"lid" form:"lid" binding:"required"`
						CatId        string  `json:"cid" bson:"cat_id" form:"cid" binding:"required"`
						Name         string  `json:"name" bson:"name" form:"name" binding:"required"`
						ShortDesc    string  `json:"short_desc" bson:"short_desc" form:"short_desc" binding:"required"`
						FullDesc     string  `json:"full_desc" bson:"full_desc" form:"full_desc" binding:"required"`
						Fee          float64 `json:"fee" bson:"fee" form:"fee" binding:"required"`
						Discount     float64 `json:"discount" bson:"discount" form:"discount"`
						ChapterCount int     `json:"chapter_count" bson:"chapter_count" form:"chapter_count" binding:"required"`
					}

					var courseInfo courseCreateInfo
					if err := c.ShouldBind(&courseInfo); err != nil {
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

					lid, err := primitive.ObjectIDFromHex(courseInfo.LecId)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
						})
						return
					}

					cid, err := primitive.ObjectIDFromHex(courseInfo.CatId)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
						})
						return
					}

					course := models.Course{
						LecId:        lid,
						CatId:        cid,
						Ava:          base64.StdEncoding.EncodeToString(buff.Bytes()),
						Name:         courseInfo.Name,
						ShortDesc:    courseInfo.ShortDesc,
						FullDesc:     courseInfo.FullDesc,
						Fee:          courseInfo.Fee,
						Discount:     courseInfo.Discount,
						ChapterCount: courseInfo.ChapterCount,
					}

					if err := course.Save(); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}
					c.JSON(http.StatusOK, courseInfo)
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

			}

		}
	}
}
