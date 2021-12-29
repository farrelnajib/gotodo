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

// var cache ttlcache.SimpleCache = ttlcache.NewCache()
// var notFound = ttlcache.ErrNotFound

var activityCache = []*models.Activity{}
var singleActivityCache = map[uint]*models.Activity{}

var GetActivities = func(w http.ResponseWriter, r *http.Request) {
	data := activityCache
	if len(data) == 0 {
		data = models.GetAllActivities()
		activityCache = data
	}
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

	data := singleActivityCache[uint(id)]
	if data != nil {
		response := utils.Response{Status: "Success", Message: "Success", Data: data}
		utils.Respond(w, 200, response)
		return
	}

	query := models.GetActivityById(uint(id))

	if query == nil {
		utils.Respond(w, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	singleActivityCache[uint(id)] = query

	response := utils.Response{Status: "Success", Message: "Success", Data: query}
	utils.Respond(w, 200, response)
}

var CreateActivity = func(w http.ResponseWriter, r *http.Request) {
	activity := &models.Activity{}

	if err := json.NewDecoder(r.Body).Decode(activity); err != nil {
		utils.Respond(w, 400, utils.Response{Status: "Bad request", Message: ""})
		return
	}

	response, status := activity.CreateActivity()
	if status == 201 {
		singleActivityCache[uint(activity.ID)] = activity
		activityCache = append(activityCache, activity)
	}
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

	activityCache = []*models.Activity{}
	singleActivityCache[uint(id)] = nil

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

	response, status, data := activity.EditActivity(uint(id))

	if status == 200 {
		singleActivityCache[uint(id)] = data
		activityCache = []*models.Activity{}
	}

	utils.Respond(w, status, response)
}
