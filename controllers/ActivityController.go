package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/farrelnajib/gotodo/models"
	"github.com/farrelnajib/gotodo/utils"
)

var activityCache = []*models.Activity{}
var singleActivityCache = map[uint]*models.Activity{}
var latestActivityId = 0

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func RemoveActivity(slice []*models.Activity, s int) []*models.Activity {
	return append(slice[:s], slice[s+1:]...)
}

func DeleteSingleActivityFromCache(id int) {
	idx := SliceIndex(len(activityCache), func(i int) bool { return activityCache[i].ID == uint64(id) })
	activityCache = RemoveActivity(activityCache, idx)
}

func EditActivityInCache(activity *models.Activity) {
	idx := SliceIndex(len(activityCache), func(i int) bool { return activityCache[i].ID == activity.ID })
	activityCache[idx] = activity
}

var GetActivities = func(c *fasthttp.RequestCtx) {
	data := activityCache
	// if len(data) == 0 {
	// 	data = models.GetAllActivities()
	// 	activityCache = data
	// }
	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	utils.Respond(c, 200, response)
}

var GetActivitiesById = func(c *fasthttp.RequestCtx) {
	params := fmt.Sprintf("%s", c.UserValue("id"))
	id, err := strconv.Atoi(params)
	if err != nil {
		utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params), Data: map[string]string{}})
		return
	}

	data := singleActivityCache[uint(id)]
	if data == nil {
		utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
		return
	}
	// if data != nil {
	// 	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	// 	return utils.Respond(c, 200, response)
	// }

	// query := models.GetActivityById(uint(id))

	// if query == nil {
	// 	return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
	// }

	// singleActivityCache[uint(id)] = query

	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	utils.Respond(c, 200, response)
}

var CreateActivity = func(c *fasthttp.RequestCtx) {
	activity := &models.Activity{}

	if err := json.Unmarshal(c.PostBody(), &activity); err != nil {
		utils.Respond(c, 400, utils.Response{Status: "Bad Request", Message: ""})
		return
	}

	if response, valid := activity.ValidateActivity(); !valid {
		utils.Respond(c, 400, response)
		return
	}

	now := time.Now()

	activity.ID = uint64(latestActivityId + 1)
	activity.CreatedAt = now
	activity.UpdatedAt = now

	go func() {
		singleActivityCache[uint(activity.ID)] = activity
		activityCache = append(activityCache, activity)
	}()

	activity.CreateActivity()

	latestActivityId++
	utils.Respond(c, 201, utils.Message("Success", "Success", activity))
}

var DeleteActivity = func(c *fasthttp.RequestCtx) {
	params := fmt.Sprintf("%s", c.UserValue("id"))
	id, err := strconv.Atoi(params)
	if err != nil {
		utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params), Data: map[string]string{}})
		return
	}

	activity := singleActivityCache[uint(id)]
	// if activity == nil {
	// 	activity = models.GetActivityById(uint(id))
	// }

	if activity == nil {
		utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	go activity.DeleteActivity()
	go DeleteSingleActivityFromCache(id)
	go func() {
		singleActivityCache[uint(id)] = nil
	}()

	utils.Respond(c, 200, utils.Response{Status: "Success", Message: "Success", Data: map[string]string{}})
}

var EditActivity = func(c *fasthttp.RequestCtx) {
	params := fmt.Sprintf("%s", c.UserValue("id"))
	id, err := strconv.Atoi(params)
	if err != nil {
		utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params), Data: map[string]string{}})
		return
	}

	activity := &models.Activity{}
	if err := json.Unmarshal(c.PostBody(), &activity); err != nil {
		utils.Respond(c, 400, utils.Response{Status: "Bad request", Message: err.Error()})
		return
	}

	if response, valid := activity.ValidateActivity(); !valid {
		utils.Respond(c, 400, response)
		return
	}

	existing := singleActivityCache[uint(id)]
	// if existing == nil {
	// 	existing = models.GetActivityById(uint(id))
	// }

	if existing == nil {
		utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
		return
	}

	existing.Title = activity.Title
	existing.UpdatedAt = time.Now()

	go activity.EditActivity(existing)
	go EditActivityInCache(existing)
	go func() {
		singleActivityCache[uint(id)] = existing
	}()

	utils.Respond(c, 200, utils.Message("Success", "Success", existing))
}
