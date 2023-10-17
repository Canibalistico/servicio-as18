package controllers

import (
	"errors"
	"strconv"
	"time"
	"unicode"

	"github.com/Canibalistico/servicio-as18/database"
	"github.com/Canibalistico/servicio-as18/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {

	type UserLogin struct {
		Email    string
		Password string
	}

	data := UserLogin{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := models.GetUserByEmail(database.DB, data.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "El password es incorrecto",
			"msg":     err.Error(),
		})
	}

	claims := jwt.MapClaims{
		"name":  user.Email,
		"admin": false,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": t,
	})
}

func GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	user, err := models.GetUser(database.DB, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(user)

}

func GetUsers(c *fiber.Ctx) error {
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

	user, err := models.GetUsers(database.DB, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error en la consulta",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	user := &models.User{}
	if err := c.BodyParser(user); err != nil {
		return fiber.ErrBadRequest
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = models.CreateUser(database.DB, user)

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Usuario registrado exitosamente",
	})

}

func UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	user, err := models.GetUser(database.DB, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}

	type UpdateUser struct {
		FirstName string
		LastName  string
		Email     string
		Phone     string
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName
	user.Email = updateData.Email
	user.Phone = updateData.Phone

	err = models.UpdateUser(database.DB, user)

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Usuario actualizado exitosamente",
	})
}

func UpdateUserPassword(c *fiber.Ctx) error {
	type UserPassword struct {
		Email                 string
		Password              string
		Password_Confirmation string
	}

	data := UserPassword{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := models.GetUserByEmail(database.DB, data.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "El password es incorrecto",
			"msg":     err.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password_Confirmation), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	err = models.UpdateUserPassword(database.DB, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password actualizado exitosamente",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	err = models.DeleteUser(database.DB, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Registro eliminado.",
	})
}

func ValidatePassword(password string) error {
	// Verificar que el password tenga más de 4 caracteres
	if len(password) < 4 {
		return errors.New("El password debe tener más de 4 caracteres")
	}

	// Verificar que el password contenga letras y números
	for _, char := range password {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return errors.New("El password debe contener letras y números")
		}
	}

	return nil
}
