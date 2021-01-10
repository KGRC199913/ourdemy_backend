package routes

import (
	"errors"
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/middlewares"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func UserRoutes(route *gin.Engine) {
	adminRoutesGroup := route.Group("/admin")
	{
		adminRoutesGroup.POST("/signin", func(c *gin.Context) {
			type signinAdmin struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			var curSigninAdmin signinAdmin
			if err := c.ShouldBind(&curSigninAdmin); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			if curSigninAdmin.Username != viper.GetString("ADMIN_USERNAME") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "invalid login",
				})
				return
			}

			if curSigninAdmin.Password != viper.GetString("ADMIN_PASSWORD") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "invalid login",
				})
				return
			}

			accessToken, err := ultis.CreateAdminToken()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"accessToken": accessToken,
			})
		})
	}

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

			auth := LoginAuth(viper.GetString("USERNAME"), viper.GetString("PASSWORD"))
			to := []string{user.Email}
			msg := []byte("To: " + user.Email + "\r\n" +
				"Subject: Ourdemy Announcement\r\n" +
				"\r\n" + "OTP: " +
				user.CurOtp + "\nExpired Time: " + user.CurOtpExpiredTime.Format("2006-01-02 15:04:05") + "\r\n")
			err := smtp.SendMail("smtp.gmail.com:587", auth, viper.GetString("USERNAME"), to, msg)
			if err != nil {
				panic(err.Error())
			}

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
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
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
					"error": "Login failed",
				})
				return
			}

			if curUser.IsBanned {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Your account has been banned, please contact admin",
				})
				return
			}

			accessToken, err := ultis.CreateToken(curUser.Id, curUser.IsLec)
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

			c.JSON(http.StatusOK, curUpdateUser)
		})

		userRoutesGroup.POST("/updatePassword", middlewares.Authenticate, func(c *gin.Context) {
			type userPasswordUpdate struct {
				OldPassword string `json:"old_password" binding:"required"`
				NewPassword string `json:"new_password" binding:"required"`
			}
			var curUpdateUserPassword userPasswordUpdate
			if err := c.ShouldBind(&curUpdateUserPassword); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			curUser := models.User{}

			curUserId, _ := c.Get("id")
			if err := curUser.FindById(curUserId.(primitive.ObjectID)); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			err := scrypt.CompareHashAndPassword([]byte(curUser.HPassword), []byte(curUpdateUserPassword.OldPassword))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Wrong old password",
				})
				return
			}

			hashed, err := scrypt.GenerateFromPassword([]byte(curUpdateUserPassword.NewPassword), scrypt.DefaultParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "something went wrong",
				})
				return
			}

			if err := curUser.UpdatePassword(string(hashed)); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, curUpdateUserPassword)
		})

		userRoutesGroup.GET("/profile", middlewares.Authenticate, func(c *gin.Context) {
			var curUser models.User
			curUserId, _ := c.Get("id")
			if err := curUser.FindById(curUserId.(primitive.ObjectID)); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, curUser)
		})

		userRoutesGroup.POST("/fav/:cid", middlewares.Authenticate, func(c *gin.Context) {
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

			uid, ok := c.Get("id")
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "id missing? wtf",
				})
				return
			}

			var wl models.WatchList
			if err := wl.FindByUid(uid.(primitive.ObjectID)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := wl.AddCourseToWatchList(cid); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "added to fav list",
			})
		})

		userRoutesGroup.POST("/unfav/:cid", middlewares.Authenticate, func(c *gin.Context) {
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

			uid, ok := c.Get("id")
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "id missing? wtf",
				})
				return
			}

			var wl models.WatchList
			if err := wl.FindByUid(uid.(primitive.ObjectID)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			if err := wl.RemoveCourseFromWatchList(cid); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "added to fav list",
			})
		})

		userRoutesGroup.GET("/favList", middlewares.Authenticate, func(c *gin.Context) {
			uid, ok := c.Get("id")
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "id missing? wtf",
				})
				return
			}

			var wl models.WatchList
			if err := wl.FindByUid(uid.(primitive.ObjectID)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			type extremeSimpleCourse struct {
				Id   primitive.ObjectID `json:"cid"`
				Name string             `json:"name"`
			}

			var res []extremeSimpleCourse
			var course models.Course
			for _, courseId := range wl.CoursesId {
				if err := course.FindById(courseId); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "something went wrong",
					})
					fmt.Println(err.Error())
					return
				}
				res = append(res, extremeSimpleCourse{
					Id:   courseId,
					Name: course.Name,
				})
			}

			c.JSON(http.StatusOK, res)
		})
	}
}
