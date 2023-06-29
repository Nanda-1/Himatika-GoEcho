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

	// Tambahkan rute-rute lainnya di sini

	e.Start(":8080")
}
