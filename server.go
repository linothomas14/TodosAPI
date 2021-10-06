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

type Todo struct {
	ID          string `json:"id" example:"1"`
	Title       string `json:"title" example:"Reading book"`
	Description string `json:"description" example:"Reading book at 9 A.M"`
	IsComplete  bool `json:"is_complete" example:"false"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error"`
}

var prevOrderID = 0
var todos []Todo

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
	router.HandleFunc("/todos", GetTodos).Methods("GET")
	router.HandleFunc("/todos", CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", GetTodo).Methods("GET")
	router.HandleFunc("/todos/{id}", UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")
	log.Println("Start server at localhost:8080")
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	http.ListenAndServe(":8080", router)
}

// GetTodos godoc
// @Summary Get all todos
// @Description Get details of all todos
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} Todo
// @Router /todos [get]
func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)

}

// CreateTodo godoc
// @Summary Create a todo
// @Description Create a new todo with the input payload
// @Tags todos
// @Accept json
// @Produce json
// @Param todos body Todo true "Create todo"
// @Success 201 {object} Todo
// @Router /todos [post]
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	w.Header().Set("Content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		fmt.Println("ERROR")
		log.Fatal(err)
	}
	prevOrderID++
	todo.ID = strconv.Itoa(prevOrderID)
	todos = append(todos, todo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// GetTodo godoc
// @Summary Get todo
// @Description Get details of todo
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} Todo
// @Param id path int true "Todo id"
// @Router /todos/{id} [get]
// @Failure 404 {object} ErrorResponse
func GetTodo(w http.ResponseWriter, r *http.Request){
	var todo Todo
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	indexTodo := GetIndexTodos(todos, params["id"])
	if indexTodo == -1 {
		var errorResponse ErrorResponse
		errorResponse.Error = "data not found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return 
	}
	todo = todos[indexTodo]
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo godoc
// @Summary Update todo
// @Description Update todo
// @Tags todos
// @Accept json
// @Produce json
// @Param todos body Todo true "Update todo"
// @Success 200 
// @Param id path int true "Todo id"
// @Router /todos/{id} [put]
// @Failure 404 {object} ErrorResponse
func UpdateTodo(w http.ResponseWriter, r *http.Request){
	var todo Todo
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&todo)
	indexTodo := GetIndexTodos(todos, params["id"])
	if indexTodo == -1 {
		var errorResponse ErrorResponse
		errorResponse.Error = "data not found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return 
	}
	todo.ID = params["id"]
	todos[indexTodo] = todo
	w.WriteHeader(http.StatusOK)
}

// DeleteTodo godoc
// @Summary Delete todo
// @Description Delete todo
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 
// @Param id path int true "Todo id"
// @Router /todos/{id} [delete]
// @Failure 404 {object} ErrorResponse
func DeleteTodo(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	indexTodo := GetIndexTodos(todos, params["id"])
	if indexTodo == -1 {
		var errorResponse ErrorResponse
		errorResponse.Error = "data not found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return 
	}
	todos = append(todos[:indexTodo], todos[indexTodo+1:]...)
	w.WriteHeader(http.StatusOK)
}

func GetIndexTodos(todos []Todo, idTodo string) int {
	for i, todo := range todos {
		if todo.ID == idTodo {
			return i
		}
	}
	return -1
}