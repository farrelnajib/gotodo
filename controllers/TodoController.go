package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/farrelnajib/gotodo/models"
	"github.com/farrelnajib/gotodo/utils"
	"github.com/gorilla/mux"
)

var GetAllTodo = func(w http.ResponseWriter, r *http.Request) {
	params := r.FormValue("activity_group_id")

	activityId, err := strconv.Atoi(params)
	if err != nil {
		activityId = 0
	}

	key := fmt.Sprintf("all_todos_%d", activityId)

	data, err1 := cache.Get(key)
	if err1 == notFound {
		data = models.GetTodos(uint(activityId))
		cache.SetWithTTL(key, data, time.Hour)
	}

	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	utils.Respond(w, 200, response)
}

var GetTodoById = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %s Not Found", params["id"]), Data: map[string]string{}})
		return
	}

	key := fmt.Sprintf("todo_%d", id)

	data, err := cache.Get(key)
	if err != notFound {
		response := utils.Response{Status: "Success", Message: "Success", Data: data}
		utils.Respond(w, 200, response)
		return
	}

	query := models.GetTodoById(uint(id))
	if query == nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	cache.SetWithTTL(key, query, time.Hour)
	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	utils.Respond(w, 200, response)
}

var CreateTodo = func(w http.ResponseWriter, r *http.Request) {
	todo := &models.Todo{}

	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		utils.Respond(w, 400, utils.Response{Status: "Bad request", Message: ""})
		return
	}

	response, status, cachedData := todo.CreateTodo()

	if status == 201 {
		cache.SetWithTTL(fmt.Sprintf("todo_%d", todo.ID), cachedData, time.Hour)
		cache.Remove("all_todos_0")
		cache.Remove(fmt.Sprintf("all_todos_%d", int(todo.ActivityGroupID)))
	}

	utils.Respond(w, status, response)
}

var DeleteTodo = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %s Not Found", params["id"]), Data: map[string]string{}})
		return
	}

	deleted, activityId := models.DeleteTodo(uint(id))
	if !deleted {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	cache.Remove("all_todos_0")
	cache.Remove(fmt.Sprintf("all_todos_%d", activityId))
	cache.Remove(fmt.Sprintf("todo_%d", id))

	utils.Respond(w, 200, utils.Response{Status: "Success", Message: "Success", Data: map[string]string{}})
}

var EditTodo = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %s Not Found", params["id"]), Data: map[string]string{}})
		return
	}

	todo := &models.Todo{}
	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		utils.Respond(w, 400, utils.Response{Status: "Bad request", Message: err.Error()})
		return
	}

	response, status, exsiting := todo.EditTodo(uint(id))

	if status == 200 {
		cache.SetWithTTL(fmt.Sprintf("todo_%d", id), response.Data, time.Hour)
		cache.Remove("all_todos_0")
		cache.Remove(fmt.Sprintf("all_todos_%d", int(exsiting.ActivityGroupID)))
	}

	utils.Respond(w, status, response)
}
