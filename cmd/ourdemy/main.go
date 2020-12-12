package main

import (
	app "github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals"
)

func main() {
	config := app.Config{
		Port:       8080,
		DbUsername: "ABC",
		DbPassword: "ABC",
	}

	_ = app.Run(&config)
}
