package internals

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Config struct {
	Port       int32
	DbUsername string
	DbPassword string
}

func serverInit(config *Config) (*gin.Engine, error) {
	//init DB here

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

	err = r.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}
	return nil
}
