package router

import (
	"echoGo/app/controller"
	"echoGo/db"

	"github.com/labstack/echo/v4"
)

func Router() {
	e := echo.New()
	e.Use(db.DatabaseMiddleware)
	e.POST("/register", controller.Register)
	e.POST("/login", controller.Login)
	e.POST("/refresh-token", controller.RefreshToken)

	//Workshop
	e.POST("/create-workshop", controller.CreateWorkshop)
	e.GET("/getall-workshop", controller.GetAllWorkshop)
	e.POST("/pendaftaran-workshop", controller.CreatePendaftaran)
	e.POST("/absen-workshop", controller.CreateAbsensi)
	e.POST("/props", controller.UploadProposal)

	e.Start(":8080")
}
