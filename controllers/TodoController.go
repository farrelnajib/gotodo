package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	data := models.GetTodos(uint(activityId))
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

	data := models.GetTodoById(uint(id))
	if data == nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	utils.Respond(w, 200, response)
}

var CreateTodo = func(w http.ResponseWriter, r *http.Request) {
	todo := &models.Todo{}

	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		utils.Respond(w, 400, utils.Response{Status: "Bad request", Message: ""})
		return
	}

	response, status := todo.CreateTodo()
	utils.Respond(w, status, response)
}

var DeleteTodo = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %s Not Found", params["id"]), Data: map[string]string{}})
		return
	}

	deleted := models.DeleteTodo(uint(id))
	if !deleted {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

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

	response, status := todo.EditTodo(uint(id))
	utils.Respond(w, status, response)
}
