package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"students/client"
	"students/students"

	"github.com/gorilla/mux"
)

// CRUD HANDLERS
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	options := students.CreateNewUserOptions{
		UserName:       r.FormValue("user_name"),
		Email:          r.FormValue("email"),
		HashedPassword: r.FormValue("hashed_password"),
	}

	user, err := options.CreateNewUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error creating user")
		return
	}

	ret, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error marshalling user")
		return
	}
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
		_, err := students.GetUser(ownerID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "No user with that ID")
			return
		}
		opts.OwnerID = ownerID
	}

	opts.UserName = r.FormValue("user_name")
	opts.Email = r.FormValue("email")
	opts.HashedPassword = r.FormValue("hashed_password")
	opts.AddSchoolList = strings.Split(r.FormValue("add_school_list"), ",")
	opts.RemoveSchoolList = strings.Split(r.FormValue("remove_school_list"), ",")

	err := opts.UpdateUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerID := vars["owner_id"]
	err := students.DeleteUser(ownerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error deleting user")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AUTH HANDLERS
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Read the user's credentials from the request body
	// username password
	var user client.User_API
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify the user's credentials
	authenticated, userInfo, err := students.AuthenticateUser(user.UserName, user.HashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !authenticated {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate a new JWT token
	token, err := students.GenerateToken(userInfo.OwnerID, userInfo.UserName, userInfo.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the JWT token to the client
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear authentication token or cookie
	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	// Redirect the user to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
