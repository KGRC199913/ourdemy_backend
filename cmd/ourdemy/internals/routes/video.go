package routes

import (
	"errors"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func VideoRoutes(route *gin.Engine) {
	videoRoutesGroup := route.Group("/vid")
	{
		videoRoutesGroup.GET("/chapter/:ccid", func(c *gin.Context) {
			ccid, err := primitive.ObjectIDFromHex(c.Param("ccid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("course chapter id invalid"),
				})
				return
			}

			var courseChapter models.CourseChapter
			if err := courseChapter.FindById(ccid); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("chapter not exist"),
				})
				return
			}

			vms, err := models.FindAllVideoMetadataByChapterId(ccid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, vms)
		})
		authenVidRoutesGroup := videoRoutesGroup.Group("/", middlewares.UrlAuthenticate)
		{
			authenVidRoutesGroup.GET("/download/:vid", func(c *gin.Context) {
				vid, err := primitive.ObjectIDFromHex(c.Param("vid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("vid id invalid"),
					})
					return
				}

				var vm models.VideoMetadata
				if err := vm.FindById(vid); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("video not exist"),
					})
					return
				}

				var chap models.CourseChapter
				if err := chap.FindById(vm.ChapterId); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "something went wrong",
					})
					return
				}

				if *chap.Previewable {
					c.File("vid/" + vm.Id.Hex())
					return
				}

				uid, exist := c.Get("id")
				if !exist {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": errors.New("something went wrong"),
					})
					return
				}

				ouid := uid.(primitive.ObjectID)
				if !models.IsUserJoined(vm.CourseId, ouid) {
					c.JSON(http.StatusForbidden, gin.H{
						"error": errors.New("didn't join this course"),
					})
					return
				}

				c.File("vid/" + vm.Id.Hex())
			})
		}
		lecVidRoutesGroup := videoRoutesGroup.Group("/", middlewares.Authenticate, middlewares.LecturerAuthenticate)
		{
			lecVidRoutesGroup.PUT("/upload/:cid/:ccid", func(c *gin.Context) {
				type uploadVidData struct {
					Title string `json:"title" form:"title" binding:"required"`
				}

				cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "course id invalid",
					})
					return
				}

				ccid, err := primitive.ObjectIDFromHex(c.Param("ccid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "course chapter id invalid",
					})
					return
				}

				var uploadData uploadVidData
				if err := c.ShouldBind(&uploadData); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}

				var course models.Course
				if err := course.FindById(cid); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "course not exist",
					})
					return
				}

				lecid, exist := c.Get("id")
				if !exist {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "something went wrong",
					})
					return
				}

				if course.LecId != lecid.(primitive.ObjectID) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "course does not belong to you",
					})
					return
				}

				var courseChapter models.CourseChapter
				if err := courseChapter.FindById(ccid); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "chapter not exist",
					})
					return
				}

				if courseChapter.CourseId != course.Id {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "invalid chapter id",
					})
					return
				}

				ff, err := c.FormFile("vid")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "video missing",
					})
					return
				}

				f, err := ff.Open()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "video upload failed",
					})
					return
				}

				vm := models.VideoMetadata{
					ChapterId: courseChapter.Id,
					CourseId:  course.Id,
					Title:     uploadData.Title,
				}
				if err := vm.Save(); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "video upload failed",
					})
					return
				}

				target, err := create("vid/" + vm.Id.Hex())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "video upload failed",
					})
					return
				}
				defer target.Close()

				_, err = io.Copy(target, f)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "video upload failed",
					})
					go vm.Remove()
					return
				}
				c.JSON(http.StatusOK, vm)
			})
			lecVidRoutesGroup.POST("/:cid/:ccid", func(c *gin.Context) {
				type updateVidData struct {
					Title   string             `json:"title" form:"title"`
					VideoId primitive.ObjectID `json:"video_id" form:"video_id" binding:"required"`
				}
				cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course id invalid"),
					})
					return
				}

				ccid, err := primitive.ObjectIDFromHex(c.Param("ccid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course chapter id invalid"),
					})
					return
				}

				var course models.Course
				if err := course.FindById(cid); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course not exist"),
					})
					return
				}

				lecid, exist := c.Get("id")
				if !exist {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": errors.New("something went wrong"),
					})
					return
				}

				if course.LecId != lecid.(primitive.ObjectID) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course does not belong to you"),
					})
					return
				}

				var courseChapter models.CourseChapter
				if err := courseChapter.FindById(ccid); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("chapter not exist"),
					})
					return
				}

				if courseChapter.CourseId != course.Id {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("invalid chapter id"),
					})
					return
				}

				var updateData updateVidData
				if err := c.ShouldBind(&updateData); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}

				var vm models.VideoMetadata
				if err := vm.FindById(updateData.VideoId); err != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"error": errors.New("video not found"),
					})
					return
				}

				if updateData.Title != "" {
					if err := vm.UpdateVideoTitle(updateData.Title); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}
				}

				ff, err := c.FormFile("vid")
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"message": "video title updated",
					})
					return
				}

				if err := os.Remove("vid/" + vm.Id.Hex()); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				f, err := ff.Open()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "video upload failed",
					})
					return
				}

				target, err := create("vid/" + vm.Id.Hex())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "video upload failed",
					})
					return
				}
				_, err = io.Copy(target, f)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "video upload failed",
					})
					go vm.Remove()
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"message": "video updated",
				})
			})
			lecVidRoutesGroup.DELETE("/:cid/:ccid/:vid", func(c *gin.Context) {
				cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course id invalid"),
					})
					return
				}

				ccid, err := primitive.ObjectIDFromHex(c.Param("ccid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course chapter id invalid"),
					})
					return
				}

				vid, err := primitive.ObjectIDFromHex(c.Param("vid"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course chapter id invalid"),
					})
					return
				}

				var course models.Course
				if err := course.FindById(cid); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course not exist"),
					})
					return
				}

				lecid, exist := c.Get("id")
				if !exist {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": errors.New("something went wrong"),
					})
					return
				}

				if course.LecId != lecid.(primitive.ObjectID) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("course does not belong to you"),
					})
					return
				}

				var courseChapter models.CourseChapter
				if err := courseChapter.FindById(ccid); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("chapter not exist"),
					})
					return
				}

				if courseChapter.CourseId != course.Id {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("invalid chapter id"),
					})
					return
				}

				var vm models.VideoMetadata
				if err := vm.FindById(vid); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("video not exist"),
					})
					return
				}

				if vm.ChapterId != courseChapter.Id || vm.CourseId != course.Id {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errors.New("video not belong to course/chapter"),
					})
					return
				}

				if err := os.Remove("vid/" + vm.Id.Hex()); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				if err := vm.Remove(); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				c.JSON(http.StatusOK, vm)
			})
		}
	}
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
