package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateStudentHandler(t *testing.T) {
	// Create a request body
	form := url.Values{}
	form.Add("name", "John Doe")
	form.Add("current_year", "3")
	form.Add("graduation_year", "2023")
	form.Add("gpa", "3.5")
	form.Add("age", "21")
	form.Add("dob", "2000-01-01")
	form.Add("enrolled", "true")
	reqBody := strings.NewReader(form.Encode())

	// Create a request with the body and content-type header
	req, err := http.NewRequest("POST", "/students", reqBody)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a router with the CreateStudentHandler route
	router := mux.NewRouter()
	router.HandleFunc("/students", CreateStudentHandler).Methods("POST")

	// Send the request through the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check the response body
	var respBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %s", err)
	}
	if _, ok := respBody["id"]; !ok {
		t.Errorf("response body does not contain id field")
	}
}
func TestGetStudentHandler(t *testing.T) {
	// Create a request with the student id parameter
	req, err := http.NewRequest("GET", "/students/1", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a router with the GetStudentHandler route
	router := mux.NewRouter()
	router.HandleFunc("/students/{id}", GetStudentHandler).Methods("GET")

	// Set the id parameter in the request context
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Send the request through the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check the response body
	var respBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %s", err)
	}
	if respBody["id"] != 1 {
		t.Errorf("response body has wrong id field: got %v want %v", respBody["id"], 1)
	}
}
