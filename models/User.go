package models

import (
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	gorm.Model

	FirstName string                `gorm:"type:varchar(100);not null"`
	LastName  string                `gorm:"type:varchar(100);not null"`
	Email     string                `gorm:"type:varchar(100);not null;unique"`
	Phone     string                `gorm:"type:varchar(100)"`
	Password  string                `gorm:"type:varchar(255);not null"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func GetUser(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.Select("id", "FirstName", "LastName", "Email", "Phone").First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func GetUsers(db *gorm.DB, offset int, limit int) (*[]User, error) {
	var users []User
	result := db.Select("id", "FirstName", "LastName", "Email", "Phone").Offset(offset).Limit(limit).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return &users, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	err := db.Create(&user).Error
	return err
}

func UpdateUser(db *gorm.DB, user *User) error {
	err := db.Save(&user).Error
	return err
}

func UpdateUserPassword(db *gorm.DB, user *User) error {
	err := db.Model(&user).Update("password", user.Password).Error
	return err
}

func DeleteUser(db *gorm.DB, id uint) error {

	var user User
	result := db.First(&user, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	err := db.Delete(&user, id).Error

	return err
}

/// https://github.com/ftfetter/gin-gorm-mysql/
// https://github.com/rog-golang-buddies/go-automatic-apps/blob/main/architecture/show_model.puml
// http://www.plantuml.com/plantuml/duml/
