package main

import (
	"encoding/json"
	_ "todosAPI/docs"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todosAPI/helpers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type TodoRequest struct {
	Title       string `json:"title" example:"Reading book"`
	Description string `json:"description" example:"Reading book at 9 A.M"`
	IsComplete  bool `json:"is_complete" example:"false"`
}

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsComplete  bool `json:"is_complete"`
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
// @Success 200 {object} helpers.Response{data=[]Todo}
// @Router /todos [get]
func GetTodos(w http.ResponseWriter, r *http.Request) {
	helpers.HttpResponse(w, http.StatusOK, helpers.Response{"success", todos})
}

// CreateTodo godoc
// @Summary Create a todo
// @Description Create a new todo with the input payload
// @Tags todos
// @Accept json
// @Produce json
// @Param todos body TodoRequest true "Create todo"
// @Success 201 {object} helpers.Response{data=Todo}
// @Router /todos [post]
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		fmt.Println("ERROR ", err.Error())
		helpers.HttpResponse(w, http.StatusBadRequest, helpers.Response{err.Error(), nil})
		return
	}
	prevOrderID++
	todo.ID = strconv.Itoa(prevOrderID)
	todos = append(todos, todo)
	helpers.HttpResponse(w, http.StatusCreated, helpers.Response{"success", todo})
}

// GetTodo godoc
// @Summary Get todo
// @Description Get details of todo
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response{data=Todo}
// @Param id path int true "Todo id"
// @Router /todos/{id} [get]
// @Failure 404 {object} helpers.Response
func GetTodo(w http.ResponseWriter, r *http.Request){
	var todo Todo
	params := mux.Vars(r)
	indexTodo := GetIndexTodos(todos, params["id"])
	if indexTodo == -1 {
		helpers.HttpResponse(w, http.StatusNotFound, helpers.Response{"data not found", nil})
		return 
	}
	todo = todos[indexTodo]
	helpers.HttpResponse(w, http.StatusOK, helpers.Response{"success", todo})
}

// UpdateTodo godoc
// @Summary Update todo
// @Description Update todo
// @Tags todos
// @Accept json
// @Produce json
// @Param todos body TodoRequest true "Update todo"
// @Success 200 {object} helpers.Response
// @Param id path int true "Todo id"
// @Router /todos/{id} [put]
// @Failure 404 {object} helpers.Response
func UpdateTodo(w http.ResponseWriter, r *http.Request){
	var todo Todo
	params := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		fmt.Println("ERROR ", err.Error())
		helpers.HttpResponse(w, http.StatusBadRequest, helpers.Response{err.Error(), nil})
		return
	}
	indexTodo := GetIndexTodos(todos, params["id"])
	if indexTodo == -1 {
		helpers.HttpResponse(w, http.StatusNotFound, helpers.Response{"data not found", nil})
		return 
	}
	todo.ID = params["id"]
	todos[indexTodo] = todo
	helpers.HttpResponse(w, http.StatusOK, helpers.Response{"success", nil})
}

// DeleteTodo godoc
// @Summary Delete todo
// @Description Delete todo
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Param id path int true "Todo id"
// @Router /todos/{id} [delete]
// @Failure 404 {object} helpers.Response
func DeleteTodo(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	indexTodo := GetIndexTodos(todos, params["id"])
	if indexTodo == -1 {
		helpers.HttpResponse(w, http.StatusNotFound, helpers.Response{"data not found", nil})
		return 
	}
	todos = append(todos[:indexTodo], todos[indexTodo+1:]...)
	helpers.HttpResponse(w, http.StatusOK, helpers.Response{"success", nil})
}

func GetIndexTodos(todos []Todo, idTodo string) int {
	for i, todo := range todos {
		if todo.ID == idTodo {
			return i
		}
	}
	return -1
}