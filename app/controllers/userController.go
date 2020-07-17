package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hudaprs/RESTfulAPIScotchIO/app/api/responses"
	"github.com/hudaprs/RESTfulAPIScotchIO/app/models"
	"github.com/hudaprs/RESTfulAPIScotchIO/app/utils"
)

// UserSignUp Controller for creating new user
func (app *App) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Registered Successfully"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, _ := user.GetUser(app.DB)
	if usr != nil {
		resp["status"] = "Failed"
		resp["message"] = "User already registered, please login"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	user.Prepare() // here strip the next of white spaces

	err = user.Validate("") //  default were all fields(email, lastname, firstname, password, profileimage) area validated
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	userCreated, err := user.SaveUser(app.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["user"] = userCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

// Login signs in users
func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "logged in"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body) // read user input from request
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user.Prepare() //  here strip the text of white space

	err = user.Validate("login") // fields(email, password) are validated
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	usr, err := user.GetUser(app.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if usr == nil { // User is not registered
		resp["status"] = "failed"
		resp["message"] = "Login failed, please signup"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(user.Password, usr.Password)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "Login failed, please try again"
		responses.JSON(w, http.StatusForbidden, resp)
		return
	}

	token, err := utils.EncodeAuthToken(usr.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)
	return
}
