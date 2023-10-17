package controllers

import (
	//	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func UsersRoutes(app *fiber.App) {

	app.Post("/api/v1/users", CreateUser)
	app.Post("/api/v1/login", Login)
	
		app.Use(jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte("secret")},
		}))
	
	// rutas protegidas
	app.Get("/api/v1/users/:id", GetUser)
	app.Put("/api/v1/users/:id", UpdateUser)
	app.Delete("/api/v1/users/:id", DeleteUser)
	app.Put("/api/v1/users/password", UpdateUserPassword)
	app.Get("/api/v1/users", GetUsers)
	// api/v1/users?page=1&limit=10

}

func EventRoutes(app *fiber.App) {

	// rutas protegidas
	app.Post("/api/v1/events", CreateEvent)
	app.Get("/api/v1/events/:id", GetEvent)
	app.Get("/api/v1/events", GetEvents)
	app.Put("/api/v1/events/:id", UpdateEvent)
	app.Delete("/api/v1/events/:id", DeleteEvent)

}
