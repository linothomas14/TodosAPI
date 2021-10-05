package main

import (
	"encoding/json"
	_ "final-project-1/docs"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Todos struct {
	ID          int    `json:"id" example:"1"`
	Title       string `json:"title" example:"Reading book"`
	Description string `json:"description" example:"Reading book at 9 A.M"`
	IsComplete  string `json:"is_complete" example:"false"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error"`
}

var prevOrderID = 0
var all_todos []Todos
var error ErrorResponse

// @title Todos API
// @version 1.0
// @description This is a sample service for managing todos
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email hacktiv@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/todos", GetAllTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", GetTodos).Methods("GET")
	router.HandleFunc("/todos", CreateTodos).Methods("POST")
	router.HandleFunc("/todos/{id}", UpdateTodos).Methods("PUT")
	router.HandleFunc("/todos/{id}", DeleteTodos).Methods("DELETE")
	log.Println("Start server at localhost:8080")
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	http.ListenAndServe(":8080", router)
}

// GetAllTodos godoc
// @Summary Get all todos
// @Description Get details of all todos
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} Todos
// @Router /todos [get]
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func GetAllTodos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(all_todos)

}

// CreateTodos godoc
// @Summary Create a todo
// @Description Create a new todo with the input payload
// @Tags todos
// @Accept json
// @Produce json
// @Param todos body Todos true "Create todo"
// @Success 200 {array} Todos
// @Router /todos [post]
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func CreateTodos(w http.ResponseWriter, r *http.Request) {
	var todos Todos
	w.Header().Set("Content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&todos)
	if err != nil {
		fmt.Println("ERROR")
		log.Fatal(err)
	}
	prevOrderID++
	todos.ID = prevOrderID
	all_todos = append(all_todos, todos)
	json.NewEncoder(w).Encode(todos)
}

// GetTodos godoc
// @Summary Get a todo
// @Description Get a todo
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID of the todos"
// @Success 200 {object} Todos
// @Router /todos/{id} [get]
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ID, err := strconv.Atoi(params["id"])

	if err != nil {
		msg := fmt.Sprintf("id must be integer")
		error.Error = msg
		json.NewEncoder(w).Encode(error)
		return
	}
	for _, val := range all_todos {
		if val.ID == ID {
			json.NewEncoder(w).Encode(val)
			return
		}
	}
	msg := fmt.Sprintf("Cant find id %v", ID)
	error.Error = msg
	json.NewEncoder(w).Encode(error)
	return

}

// UpdateTodos godoc
// @Summary Update a todo
// @Description Update status IsComplete a todo
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID of the todos"
// @Success 200 {object} Todos
// @Router /todos/{id} [put]
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func UpdateTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ID, err := strconv.Atoi(params["id"])

	if err != nil {
		msg := fmt.Sprintf("id must be integer")
		error.Error = msg
		json.NewEncoder(w).Encode(error)
		return
	}

	for _, val := range all_todos {
		if val.ID == ID {
			val.IsComplete = "true"
			json.NewEncoder(w).Encode(val)
			return
		}
	}
	msg := fmt.Sprintf("Cant find id %v", ID)
	error.Error = msg
	json.NewEncoder(w).Encode(error)
	return
}

// DeleteTodos godoc
// @Summary Delete a todo
// @Description Delete a todo
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID of the todos"
// @Success 200 {object} Todos
// @Router /todos/{id} [delete]
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func DeleteTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	ID, err := strconv.Atoi(param["id"])
	if err != nil {
		msg := fmt.Sprintf("id must be integer")
		error.Error = msg
		json.NewEncoder(w).Encode(error)
	}
	for i, val := range all_todos {
		if val.ID == ID {
			all_todos = append(all_todos[:i], all_todos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
