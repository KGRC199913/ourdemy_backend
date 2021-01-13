package internals

import (
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/models"
	route "github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/routes"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"time"
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

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	//cors
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Refresh", "Accept", "Accept-Language", "Content-Type"},
		ExposeHeaders:   []string{"AccessToken", "RefreshToken"},
		MaxAge:          12 * time.Hour,
	}))
	//end cors
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
	route.TimeMark(r)
	route.HomeRoutes(r)
}

func Run(config *ultis.Config) error {
	r, err := serverInit(config)
	if err != nil {
		panic(err)
	}

	routing(r)

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("ourdemy.xyz"),
		Cache:      autocert.DirCache("/certCache"),
	}

	return autotls.RunWithManager(r, &m)
}
