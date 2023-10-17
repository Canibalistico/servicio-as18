package main

import (
	"github.com/Canibalistico/servicio-as18/controllers"
	"github.com/Canibalistico/servicio-as18/database"
	"github.com/Canibalistico/servicio-as18/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.DBconnection()
	database.DB.AutoMigrate(models.User{})
	database.DB.AutoMigrate(models.Event{})

	app := fiber.New()
	app.Use(cors.New())

	controllers.UsersRoutes(app)
	controllers.EventRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
