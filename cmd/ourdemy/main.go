package main

import (
	app "github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	"log"
)

func main() {
	//config := ultis.Config{
	//	Port:       8080,
	//	DbUsername: "admin",
	//	DbPassword: "root",
	//	DbUrl:      "localhost:27017",
	//	DbName:     "ourdemy",
	//}

	config, _ := ultis.LoadConfigFile()
	log.Fatal(app.Run(config))
}
