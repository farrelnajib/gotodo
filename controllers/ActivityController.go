package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/farrelnajib/gotodo/models"
	"github.com/farrelnajib/gotodo/utils"
)

var GetActivities = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetAllActivities()
	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	utils.Respond(w, 200, response)
}

var GetActivitiesById = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params["id"]), Data: map[string]string{}})
		return
	}

	data := models.GetActivityById(uint(id))
	if data == nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	utils.Respond(w, 200, response)
}

var CreateActivity = func(w http.ResponseWriter, r *http.Request) {
	activity := &models.Activity{}

	if err := json.NewDecoder(r.Body).Decode(activity); err != nil {
		utils.Respond(w, 400, utils.Response{Status: "Bad request", Message: ""})
		return
	}

	response, status := activity.CreateActivity()
	utils.Respond(w, status, response)
}

var DeleteActivity = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params["id"]), Data: map[string]string{}})
		return
	}

	deleted := models.DeleteActivity(uint(id))
	if !deleted {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	utils.Respond(w, 200, utils.Response{Status: "Success", Message: "Success", Data: map[string]string{}})
}

var EditActivity = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params["id"]), Data: map[string]string{}})
		return
	}

	activity := &models.Activity{}
	if err := json.NewDecoder(r.Body).Decode(activity); err != nil {
		utils.Respond(w, 400, utils.Response{Status: "Bad request", Message: err.Error()})
		return
	}

	response, status := activity.EditActivity(uint(id))
	utils.Respond(w, status, response)
}
