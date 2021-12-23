package models

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func convertBool(boolean bool) string {
	if boolean {
		return "1"
	}
	return "0"
}

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

func GetTodos(activityId uint) []map[string]interface{} {
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

	var data []map[string]interface{}
	temp, _ := json.Marshal(&todos)
	json.Unmarshal(temp, &data)

	for idx, row := range data {
		row["is_active"] = convertBool(*todos[idx].IsActive)
	}

	return data
}

func GetTodoById(id uint) map[string]interface{} {
	todo := &Todo{}
	err := GetDB().Where("id = ?", id).First(&todo).Error

	if err != nil {
		return nil
	}

	var data map[string]interface{}
	temp, _ := json.Marshal(&todo)
	json.Unmarshal(temp, &data)

	data["is_active"] = convertBool(*todo.IsActive)
	data["activity_group_id"] = strconv.Itoa(int(todo.ActivityGroupID))

	return data
}

func DeleteTodo(id uint) bool {
	todo := &Todo{}
	err := GetDB().Where("id = ?", id).First(&todo).Error

	if err != nil {
		return false
	}

	err = GetDB().Delete(&todo).Error

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func (todo *Todo) EditTodo(id uint) (utils.Response, int) {
	existing := &Todo{}

	err := GetDB().Where("id = ?", id).First(&existing).Error
	if err != nil {
		response := utils.Message("Not Found", fmt.Sprintf("Todo with ID %d Not Found", id), map[string]string{})
		return response, 404
	}

	err = GetDB().Model(&existing).Updates(todo).Error
	if err != nil {
		response := utils.Message("Error", err.Error(), map[string]string{})
		return response, 500
	}

	var data map[string]interface{}
	temp, _ := json.Marshal(&existing)
	json.Unmarshal(temp, &data)

	data["is_active"] = convertBool(*existing.IsActive)

	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	return response, 200
}
