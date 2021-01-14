package middlewares

import (
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Authenticate(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	auths := strings.Split(authToken, "Bearer ")
	if len(auths) < 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		c.Abort()
		return
	}

	authToken = auths[1]
	var userClaims ultis.UserClaims
	token, err := ultis.ParseToken(authToken, &userClaims)

	if err != nil {
		validationError, _ := err.(*jwt.ValidationError)

		if validationError.Errors == jwt.ValidationErrorExpired {
			rfToken := c.GetHeader("Refresh")
			if rfToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "please log in again",
				})
				c.Abort()
				return
			}

			var rfUser models.User
			if err := rfUser.FindByIdAndRfToken(userClaims.Id, rfToken); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				c.Abort()
				return
			}

			if rfUser.IsBanned {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "banned",
				})

				c.Abort()
				return
			}

			newAccessToken, newRfToken, err := rfUser.UpdateTokens()
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				c.Abort()
				return
			}

			c.Writer.Header().Set("AccessToken", *newAccessToken)
			c.Writer.Header().Set("RefreshToken", *newRfToken)
		}

		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("id", userClaims.Id)
		c.Next()
		return
	}

	if !token.Valid {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
	}

	if models.IsBanned(userClaims.Id) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "banned",
		})

		c.Abort()
		return
	}

	c.Set("id", userClaims.Id)
	c.Set("is_lec", userClaims.IsLec)
	c.Next()
}

func LecturerAuthenticate(c *gin.Context) {
	if !c.GetBool("is_lec") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "lecturer unauthorized",
		})
		c.Abort()
	}
	c.Next()
}

func UrlAuthenticate(c *gin.Context) {
	authToken := c.Query("auth")

	if authToken == "none" {
		c.Next()
		return
	}

	var userClaims ultis.UserClaims
	token, err := ultis.ParseToken(authToken, &userClaims)

	if err != nil {
		validationError, _ := err.(*jwt.ValidationError)
		if validationError.Errors == jwt.ValidationErrorExpired {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			c.Abort()
			return
		}

		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		c.Abort()
		return
	}

	if !token.Valid {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
	}

	if models.IsBanned(userClaims.Id) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "banned",
		})

		c.Abort()
		return
	}

	c.Set("id", userClaims.Id)
	c.Set("is_lec", userClaims.IsLec)
	c.Next()
}

func AdminAuthenticate(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	auths := strings.Split(authToken, "Bearer ")
	if len(auths) < 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		c.Abort()
		return
	}

	authToken = auths[1]
	var adminClaims ultis.AdminClaims
	token, err := ultis.ParseAdminToken(authToken, &adminClaims)

	if err != nil {
		validationError, _ := err.(*jwt.ValidationError)
		if validationError.Errors == jwt.ValidationErrorExpired {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			c.Abort()
			return
		}

		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		c.Abort()
		return
	}

	if !token.Valid {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
	}

	c.Next()
}
