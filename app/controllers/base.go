package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres

	"github.com/hudaprs/RESTfulAPIScotchIO/app/api/middlewares"
	"github.com/hudaprs/RESTfulAPIScotchIO/app/api/responses"
	"github.com/hudaprs/RESTfulAPIScotchIO/app/models"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (app *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error
	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	app.DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Printf("\n Cannot connect to database %s", DbName)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the database %s", DbName)
	}

	app.DB.Debug().AutoMigrate(&models.User{}) // database migration

	app.Router = mux.NewRouter().StrictSlash(true)
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.Use(middlewares.SetContentTypeMiddleware) // setting content-type to json

	app.Router.HandleFunc("/", home).Methods("GET")
	app.Router.HandleFunc("/api/register", app.UserSignUp).Methods("POST")
	app.Router.HandleFunc("/api/login", app.Login).Methods("POST")
}

func (app *App) RunServer() {
	log.Printf("\n Server starting at port 5000")
	log.Fatal(http.ListenAndServe(":5000", app.Router))
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To RESTfulAPI with GO")
}
