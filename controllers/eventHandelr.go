package controllers

import (
	"strconv"

	"github.com/Canibalistico/servicio-as18/database"
	"github.com/Canibalistico/servicio-as18/models"
	"github.com/gofiber/fiber/v2"
)

func CreateEvent(c *fiber.Ctx) error {
	event := &models.Event{}
	if err := c.BodyParser(event); err != nil {
		return fiber.ErrBadRequest
	}

	err := models.CreateEvent(database.DB, event)

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
			"index": " after crete",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Evento registrado exitosamente",
	})

}

func GetEvent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	event, err := models.GetEvent(database.DB, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(event)

}

func GetEvents(c *fiber.Ctx) error {
	// Obtener los parámetros de paginación
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "El parámetro `page` no es válido",
		})
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "El parámetro `limit` no es válido",
		})
	}

	// Crear el cursor
	offset := (page - 1) * limit

	event, err := models.GetEvents(database.DB, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error en la consulta",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(event)
}

func UpdateEvent(c *fiber.Ctx) error {

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	event, err := models.GetEvent(database.DB, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}

	if err := c.BodyParser(&event); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = models.UpdateEvent(database.DB, event)

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registro actualizado exitosamente",
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	err = models.DeleteEvent(database.DB, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Registro eliminado.",
	})
}
