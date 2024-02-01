package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"students/client"
	"students/students"
	th "students/testhelpers"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// TODO: defer deletes on all the created information

func TestCreateStudentHandler(t *testing.T) {
	// Create a request body
	th.TestingInit()
	defer th.TestingEnvDif()
	form := url.Values{}
	form.Add("name", "John Doe")
	form.Add("current_year", "3")
	form.Add("graduation_year", "2023")
	form.Add("avg_gpa", "3.5")
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

	studentapi := client.Student_API{}
	err = json.Unmarshal(rr.Body.Bytes(), &studentapi)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %s", err)
	}
	fmt.Println(studentapi)
	th.AssertEqual(t, "student name", studentapi.Name, "John Doe")
	th.AssertEqual(t, "current year", studentapi.CurrentYear, 3)
	th.AssertEqual(t, "graduation year", studentapi.GraduationYear, 2023)
	th.AssertEqual(t, "gpa", studentapi.AvgGPA, 3.5)
	th.AssertEqual(t, "api test", studentapi.Age, 21)
	//th.AssertEqual(t, "dob", studentapi.Dob, "2000-01-01")
	th.AssertEqual(t, "enrolled", studentapi.Enrolled, true)

}
func TestGetStudentHandler(t *testing.T) {
	th.TestingInit()
	defer th.TestingEnvDif()
	// Create a request with the student id parameter
	dob, err := time.Parse("2006-01-02", "2013-01-07")
	if err != nil {
		t.Error(err)
	}

	student := client.Student_API{
		Name:           "vesemir",
		CurrentYear:    5,
		GraduationYear: 2031,
		AvgGPA:         3.8,
		Age:            11,
		Dob:            dob,
		Enrolled:       true,
	}
	newStudent, err := students.CreateNewStudent(student.Name, student.CurrentYear, student.GraduationYear, student.AvgGPA, student.Age, student.Dob, student.Enrolled)
	if err != nil {
		t.Error(err)
	}
	requestString := fmt.Sprintf("/students/%s", newStudent.StudentID)
	req, err := http.NewRequest("GET", requestString, nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a router with the GetStudentHandler route
	router := mux.NewRouter()
	router.HandleFunc("/students/{student_id}", GetStudentHandler).Methods("GET")

	// Set the id parameter in the request context
	req = mux.SetURLVars(req, map[string]string{"student_id": newStudent.StudentID})

	// Send the request through the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
	var respBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %s", err)
	}

	studentapi := client.Student_API{}
	err = json.Unmarshal(rr.Body.Bytes(), &studentapi)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %s", err)
	}
	th.AssertEqual(t, "student name", studentapi.Name, student.Name)
	th.AssertEqual(t, "current year", studentapi.CurrentYear, student.CurrentYear)
	th.AssertEqual(t, "graduation year", studentapi.GraduationYear, student.GraduationYear)
	th.AssertEqual(t, "gpa", studentapi.AvgGPA, student.AvgGPA)
	th.AssertEqual(t, "api test", studentapi.Age, student.Age)
	th.AssertEqual(t, "enrolled", studentapi.Enrolled, true)
}
func TestUpdateStudentHandler(t *testing.T) {
	th.TestingInit()
	defer th.TestingEnvDif()
	// Create a request with the student id parameter
	dob, err := time.Parse("2006-01-02", "2013-01-07")
	if err != nil {
		t.Error(err)
	}
	form := url.Values{}
	//form.Add("name", "mittens")
	form.Add("current_year", "6")
	form.Add("graduation_year", "2032")
	form.Add("avg_gpa", "3.8")
	form.Add("age", "12")
	form.Add("dob", "2013-01-07")
	form.Add("enrolled", "true")
	reqBody := strings.NewReader(form.Encode())
	student := client.Student_API{
		Name:           "vesemir",
		CurrentYear:    5,
		GraduationYear: 2031,
		AvgGPA:         3.8,
		Age:            11,
		Dob:            dob,
		Enrolled:       true,
	}
	newStudent, err := students.CreateNewStudent(student.Name, student.CurrentYear, student.GraduationYear, student.AvgGPA, student.Age, student.Dob, student.Enrolled)
	if err != nil {
		t.Error(err)
	}
	requestString := fmt.Sprintf("/students/%s", newStudent.StudentID)
	req, err := http.NewRequest("PUT", requestString, reqBody)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a router with the GetStudentHandler route
	router := mux.NewRouter()
	router.HandleFunc("/students/{student_id}", UpdateStudentHandler).Methods("PUT")

	// Set the id parameter in the request context
	req = mux.SetURLVars(req, map[string]string{"student_id": newStudent.StudentID})

	// Send the request through the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	retStudent, err := students.GetStudent(newStudent.StudentID)
	if err != nil {
		t.Error(err)
	}
	//th.AssertEqual(t, "student name", retStudent.Name, "mittens")
	th.AssertEqual(t, "current year", retStudent.CurrentYear, 6)
	th.AssertEqual(t, "graduation year", retStudent.GraduationYear, 2032)
	th.AssertEqual(t, "gpa", retStudent.AvgGPA, 3.8)
	th.AssertEqual(t, "api test", retStudent.Age, 12)
	//th.AssertEqual(t, "dob", studentapi.Dob, "2000-01-01")
	th.AssertEqual(t, "enrolled", retStudent.Enrolled, true)

}
func TestDeleteStudentHandler(t *testing.T) {
	// Create a new student
	th.TestingInit()
	defer th.TestingEnvDif()
	dob, _ := time.Parse("2006-01-02", "2013-01-07")
	newStudent, err := students.CreateNewStudent("John Doe", 6, 2031, 3.7, 12, dob, true)
	if err != nil {
		t.Error(err)
	}

	// Create a DELETE request with the student id parameter
	requestString := fmt.Sprintf("/students/%s", newStudent.StudentID)
	req, err := http.NewRequest("DELETE", requestString, nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a router with the DeleteStudentHandler route
	router := mux.NewRouter()
	router.HandleFunc("/students/{student_id}", DeleteStudentHandler).Methods("DELETE")

	// Set the id parameter in the request context
	req = mux.SetURLVars(req, map[string]string{"student_id": newStudent.StudentID})

	// Send the request through the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check that the student was deleted
	_, err = students.GetStudent(newStudent.StudentID)
	if err == nil {
		t.Errorf("handler did not delete student: %s", newStudent.StudentID)
	}
}
