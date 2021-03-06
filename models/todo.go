package models

import (
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

func (todo *Todo) CreateTodo() {
	GetDB().Create(&todo)
}

func GetTodos(activityId uint) []*Todo {
	todos := make([]*Todo, 0)

	query := GetDB()
	if activityId > 0 {
		query = query.Where("activity_group_id = ?", activityId)
	}
	err := query.Where("is_active", 1).Find(&todos).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return todos
}

func GetTodoById(id uint) *Todo {
	todo := &Todo{}
	err := GetDB().First(&todo, id).Error

	if err != nil {
		return nil
	}

	return todo
}

func (todo *Todo) DeleteTodo() {
	GetDB().Delete(&todo)
}

func (todo *Todo) EditTodo(existing *Todo) {
	GetDB().Model(&existing).Updates(todo)
}
