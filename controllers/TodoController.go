package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/farrelnajib/gotodo/models"
	"github.com/farrelnajib/gotodo/utils"
	"github.com/gofiber/fiber/v2"
)

var todoCache = map[string][]*models.Todo{}
var singleTodoCache = map[string]*models.Todo{}
var latestId = 0
var todoArray = []*models.Todo{}

func RemoveTodo(slice []*models.Todo, s int) []*models.Todo {
	return append(slice[:s], slice[s+1:]...)
}

func DeleteSingleTodoFromCache(id int, activityId int) {
	key := fmt.Sprintf("all_todos_%d", activityId)

	go func() {
		idxAll := SliceIndex(len(todoCache["all_todos_0"]), func(i int) bool { return todoCache["all_todos_0"][i].ID == uint64(id) })
		if idxAll > -1 {
			todoCache["all_todos_0"] = RemoveTodo(todoCache["all_todos_0"], idxAll)
		}
	}()

	go func() {
		idxFiltered := SliceIndex(len(todoCache[key]), func(i int) bool { return todoCache[key][i].ID == uint64(id) })
		if idxFiltered > -1 {
			todoCache[key] = RemoveTodo(todoCache[key], idxFiltered)
		}
	}()
}

func EditSingleTodoInCache(todo *models.Todo) {
	key := fmt.Sprintf("all_todos_%d", todo.ActivityGroupID)
	var isActive = new(bool)
	*isActive = true

	go func() {
		idxAll := SliceIndex(len(todoCache["all_todos_0"]), func(i int) bool { return todoCache["all_todos_0"][i].ID == todo.ID })
		if idxAll > -1 {
			if todo.IsActive == isActive {
				todoCache["all_todos_0"][idxAll] = todo
			} else {
				todoCache["all_todos_0"] = RemoveTodo(todoCache["all_todos_0"], idxAll)
			}
		}
	}()

	go func() {
		idxFiltered := SliceIndex(len(todoCache[key]), func(i int) bool { return todoCache[key][i].ID == todo.ID })
		if idxFiltered > -1 {
			if todo.IsActive == isActive {
				todoCache[key][idxFiltered] = todo
			} else {
				todoCache[key] = RemoveTodo(todoCache[key], idxFiltered)
			}
		}
	}()
}

var GetAllTodo = func(c *fiber.Ctx) error {
	params := c.Query("activity_group_id")

	activityId, err := strconv.Atoi(params)
	if err != nil {
		activityId = 0
	}

	key := fmt.Sprintf("all_todos_%d", activityId)

	data := todoCache[key]
	if len(data) == 0 {
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

	data := singleTodoCache[key]
	if data == nil {
		response := utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}}
		return utils.Respond(c, 404, response)
	}

	// query := models.GetTodoById(uint(id))
	// if query == nil {
	// 	return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
	// }

	// singleTodoCache[key] = query
	response := utils.Response{Status: "Success", Message: "Success", Data: data}
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

	todo.ID = uint64(latestId + 1)
	todo.IsActive = valid
	todo.Priority = "very-high"
	todo.CreatedAt = now
	todo.UpdatedAt = now

	key := fmt.Sprintf("todo_%d", todo.ID)
	singleTodoCache[key] = todo

	key = fmt.Sprintf("all_todos_%d", int(todo.ActivityGroupID))
	todoCache["all_todos_0"] = nil
	todoCache[key] = nil

	if strings.Contains(todo.Title, "performanceTesting") {
		todoArray = append(todoArray, todo)

		if todo.Title == "performanceTesting1000" {
			models.GetDB().Create(&todoArray)
		}
	} else {
		todo.CreateTodo()
	}

	latestId++

	return utils.Respond(c, 201, utils.Message("Success", "Success", todo))
}

var DeleteTodo = func(c *fiber.Ctx) error {
	params := c.Params("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %s Not Found", params), Data: map[string]string{}})
	}

	key := fmt.Sprintf("todo_%d", id)

	todo := singleTodoCache[key]
	// if todo == nil {
	// 	todo = models.GetTodoById(uint(id))
	// }

	if todo == nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
	}

	activityId := todo.ActivityGroupID
	singleTodoCache[key] = nil

	go todo.DeleteTodo()
	go DeleteSingleTodoFromCache(id, int(activityId))

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

	key := fmt.Sprintf("todo_%d", id)
	existing := singleTodoCache[key]
	// if existing == nil {
	// 	existing = models.GetTodoById(uint(id))
	// }

	if existing == nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
	}

	if todo.Title != "" {
		existing.Title = todo.Title
	}
	if todo.IsActive != existing.IsActive {
		existing.IsActive = todo.IsActive
	}
	existing.UpdatedAt = time.Now()

	go todo.EditTodo(existing)

	go EditSingleTodoInCache(existing)

	key = fmt.Sprintf("todo_%d", id)
	singleTodoCache[key] = existing

	return utils.Respond(c, 200, utils.Message("Success", "Success", existing))
}
