package routes

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
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
		courseRoutesGroup.GET("/all", func(c *gin.Context) {
			res, err := models.GetAllCourseAsSimple()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})

				return
			}

			c.JSON(http.StatusOK, res)
		})
		courseRoutesGroup.GET("/all/cat/:cat_id", func(c *gin.Context) {
			catId, err := primitive.ObjectIDFromHex(c.Param("cat_id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "cat id invalid",
				})
				return
			}

			res, err := models.GetAllCourseByCatIdAsSimple(catId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})

				return
			}

			c.JSON(http.StatusOK, res)
		})
		courseRoutesGroup.GET("/all/subcat/:subcat_id", func(c *gin.Context) {
			subcatId, err := primitive.ObjectIDFromHex(c.Param("subcat_id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "cat id invalid",
				})
				return
			}

			res, err := models.GetAllCourseBySubcatIdAsSimple(subcatId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})

				return
			}

			c.JSON(http.StatusOK, res)
		})

		courseRoutesGroup.GET("/chapter/:cid", func(c *gin.Context) {
			cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			chapters, err := models.FindAllChapterByCourseId(cid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, chapters)
		})
		courseRoutesGroup.GET("/simple/:cid", func(c *gin.Context) {
			cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "course id invalid",
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
					"error": "course id invalid",
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
			err = course.UpdateWatchCount()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, full)
		})

		courseRoutesGroup.GET("/search", func(c *gin.Context) {
			keyword := c.DefaultQuery("keyword", "")

			limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "limit invalid",
				})

				return
			}
			offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "offset invalid",
				})

				return
			}

			res, err := models.SearchByKeyword(keyword, int64(limit), int64(offset))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "something when wrong",
				})

				return
			}

			catId := c.DefaultQuery("catId", "")
			if catId == "" {
				subcatId := c.DefaultQuery("subcatId", "")
				if subcatId == "" {
					var r []models.SimpleCourse
					for _, item := range res {
						data, err := item.ConvertToSimpleCourse()
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": "something went wrong",
							})

							return
						}

						r = append(r, *data)
					}

					if r == nil {
						r = []models.SimpleCourse{}
					}

					c.JSON(http.StatusOK, r)
					return
				}

				var filteredRes []models.Course
				for _, item := range res {
					if item.CatId.Hex() == subcatId {
						filteredRes = append(filteredRes, item)
					}
				}

				if filteredRes == nil {
					filteredRes = []models.Course{}
				}

				var r []models.SimpleCourse
				for _, item := range filteredRes {
					data, err := item.ConvertToSimpleCourse()
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": "something went wrong",
						})

						return
					}

					r = append(r, *data)
				}

				if r == nil {
					r = []models.SimpleCourse{}
				}

				c.JSON(http.StatusOK, r)
				return
			}

			var filteredRes []models.Course
			catOid, err := primitive.ObjectIDFromHex(catId)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "not found",
				})
				return
			}

			subcats, err := models.FindByParentCategoryId(catOid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "something went wrong",
				})
				return
			}

			for _, item := range res {
				if ContainsSubcatId(subcats, item.CatId) {
					filteredRes = append(filteredRes, item)
				}
			}

			if filteredRes == nil {
				filteredRes = []models.Course{}
			}

			var r []models.SimpleCourse
			for _, item := range filteredRes {
				data, err := item.ConvertToSimpleCourse()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "something went wrong",
					})

					return
				}

				r = append(r, *data)
			}

			if r == nil {
				r = []models.SimpleCourse{}
			}

			c.JSON(http.StatusOK, r)
		})

		courseRoutesGroup.GET("/searchByCatId", func(c *gin.Context) {
			catId, err := primitive.ObjectIDFromHex(c.Query("catId"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("category id invalid"),
				})
			}

			limitStr := c.DefaultQuery("limit", "5")
			offsetStr := c.DefaultQuery("offset", "0")

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

		courseRoutesGroup.GET("/search/subcat", func(c *gin.Context) {
			subcatId, err := primitive.ObjectIDFromHex(c.Query("subcatId"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("category id invalid"),
				})
			}

			limitStr := c.DefaultQuery("limit", "5")
			offsetStr := c.DefaultQuery("offset", "0")

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

			res, err := models.FindBySubcatId(subcatId, limit, offset)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
			}
			c.JSON(http.StatusOK, res)
		})

		courseRoutesGroup.GET("/relevance/:cid", func(c *gin.Context) {
			cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "course id invalid",
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

			res, err := models.Get5RandomCourseBySubcat(course.CatId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "something wrong",
				})

				return
			}

			var scRes []models.SimpleCourse
			for _, course := range res {
				sc, err := course.ConvertToSimpleCourse()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "something went wrong",
					})

					return
				}
				scRes = append(scRes, *sc)
			}

			if scRes == nil {
				scRes = []models.SimpleCourse{}
			}

			c.JSON(http.StatusOK, scRes)
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
			authCourseRoutesGroup.GET("/checkJoined/:cid", func(c *gin.Context) {
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

				c.JSON(http.StatusOK, models.IsUserJoined(cid, uid.(primitive.ObjectID)))
			})
			lecturerCourseRoutesGroup := authCourseRoutesGroup.Group("/", middlewares.LecturerAuthenticate)
			{
				lecturerCourseRoutesGroup.GET("/allByMe", func(c *gin.Context) {
					uid, ok := c.Get("id")
					if !ok {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": "id missing? wtf",
						})
						return
					}

					var courses []models.Course
					courses, err := models.FindByLecId(uid.(primitive.ObjectID))
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					var full []models.FullCourse
					for _, course := range courses {
						fullCourse, err := course.ConvertToFullCourse()
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err.Error(),
							})
							return
						}
						full = append(full, *fullCourse)
					}

					if full == nil {
						full = []models.FullCourse{}
					}

					c.JSON(http.StatusOK, full)
				})
				lecturerCourseRoutesGroup.POST("/create", func(c *gin.Context) {
					type courseCreateInfo struct {
						CatId     string  `json:"cid" bson:"cat_id" form:"cid" binding:"required"`
						Name      string  `json:"name" bson:"name" form:"name" binding:"required"`
						ShortDesc string  `json:"short_desc" bson:"short_desc" form:"short_desc"`
						FullDesc  string  `json:"full_desc" bson:"full_desc" form:"full_desc"`
						Fee       float64 `json:"fee" bson:"fee" form:"fee" binding:"required"`
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

					uid, ok := c.Get("id")
					if !ok {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": "id missing? wtf",
						})
						return
					}
					lid := uid.(primitive.ObjectID)

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
						Discount:     0.0,
						ChapterCount: 0,
					}

					if err := course.Save(); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}
					c.JSON(http.StatusOK, courseInfo)
				})
				lecturerCourseRoutesGroup.POST("/update/:cid", func(c *gin.Context) {
					type updateCourseDescData struct {
						Short    string  `json:"short_desc" form:"short_desc" binding:"required"`
						Full     string  `json:"full_desc" form:"full_desc" binding:"required"`
						Discount float64 `json:"discount" form:"discount"`
					}

					cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": "course id invalid",
						})
						return
					}

					var updateData updateCourseDescData
					if err := c.ShouldBind(&updateData); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
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

					uid, ok := c.Get("id")
					if !ok {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": "id missing? wtf",
						})
						return
					}
					lid := uid.(primitive.ObjectID)

					if course.LecId != lid {
						c.JSON(http.StatusUnauthorized, gin.H{
							"error": "not owned this course",
						})

						return
					}

					if err := course.UpdateCourseDesc(updateData.Short, updateData.Full); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					if err := course.UpdateDiscount(updateData.Discount); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

					full, err := course.ConvertToFullCourse()
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": "something when wrong",
						})
						fmt.Println(err.Error())
						return
					}

					c.JSON(http.StatusOK, full)
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

					chapter.Videos = []models.VideoMetadata{}

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

func ContainsSubcatId(s []models.SubCategory, e primitive.ObjectID) bool {
	for _, a := range s {
		if a.Id == e {
			return true
		}
	}
	return false
}
