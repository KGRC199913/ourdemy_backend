package routes

import (
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

			//TODO GEN OTP

			//FAKE OTP
			user.CurOtp = "1234"
			user.CurOtpExpiredTime = time.Now().Add(time.Minute * 30)
			//END FAKE OTP
			//END GEN OTP

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

			c.JSON(http.StatusOK, gin.H{
				"message": "otp confirmed",
			})
		})
	}
}
