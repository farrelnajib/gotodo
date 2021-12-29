package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/farrelnajib/gotodo/models"
	"github.com/farrelnajib/gotodo/utils"
)

var activityCache = []*models.Activity{}
var singleActivityCache = map[uint]*models.Activity{}

var GetActivities = func(c *fiber.Ctx) error {
	data := activityCache
	if len(data) == 0 {
		data = models.GetAllActivities()
		activityCache = data
	}
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
	if data != nil {
		response := utils.Response{Status: "Success", Message: "Success", Data: data}
		return utils.Respond(c, 200, response)
	}

	query := models.GetActivityById(uint(id))

	if query == nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
	}

	singleActivityCache[uint(id)] = query

	response := utils.Response{Status: "Success", Message: "Success", Data: query}
	return utils.Respond(c, 200, response)
}

var CreateActivity = func(c *fiber.Ctx) error {
	activity := &models.Activity{}

	if err := c.BodyParser(activity); err != nil {
		return utils.Respond(c, 400, utils.Response{Status: "Bad Request", Message: ""})
	}

	response, status := activity.CreateActivity()
	if status == 201 {
		singleActivityCache[uint(activity.ID)] = activity
		activityCache = append(activityCache, activity)
	}
	return utils.Respond(c, status, response)
}

var DeleteActivity = func(c *fiber.Ctx) error {
	params := c.Params("id")
	id, err := strconv.Atoi(params)
	if err != nil {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %s Not Found", params), Data: map[string]string{}})
	}

	deleted := models.DeleteActivity(uint(id))
	if !deleted {
		return utils.Respond(c, 404, utils.Response{Status: "Not Found", Message: fmt.Sprintf("Activity with ID %d Not Found", id), Data: map[string]string{}})
	}

	activityCache = []*models.Activity{}
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

	response, status, data := activity.EditActivity(uint(id))

	if status == 200 {
		singleActivityCache[uint(id)] = data
		activityCache = []*models.Activity{}
	}

	return utils.Respond(c, status, response)
}
