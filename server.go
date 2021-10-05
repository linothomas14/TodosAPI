package main

import (
	"encoding/json"
	_ "final-project-1/docs"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todos struct {
	ID          int    `json:"id" gorm:"primaryKey" example:"1"`
	Title       string `json:"title" example:"Reading book"`
	Description string `json:"description" example:"Reading book at 9 A.M"`
	IsComplete  string `json:"is_complete" gorm:"default:false" example:"false"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error"`
}

var db *gorm.DB
var prevOrderID = 0
var all_todos []Todos

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

func initDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/todosapp?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	log.Println("SUCCESS ACCESS DB")
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Todos{})

}

func main() {
	initDB()
	router := mux.NewRouter()
	router.HandleFunc("/todos", GetTodos).Methods("GET")
	router.HandleFunc("/todos", CreateTodos).Methods("POST")
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
// @Success 200 {array} Todos
// @Router /todos [get]
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todos

	w.Header().Set("Content-Type", "application/json")

	db.Find(&todos)

	json.NewEncoder(w).Encode(todos)

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
	print(todos.ID)
	if err != nil {
		fmt.Println("ERROR")
		log.Fatal(err)
	}

	db.Create(&todos)
	json.NewEncoder(w).Encode(todos)
}
