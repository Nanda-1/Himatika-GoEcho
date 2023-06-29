package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func InitDb() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	Db = connectDB()
}

func connectDB() *gorm.DB {
	var err error
	dsn := os.Getenv("DB_DEV_USERNAME") + ":" + os.Getenv("DB_DEV_PASSWORD") + "@tcp" + "(" + os.Getenv("DB_DEV_HOST") + ":" + os.Getenv("DB_DEV_PORT") + ")/" + os.Getenv("DB_DEV_NAME") + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn : ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		fmt.Println("Error connecting to database : error=", err)
		return nil
	}

	return db
}

// Middleware function to add database connection to Echo context
func DatabaseMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("db", Db)
		return next(c)
	}
}

// Function to retrieve database connection from Echo context
func GetDB(c echo.Context) *gorm.DB {
	db, _ := c.Get("db").(*gorm.DB)
	return db
}
