package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model

	Title       string
	Description string
	Location    string
	Price       float32
	Date        time.Time
	Capability  int
	Schedule    string
	Likes       int
	Latitude    float32
	Longitude   float32
	Image       string
}

func CreateEvent(db *gorm.DB, event *Event) error {
	err := db.Create(&event).Error
	return err
}

func GetEvent(db *gorm.DB, id uint) (*Event, error) {
	var event Event
	result := db.First(&event, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

func GetEvents(db *gorm.DB, offset int, limit int) (*[]Event, error) {
	var events []Event
	result := db.Offset(offset).Limit(limit).Find(&events)

	if result.Error != nil {
		return nil, result.Error
	}

	return &events, nil
}

func UpdateEvent(db *gorm.DB, event *Event) error {
	err := db.Save(&event).Error
	return err
}

func DeleteEvent(db *gorm.DB, id uint) error {

	var event Event
	result := db.First(&event, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	err := db.Delete(&event, id).Error

	return err
}
