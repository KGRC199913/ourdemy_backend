package main

import (
	app "github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals"
)

func main() {
	config := app.Config{
		Port:       8080,
		DbUsername: "admin",
		DbPassword: "root",
		DbUrl:      "localhost:27017",
		DbName:     "ourdemy",
	}

	_ = app.Run(&config)
}
