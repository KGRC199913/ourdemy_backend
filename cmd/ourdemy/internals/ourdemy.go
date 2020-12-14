package internals

import (
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type Config struct {
	Port       int32
	DbUsername string
	DbPassword string
	DbUrl      string
	DbName     string
}

func dbInit(config *Config) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", config.DbUsername, config.DbPassword, config.DbUrl)
	fmt.Println(uri)
	err := mgm.SetDefaultConfig(nil, config.DbName, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	return nil
}

func serverInit(config *Config) (*gin.Engine, error) {
	//init DB here
	err := dbInit(config)
	if err != nil {
		return nil, err
	}
	//end DB init

	r := gin.New()
	//add global mdw here
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	//end global mdw

	return r, nil
}

func Run(config *Config) error {
	r, err := serverInit(config)
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK!")
	})

	r.GET("/test", func(c *gin.Context) {
		testUser := models.NewUser("ABC", "ABC@Meow.com", "123")
		err := mgm.Coll(testUser).Create(testUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "ok!")
		}
	})

	r.GET("/testRes", func(c *gin.Context) {
		user := &models.User{}
		coll := mgm.Coll(user)

		err := coll.First(bson.M{"fullname": "ABC"}, user)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, "Err")
		} else {
			c.JSON(http.StatusOK, user)
		}
	})

	err = r.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}
	return nil
}
