package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"students/client"
	"students/logger"
	"students/students"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// CRUD HANDLERS
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user client.User_API

	fmt.Println("got here")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	// TODO: I need to do a check here for usernames that already exist? or I could just not and check the sql error?
	// checking the sql error would save a request. Will add this feature once this is functional.

	fmt.Println()
	fmt.Printf("user: %v", user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	options := students.CreateNewUserOptions{
		UserName:       user.UserName,
		Email:          user.Email,
		HashedPassword: string(hashedPassword),
	}

	newUser, err := options.CreateNewUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error creating user")
		fmt.Println(err)
		return
	}

	ret, err := json.Marshal(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error marshalling user")
		fmt.Println(err)
		return
	}
	fmt.Printf("created user: %v", user)
	fmt.Println()
	w.Write(ret)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerID := vars["owner_id"]
	user, err := students.GetUser(ownerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error retrieving user")
		return
	}
	if len(user.OwnerID) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "User not found")
		return
	}
	fmt.Printf("returned user %v", user)
	ret, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error marshalling user")
		return
	}
	w.Write(ret)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	opts := students.UpdateUserOptions{}
	vars := mux.Vars(r)

	ownerID := vars["owner_id"]
	if len(ownerID) <= 30 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid owner ID")
		return
	} else {
		// Check if user exists
		userID := r.Context().Value("user_id")
		if userID == nil {
			// Handle error: no user ID in context
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Please log in")
			return
		}
		if ownerID != userID {
			logger.Log.WithFields(logrus.Fields{
				"owner_id": ownerID,
				"user_id":  userID,
			}).Warn("attempted update of non owned user")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "non authorized get attempt of unowned user")
			return
		}
		_, err := students.GetUser(ownerID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "No user with that ID")
			return
		}
		opts.OwnerID = ownerID
	}
	fmt.Println("Update user")
	opts.UserName = r.FormValue("user_name")
	opts.Email = r.FormValue("email")
	newPassword := r.FormValue("password")
	if len(newPassword) != 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		opts.HashedPassword = string(hashedPassword)

	}
	// opts.AddSchoolList = strings.Split(r.FormValue("add_school_list"), ",")
	// opts.RemoveSchoolList = strings.Split(r.FormValue("remove_school_list"), ",")
	fmt.Println(opts)
	err := opts.UpdateUser()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(err.Error(), "duplicate key value") {
			http.Error(w, `{"errorType": "DuplicateKey", "message": "Duplicate Key Value"}`, http.StatusConflict)
		} else {
			http.Error(w, `{"errorType": "General", "message": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerID := vars["owner_id"]
	userID := r.Context().Value("user_id")
	if userID == nil {
		// Handle error: no user ID in context
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Please log in")
		return
	}
	if ownerID != userID {
		logger.Log.WithFields(logrus.Fields{
			"owner_id": ownerID,
			"user_id":  userID,
		}).Warn("attempted delete of non owned user")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "non authorized delete attempt of unowned user")
		return
	}
	err := students.DeleteUser(ownerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error deleting user")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// This function normally wouldn't be included here but I dont have enough functions
// to reasonably create a status file. If that changes the functions below this point will be moved.

func BackendStatus(w http.ResponseWriter, r *http.Request) {

	// If I want later on I can write a better test to see if it stuff is running
	isRunning := true

	response := struct {
		IsRunning bool `json:"isRunning"`
	}{
		IsRunning: isRunning,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
