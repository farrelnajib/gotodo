package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/farrelnajib/gotodo/models"
	"github.com/farrelnajib/gotodo/utils"
	"github.com/gofiber/fiber/v2"
)

var todoCache = map[string][]*models.Todo{}
var singleTodoCache = map[string]*models.Todo{}

var GetAllTodo = func(c *fiber.Ctx) error {
	params := c.Query("activity_group_id")

	activityId, err := strconv.Atoi(params)
	if err != nil {
		activityId = 0
	}

	key := fmt.Sprintf("all_todos_%d", activityId)

	data, _ := todoCache[key]
	if data == nil || len(data) == 0 {
		data = models.GetTodos(uint(activityId))
		todoCache[key] = data
	}

	response := utils.Message("Success", "Success", data)
	return utils.Respond(c, 200, response)
}

var GetTodoById = func(c *fiber.Ctx) error {
	params := c.Params("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %s Not Found", params), Data: map[string]string{}})
	}

	key := fmt.Sprintf("todo_%d", id)

	data, isFound := singleTodoCache[key]
	if isFound && data != nil {
		response := utils.Response{Status: "Success", Message: "Success", Data: data}
		return utils.Respond(c, 200, response)
	}

	query := models.GetTodoById(uint(id))
	if query == nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
	}

	singleTodoCache[key] = query
	response := utils.Response{Status: "Success", Message: "Success", Data: query}
	return utils.Respond(c, 200, response)
}

var CreateTodo = func(c *fiber.Ctx) error {
	todo := &models.Todo{}

	if err := c.BodyParser(todo); err != nil {
		return utils.Respond(c, 400, utils.Response{Status: "Bad Request", Message: ""})
	}

	response, isValid := todo.ValidateTodo()
	if !isValid {
		return utils.Respond(c, 400, response)
	}

	now := time.Now()
	valid := new(bool)
	*valid = true

	todo.ID = uint64(len(todoCache["all_todos_0"]) + 1)
	todo.IsActive = valid
	todo.Priority = "very-high"
	todo.CreatedAt = now
	todo.UpdatedAt = now

	key := fmt.Sprintf("todo_%d", todo.ID)
	singleTodoCache[key] = todo

	if cache, _ := todoCache["all_todos_0"]; cache == nil || len(cache) == 0 {
		todoCache["all_todos_0"] = []*models.Todo{todo}
	} else {
		todoCache["all_todos_0"] = append(todoCache["all_todos_0"], todo)
	}

	key = fmt.Sprintf("all_todos_%d", int(todo.ActivityGroupID))
	if cache, _ := todoCache[key]; cache == nil || len(cache) == 0 {
		todoCache[key] = []*models.Todo{todo}
	} else {
		todoCache[key] = append(todoCache[key], todo)
	}

	go todo.CreateTodo()

	return utils.Respond(c, 201, utils.Message("Success", "Success", todo))
}

var DeleteTodo = func(c *fiber.Ctx) error {
	params := c.Params("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %s Not Found", params), Data: map[string]string{}})
	}

	deleted, activityId := models.DeleteTodo(uint(id))
	if !deleted {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
	}

	todoCache["all_todos_0"] = []*models.Todo{}

	key := fmt.Sprintf("all_todos_%d", activityId)
	todoCache[key] = []*models.Todo{}

	key = fmt.Sprintf("todo_%d", id)
	singleTodoCache[key] = nil

	return utils.Respond(c, 200, utils.Response{Status: "Success", Message: "Success", Data: map[string]string{}})
}

var EditTodo = func(c *fiber.Ctx) error {
	params := c.Params("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %s Not Found", params), Data: map[string]string{}})
	}

	todo := &models.Todo{}
	if err := c.BodyParser(todo); err != nil {
		return utils.Respond(c, 400, utils.Response{Status: "Bad Request", Message: err.Error()})
	}

	response, status, exsiting := todo.EditTodo(uint(id))

	if status == 200 {
		todoCache = map[string][]*models.Todo{}

		key := fmt.Sprintf("todo_%d", id)
		singleTodoCache[key] = exsiting
	}

	return utils.Respond(c, status, response)
}
