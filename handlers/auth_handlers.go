package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"students/auth"
)

type User_Login struct {
	UserName string
	Password string
}

// AUTH HANDLERS
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Read the user's credentials from the request body
	// username password
	var user User_Login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"errorType": "BadRequest", "message": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	if len(user.UserName) == 0 {
		http.Error(w, `{"errorType": "BadRequest", "message": "Username is required"}`, http.StatusBadRequest)
		fmt.Println("empty request")
		return
	}
	if len(user.Password) == 0 {
		http.Error(w, `{"errorType": "BadRequest", "message": "Password is required"}`, http.StatusBadRequest)
		fmt.Println("empty request")
		return
	}
	// Verify the user's credentials
	authenticated, userInfo, err := auth.AuthenticateUser(user.UserName, user.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, `{"errorType": "Unauthorized", "message": "Incorrect Username or Password"}`, http.StatusUnauthorized)
			return
		} else if strings.Contains(err.Error(), "hash") {
			http.Error(w, `{"errorType": "IncorrectPassword", "message": "Invalid username or password"}`, http.StatusUnauthorized)
		} else {
			http.Error(w, `{"errorType": "General", "message": "Unexpected Error"}`, http.StatusInternalServerError)
		}
		fmt.Println(err)
		return
	}
	// redundant can probs delete check tho
	if !authenticated {
		http.Error(w, `{"errorType": "Unauthorized", "message": "Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	// Generate a new JWT token
	token, err := auth.GenerateToken(userInfo.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	fmt.Println(userInfo)
	// Return the JWT token to the client
	json.NewEncoder(w).Encode(map[string]string{
		"token":    token,
		"username": userInfo.UserName,
		"ownerID":  userInfo.OwnerID,
		"email":    maskEmail(userInfo.Email),
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// going to use local storage for funzies but ill leave this here for now
	// cookie := &http.Cookie{
	// 	Name:     "jwt_token",
	// 	Value:    "",
	// 	Expires:  time.Unix(0, 0),
	// 	HttpOnly: true,
	// }
	// http.SetCookie(w, cookie)

	// Redirect the user to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func maskEmail(email string) string {
	atIndex := strings.LastIndex(email, "@")
	if atIndex == -1 {
		return email
	}
	domainStart := atIndex + 1
	domainIndex := domainStart + strings.Index(email[domainStart:], ".")
	if domainIndex == -1 {
		return email
	}
	maskedUser := email[:3] + strings.Repeat("*", atIndex-3)
	maskedDomain := email[domainStart:domainStart+1] + strings.Repeat("*", domainIndex-domainStart-1) + email[domainIndex:]
	return maskedUser + "@" + maskedDomain
}
