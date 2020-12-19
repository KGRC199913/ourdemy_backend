package routes

import (
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func UserRoutes(route *gin.Engine) {
	userRoutesGroup := route.Group("/users")
	{
		userRoutesGroup.POST("/signup", func(c *gin.Context) {
			var user models.User
			if err := c.ShouldBind(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := user.GenerateOtp(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := user.Save(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			//TODO SEND OTP
			//END SEND OTP

			c.JSON(http.StatusOK, user)
		})

		userRoutesGroup.PATCH("/otp", func(c *gin.Context) {
			type otpValidate struct {
				Username string `json:"username"`
				Otp      string `json:"otp"`
			}

			var otpVal otpValidate
			if err := c.ShouldBind(&otpVal); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			curValidationUser := models.User{}
			if err := curValidationUser.ConfirmOtp(otpVal.Username, otpVal.Otp); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := curValidationUser.GenerateRfToken(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "otp confirmed",
			})
		})

		userRoutesGroup.POST("/resendOtp", func(c *gin.Context) {
			type resendOtpVal struct {
				Username string `json:"username"`
			}
			var resendOtp resendOtpVal
			if err := c.ShouldBind(&resendOtp); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			curUser := models.User{}
			var newOtp *string
			var err error
			if newOtp, err = curUser.UpdateOtp(resendOtp.Username); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"new_otp": newOtp,
			})
		})

		userRoutesGroup.POST("/signin", func(c *gin.Context) {
			type signinUser struct {
				Username string `json:"username"`
				Password string `json:"pass"`
			}
			var curSigninUser signinUser
			if err := c.ShouldBind(&curSigninUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			var curUser models.User
			if err := curUser.FindByUsername(curSigninUser.Username); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Login failed",
				})
				return
			}

			if curUser.CurOtp != "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "OTP not confirmed",
				})
				return
			}

			err := scrypt.CompareHashAndPassword([]byte(curUser.HPassword), []byte(curSigninUser.Password))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Wrong Password",
				})
				return
			}

			accessToken, err := ultis.CreateToken(curUser.Id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"accessToken":  accessToken,
				"refreshToken": curUser.RefreshToken,
			})
		})

		userRoutesGroup.POST("/update", middlewares.Authenticate, func(c *gin.Context) {
			type userUpdate struct {
				Fullname string `json:"fullname" binding:"required"`
				Email    string `json:"email" binding:"required"`
			}
			var curUpdateUser userUpdate
			if err := c.ShouldBind(&curUpdateUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			fmt.Println(curUpdateUser.Fullname)
			fmt.Println(curUpdateUser.Email)

			curUser := models.User{}

			curUserId, _ := c.Get("id")
			if err := curUser.FindById(curUserId.(primitive.ObjectID)); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := curUser.UpdateProfile(curUpdateUser.Fullname, curUpdateUser.Email); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Update successful.",
			})
		})
	}
}
