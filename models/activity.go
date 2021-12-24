package models

import (
	"encoding/json"
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

func (activity *Activity) CreateActivity() (utils.Response, int) {
	if response, ok := activity.ValidateActivity(); !ok {
		return response, 400
	}

	GetDB().Create(activity)

	var data map[string]interface{}
	temp, _ := json.Marshal(&activity)
	json.Unmarshal(temp, &data)

	delete(data, "deleted_at")

	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	return response, 201
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

func DeleteActivity(id uint) bool {
	activity := GetActivityById(id)
	if activity == nil {
		return false
	}

	err := GetDB().Delete(&activity).Error
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (activity *Activity) EditActivity(id uint) (utils.Response, int, *Activity) {
	if response, ok := activity.ValidateActivity(); !ok {
		return response, 400, nil
	}

	existing := GetActivityById(id)
	if existing == nil {
		response := utils.Message("Not Found", fmt.Sprintf("Activity with ID %d Not Found", id), map[string]string{})
		return response, 404, nil
	}

	err := GetDB().Model(&existing).Updates(activity).Error
	if err != nil {
		response := utils.Message("Error", err.Error(), map[string]string{})
		return response, 500, nil
	}

	response := utils.Response{Status: "Success", Message: "Success", Data: existing}
	return response, 200, existing
}
