package routes

import (
	"errors"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func ReviewRoutes(route *gin.Engine) {
	reviewRoutesGroup := route.Group("/reviews")
	{
		reviewRoutesGroup.GET("/:cid", func(c *gin.Context) {
			cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("course id invalid"),
				})
				return
			}

			revs, err := models.FindByCourseId(cid)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			var dpRvs []models.DisplayableReview
			for _, rev := range revs {
				dpRv, err := rev.ConvertToDisplayableReview()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "something went wrong",
					})

					return
				}

				dpRvs = append(dpRvs, *dpRv)
			}

			if dpRvs == nil {
				dpRvs = []models.DisplayableReview{}
			}

			c.JSON(http.StatusOK, dpRvs)
		})

		reviewRoutesGroup.POST("/create", middlewares.Authenticate, func(c *gin.Context) {
			type createReview struct {
				CourseId string  `json:"cid" binding:"required"`
				Content  string  `json:"content" binding:"required"`
				Score    float32 `json:"score" binding:"required"`
			}

			var newReview createReview

			if err := c.ShouldBind(&newReview); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			if newReview.Score < 0 || newReview.Score > 5 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("invalid score"),
				})
				return
			}

			uid, _ := c.Get("id")
			cid, err := primitive.ObjectIDFromHex(newReview.CourseId)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			curReview := models.Review{
				UserId:   uid.(primitive.ObjectID),
				CourseId: cid,
				Content:  newReview.Content,
				Score:    newReview.Score,
			}

			if err := curReview.Save(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, curReview)
		})

		reviewRoutesGroup.POST("/update/:rid", middlewares.Authenticate, func(c *gin.Context) {
			type updateReview struct {
				Content string  `json:"content" binding:"required"`
				Score   float32 `json:"score" binding:"required"`
			}

			var curUpdateReview updateReview

			if err := c.ShouldBind(&curUpdateReview); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			curReview := models.Review{}

			rid, err := primitive.ObjectIDFromHex(c.Param("rid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("review id invalid"),
				})
				return
			}

			if err := curReview.FindById(rid); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := curReview.UpdateReview(curUpdateReview.Content, curUpdateReview.Score); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Update successful.",
			})
		})

		reviewRoutesGroup.DELETE("/delete/:rid", middlewares.Authenticate, func(c *gin.Context) {
			curReview := models.Review{}
			rid, err := primitive.ObjectIDFromHex(c.Param("rid"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("review id invalid"),
				})
				return
			}

			if err := curReview.FindById(rid); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			curUserId, _ := c.Get("id")

			if curUserId.(primitive.ObjectID) != curReview.UserId {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": errors.New("not your review"),
				})
				return
			}

			if err := curReview.DeleteReview(rid); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Delete successful.",
			})
		})
	}
}
