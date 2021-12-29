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

var todoCache = map[string][]*models.Todo{}
var singleTodoCache = map[string]*models.Todo{}

var GetAllTodo = func(w http.ResponseWriter, r *http.Request) {
	params := r.FormValue("activity_group_id")

	activityId, err := strconv.Atoi(params)
	if err != nil {
		activityId = 0
	}

	key := fmt.Sprintf("all_todos_%d", activityId)

	data, isFound := todoCache[key]
	if !isFound {
		data = models.GetTodos(uint(activityId))
		todoCache[key] = data
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

	data, isFound := singleTodoCache[key]
	if isFound && data != nil {
		response := utils.Response{Status: "Success", Message: "Success", Data: data}
		utils.Respond(w, 200, response)
		return
	}

	query := models.GetTodoById(uint(id))
	if query == nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Todo with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	singleTodoCache[key] = query
	response := utils.Response{Status: "Success", Message: "Success", Data: query}
	utils.Respond(w, 200, response)
}

var CreateTodo = func(w http.ResponseWriter, r *http.Request) {
	todo := &models.Todo{}

	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		utils.Respond(w, 400, utils.Response{Status: "Bad request", Message: ""})
		return
	}

	response, status := todo.CreateTodo()

	if status == 201 {
		key := fmt.Sprintf("todo_%d", todo.ID)
		singleTodoCache[key] = todo

		if _, isFound := todoCache["all_todos_0"]; isFound {
			todoCache["all_todos_0"] = append(todoCache["all_todos_0"], todo)
		} else {
			todoCache["all_todos_0"] = []*models.Todo{todo}
		}

		key = fmt.Sprintf("all_todos_%d", int(todo.ActivityGroupID))
		if _, isFound := todoCache[key]; isFound {
			todoCache[key] = append(todoCache[key], todo)
		} else {
			todoCache[key] = []*models.Todo{todo}
		}
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

	if _, isFound := todoCache["all_todos_0"]; isFound {
		todoCache["all_todos_0"] = nil
	}

	key := fmt.Sprintf("all_todos_%d", activityId)
	if _, isFound := todoCache[key]; isFound {
		todoCache[key] = nil
	}

	key = fmt.Sprintf("todo_%d", id)
	if _, isFound := singleTodoCache[key]; isFound {
		singleTodoCache[key] = nil
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

	response, status, exsiting := todo.EditTodo(uint(id))

	if status == 200 {
		if _, isFound := todoCache["all_todos_0"]; isFound {
			todoCache["all_todos_0"] = nil
		}

		key := fmt.Sprintf("all_todos_%d", exsiting.ActivityGroupID)
		if _, isFound := todoCache[key]; isFound {
			todoCache[key] = nil
		}

		key = fmt.Sprintf("todo_%d", id)
		if _, isFound := singleTodoCache[key]; isFound {
			singleTodoCache[key] = exsiting
		}
	}

	utils.Respond(w, status, response)
}
