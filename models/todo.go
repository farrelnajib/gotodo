package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/farrelnajib/gotodo/utils"
	"gorm.io/gorm"
)

type Todo struct {
	ID              uint64         `gorm:"primary_key" json:"id"`
	ActivityGroupID uint64         `json:"activity_group_id"`
	Title           string         `json:"title"`
	IsActive        *bool          `gorm:"default:1" json:"is_active"`
	Priority        string         `gorm:"default:'very-high'" json:"priority"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}

func (todo *Todo) ValidateTodo() (utils.Response, bool) {
	if todo.Title == "" {
		return utils.Message("Bad Request", "title cannot be null", map[string]string{}), false
	}
	if todo.ActivityGroupID == 0 {
		return utils.Message("Bad Request", "activity_group_id cannot be null", map[string]string{}), false
	}

	return utils.Message("Success", "Success", map[string]string{}), true
}

// func convertBool(boolean bool) string {
// 	if boolean {
// 		return "1"
// 	}
// 	return "0"
// }

func (todo *Todo) CreateTodo() (utils.Response, int) {
	if response, ok := todo.ValidateTodo(); !ok {
		return response, 400
	}

	GetDB().Create(&todo)

	var data map[string]interface{}
	temp, _ := json.Marshal(&todo)
	json.Unmarshal(temp, &data)

	delete(data, "deleted_at")

	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	return response, 201
}

func GetTodos(activityId uint) []*Todo {
	todos := make([]*Todo, 0)

	query := GetDB()
	if activityId > 0 {
		query = query.Where("activity_group_id = ?", activityId)
	}
	err := query.Find(&todos).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return todos
}

func GetTodoById(id uint) *Todo {
	todo := &Todo{}
	err := GetDB().Where("id = ?", id).First(&todo).Error

	if err != nil {
		return nil
	}

	return todo
}

func DeleteTodo(id uint) (bool, uint64) {
	todo := &Todo{}
	err := GetDB().Where("id = ?", id).First(&todo).Error

	if err != nil {
		return false, 0
	}

	activityId := todo.ActivityGroupID

	err = GetDB().Delete(&todo).Error

	if err != nil {
		fmt.Println(err.Error())
		return false, 0
	}

	return true, activityId
}

func (todo *Todo) EditTodo(id uint) (utils.Response, int, *Todo) {
	existing := &Todo{}

	err := GetDB().Where("id = ?", id).First(&existing).Error
	if err != nil {
		response := utils.Message("Not Found", fmt.Sprintf("Todo with ID %d Not Found", id), map[string]string{})
		return response, 404, nil
	}

	err = GetDB().Model(&existing).Updates(todo).Error
	if err != nil {
		response := utils.Message("Error", err.Error(), map[string]string{})
		return response, 500, nil
	}

	response := utils.Response{Status: "Success", Message: "Success", Data: existing}
	return response, 200, existing
}
