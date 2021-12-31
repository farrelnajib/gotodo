package models

import (
	"fmt"
	"time"

	"github.com/farrelnajib/gotodo/utils"
	"gorm.io/gorm"
)

type Activity struct {
	ID        uint64         `gorm:"primary_key" json:"id"`
	Email     string         `json:"email"`
	Title     string         `json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (activity *Activity) ValidateActivity() (utils.Response, bool) {
	if activity.Title == "" {
		return utils.Message("Bad Request", "title cannot be null", map[string]string{}), false
	}

	return utils.Message("Success", "Success", map[string]string{}), true
}

func (activity *Activity) CreateActivity() {
	GetDB().Create(activity)
}

func GetAllActivities() []*Activity {
	activities := make([]*Activity, 0)
	err := GetDB().Find(&activities).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return activities
}

func GetActivityById(id uint) *Activity {
	activity := &Activity{}
	err := GetDB().Where("id = ?", id).First(&activity).Error

	if err != nil {
		return nil
	}

	return activity
}

func (activity *Activity) DeleteActivity() {
	GetDB().Delete(&activity)
}

func (activity *Activity) EditActivity(existing *Activity) {
	GetDB().Model(&existing).Updates(activity)
}
