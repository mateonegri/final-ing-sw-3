package main

import (
	"user-api/app"
	"user-api/db"
)

func main() {
	db.StartDbEngine()

	app.StartRoute()

}
