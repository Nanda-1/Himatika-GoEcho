package main

import (
	"echoGo/app/models"
	"echoGo/app/router"
	"echoGo/db"
)

func main() {
	db.InitDb()
	// db.Db.AutoMigrate(&models.Anggota{}, &models.AnggotaDetails{}, &models.Roles{})
	db.Db.AutoMigrate(&models.Workshop{}, &models.AbsensiWorkshop{}, &models.PendaftaranWorkshop{},
		&models.PropWorkshop{}, &models.LpjWorkshop{}, &models.DokumentasiWorkshop{})
	router.Router()
}
