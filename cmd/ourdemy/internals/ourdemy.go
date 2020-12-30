package internals

import (
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	route "github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/routes"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	"github.com/gin-gonic/gin"
)

func dbInit(config *ultis.Config) error {
	//uri := fmt.Sprintf("mongodb://%s:%s@%s", config.json.DbUsername, config.json.DbPassword, config.json.DbUrl)
	//fmt.Println(uri)
	//err := mgm.SetDefaultConfig(nil, config.json.DbName, options.Client().ApplyURI(uri))
	//if err != nil {
	//	return err
	//}

	return models.InitDb(config)
}

func serverInit(config *ultis.Config) (*gin.Engine, error) {
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

func routing(r *gin.Engine) {
	route.TestRoutes(r)
	route.CategoryRoutes(r)
	route.CourseRoutes(r)
	route.LecturerRoutes(r)
	route.AdminRoutes(r)
	route.UserRoutes(r)
	route.ReviewRoutes(r)
	route.VideoRoutes(r)
}

func Run(config *ultis.Config) error {
	r, err := serverInit(config)
	if err != nil {
		panic(err)
	}

	routing(r)

	fmt.Printf("App started, listentning on port %s\n", config.Port)
	err = r.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}
	return nil
}
