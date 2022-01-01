package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

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

var GetActivities = func(c *fiber.Ctx) error {
	data := activityCache
	// if len(data) == 0 {
	// 	data = models.GetAllActivities()
	// 	activityCache = data
	// }
	response := utils.Response{Status: "Success", Message: "Success", Data: data}
	return utils.Respond(c, 200, response)
}

var GetActivitiesById = func(c *fiber.Ctx) error {
	params := c.Params("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params), Data: map[string]string{}})
	}

	data := singleActivityCache[uint(id)]
	if data == nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
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
	return utils.Respond(c, 200, response)
}

var CreateActivity = func(c *fiber.Ctx) error {
	activity := &models.Activity{}

	if err := c.BodyParser(activity); err != nil {
		return utils.Respond(c, 400, utils.Response{Status: "Bad Request", Message: ""})
	}

	if response, valid := activity.ValidateActivity(); !valid {
		return utils.Respond(c, 400, response)
	}

	now := time.Now()

	activity.ID = uint64(latestActivityId + 1)
	activity.CreatedAt = now
	activity.UpdatedAt = now

	singleActivityCache[uint(activity.ID)] = activity
	activityCache = append(activityCache, activity)

	activity.CreateActivity()

	latestActivityId++
	return utils.Respond(c, 201, utils.Message("Success", "Success", activity))
}

var DeleteActivity = func(c *fiber.Ctx) error {
	params := c.Params("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params), Data: map[string]string{}})
	}

	activity := singleActivityCache[uint(id)]
	// if activity == nil {
	// 	activity = models.GetActivityById(uint(id))
	// }

	if activity == nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
	}

	go activity.DeleteActivity()
	go DeleteSingleActivityFromCache(id)
	singleActivityCache[uint(id)] = nil

	return utils.Respond(c, 200, utils.Response{Status: "Success", Message: "Success", Data: map[string]string{}})
}

var EditActivity = func(c *fiber.Ctx) error {
	params := c.Params("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params), Data: map[string]string{}})
	}

	activity := &models.Activity{}
	if err := c.BodyParser(activity); err != nil {
		return utils.Respond(c, 400, utils.Response{Status: "Bad request", Message: err.Error()})
	}

	if response, valid := activity.ValidateActivity(); !valid {
		return utils.Respond(c, 400, response)
	}

	existing := singleActivityCache[uint(id)]
	// if existing == nil {
	// 	existing = models.GetActivityById(uint(id))
	// }

	if existing == nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
	}

	existing.Title = activity.Title
	existing.UpdatedAt = time.Now()
	singleActivityCache[uint(id)] = existing

	go activity.EditActivity(existing)
	go EditActivityInCache(existing)

	return utils.Respond(c, 200, utils.Message("Success", "Success", existing))
}
