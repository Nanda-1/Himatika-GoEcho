package main

import (
	"echoGo/app/router"
	"echoGo/db"
)

func main() {
	db.InitDb()
	// db.Db.AutoMigrate(&models.Anggota{}, &models.AnggotaDetails{}, &models.Roles{})
	router.Router()
}
