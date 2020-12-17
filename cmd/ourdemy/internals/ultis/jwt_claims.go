package ultis

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserClaims struct {
	Id primitive.ObjectID
	jwt.StandardClaims
}

func CreateToken(oid primitive.ObjectID) (string, error) {
	userClaims := UserClaims{
		Id: oid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ParseToken(authToken string, claims *UserClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong method")
		}
		return []byte(secret), nil
	})
}
