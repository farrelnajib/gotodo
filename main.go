package main

import (
	"runtime"

	// "github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cache"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/farrelnajib/gotodo/controllers"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	runtime.GOMAXPROCS(2)
	router := router.New()

	router.GET("/activity-groups", controllers.GetActivities)
	router.GET("/activity-groups/{id}", controllers.GetActivitiesById)
	router.POST("/activity-groups", controllers.CreateActivity)
	router.DELETE("/activity-groups/{id}", controllers.DeleteActivity)
	router.PATCH("/activity-groups/{id}", controllers.EditActivity)

	router.GET("/todo-items", controllers.GetAllTodo)
	router.GET("/todo-items/{id}", controllers.GetTodoById)
	router.POST("/todo-items", controllers.CreateTodo)
	router.DELETE("/todo-items/{id}", controllers.DeleteTodo)
	router.PATCH("/todo-items/{id}", controllers.EditTodo)

	fasthttp.ListenAndServe(":3030", router.Handler)
}

// func main() {
// 	runtime.GOMAXPROCS(2)
// 	router := fiber.New(fiber.Config{
// 		DisableStartupMessage: true,
// 	})

// 	// router.Use(logger.New())
// 	router.Use(cache.New())

// 	router.Get("/activity-groups", controllers.GetActivities)
// 	router.Get("/activity-groups/:id", controllers.GetActivitiesById)
// 	router.Post("/activity-groups/", controllers.CreateActivity)
// 	router.Delete("/activity-groups/:id", controllers.DeleteActivity)
// 	router.Patch("/activity-groups/:id", controllers.EditActivity)

// 	router.Get("/todo-items", controllers.GetAllTodo)
// 	router.Get("/todo-items/:id", controllers.GetTodoById)
// 	router.Post("/todo-items", controllers.CreateTodo)
// 	router.Delete("/todo-items/:id", controllers.DeleteTodo)
// 	router.Patch("/todo-items/:id", controllers.EditTodo)

// 	router.Listen(":3030")
// }
