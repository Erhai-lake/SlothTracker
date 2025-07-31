package main

import (
	"sloth-tracker/api/router"
	"sloth-tracker/api/storage"
)

func main() {
	db := storage.InitDB()
	r := router.SetupRouter(db)
	r.Run(":8080")
}
