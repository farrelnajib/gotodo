package main

import (
	"fmt"
	"net/http"

	"github.com/farrelnajib/gotodo/controllers"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", controllers.GetActivities)
	router.HandleFunc("/activity-groups", controllers.GetActivities).Methods("GET")
	router.HandleFunc("/activity-groups", controllers.CreateActivity).Methods("POST")
	router.HandleFunc("/activity-groups/{id}", controllers.GetActivitiesById).Methods("GET")
	router.HandleFunc("/activity-groups/{id}", controllers.DeleteActivity).Methods("DELETE")
	router.HandleFunc("/activity-groups/{id}", controllers.EditActivity).Methods("PATCH")

	router.HandleFunc("/todo-items", controllers.GetAllTodo).Methods("GET")
	router.HandleFunc("/todo-items", controllers.CreateTodo).Methods("POST")
	router.HandleFunc("/todo-items/{id}", controllers.GetTodoById).Methods("GET")
	router.HandleFunc("/todo-items/{id}", controllers.DeleteTodo).Methods("DELETE")
	router.HandleFunc("/todo-items/{id}", controllers.EditTodo).Methods("PATCH")

	router.Use(middleware.Logger)

	err := http.ListenAndServe(":3030", router)
	if err != nil {
		fmt.Println(err.Error())
	}
}
