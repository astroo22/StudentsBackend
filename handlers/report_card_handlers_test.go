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

func Test_CrudReportCardHandlers(t *testing.T) {
	th.TestingInit()
	defer th.TestingEnvDif()
	dob, err := time.Parse("2006-01-02", "2020-02-07")
	if err != nil {
		t.Error(err)
	}

	student := client.Student_API{
		Name:           "vesemir",
		CurrentYear:    5,
		GraduationYear: 2031,
		AvgGPA:         3.1,
		Age:            3,
		Dob:            dob,
		Enrolled:       true,
	}
	newStudent, err := students.CreateNewStudent(student.Name, student.CurrentYear, student.GraduationYear, student.AvgGPA, student.Age, student.Dob, student.Enrolled)
	if err != nil {
		t.Error(err)
	}

	// response recorder
	rr := httptest.NewRecorder()
	// router
	r := mux.NewRouter()

	r.HandleFunc("/reportcard", CreateReportCardHandler).Methods("POST")
	r.HandleFunc("/reportcard/{student_id}", GetReportCardHandler).Methods("GET")
	r.HandleFunc("/reportcard/{student_id}", UpdateReportCardHandler).Methods("PUT")
	r.HandleFunc("/reportcard/{student_id}", DeleteReportCardHandler).Methods("DELETE")

	form := url.Values{}
	form.Add("student_id", newStudent.StudentID)
	reqBody := strings.NewReader(form.Encode())
	reqPOST, err := http.NewRequest("POST", "/reportcard", reqBody)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	reqPOST.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqPOST.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	r.ServeHTTP(rr, reqPOST)
	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	reportCard := client.ReportCard_API{}
	err = json.Unmarshal(rr.Body.Bytes(), &reportCard)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %s", err)
	} else {
		th.AssertNotEqual(t, "math", reportCard.Math, 0.0)
		th.AssertNotEqual(t, "english", reportCard.English, 0.0)
		th.AssertNotEqual(t, "science", reportCard.Science, 0.0)
		th.AssertNotEqual(t, "physicaled", reportCard.PhysicalED, 0.0)
		th.AssertNotEqual(t, "lunch", reportCard.Lunch, 0.0)
	}

	// create a new GET request
	requestString := fmt.Sprintf("/reportcard/%s", newStudent.StudentID)
	reqGET, err := http.NewRequest("GET", requestString, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr.Body.Reset()
	// serve the request and record the response
	r.ServeHTTP(rr, reqGET)
	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
	reportCardGET := client.ReportCard_API{}
	err = json.Unmarshal(rr.Body.Bytes(), &reportCardGET)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %s", err)
	} else {
		th.AssertEqual(t, "math", reportCardGET.Math, reportCard.Math)
		th.AssertEqual(t, "english", reportCardGET.English, reportCard.English)
		th.AssertEqual(t, "science", reportCardGET.Science, reportCard.Science)
		th.AssertEqual(t, "physicaled", reportCardGET.PhysicalED, reportCard.PhysicalED)
		th.AssertEqual(t, "lunch", reportCardGET.Lunch, reportCard.Lunch)
	}

	//update
	updateCard := client.ReportCard_API{
		Math:       4.00,
		Science:    0.01,
		English:    2.81,
		PhysicalED: .5,
		Lunch:      1,
	}

	rr.Body.Reset()

	form = url.Values{}
	form.Add("math", fmt.Sprintf("%f", updateCard.Math))
	form.Add("english", fmt.Sprintf("%f", updateCard.English))
	form.Add("science", fmt.Sprintf("%f", updateCard.Science))
	form.Add("physicaled", fmt.Sprintf("%f", updateCard.PhysicalED))
	form.Add("lunch", fmt.Sprintf("%f", updateCard.Lunch))
	reqBody = strings.NewReader(form.Encode())

	requestString = fmt.Sprintf("/reportcard/%s", newStudent.StudentID)
	reqUPDATE, err := http.NewRequest("PUT", requestString, reqBody)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}
	reqUPDATE.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqUPDATE.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	r.ServeHTTP(rr, reqUPDATE)
	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	reportCardUpdate, err := students.GetReportCard(newStudent.StudentID)
	if err != nil {
		t.Errorf("error getting update: %d", err)
	} else {
		th.AssertEqual(t, "math", reportCardUpdate.Math, updateCard.Math)
		th.AssertEqual(t, "english", reportCardUpdate.English, updateCard.English)
		th.AssertEqual(t, "science", reportCardUpdate.Science, updateCard.Science)
		th.AssertEqual(t, "physicaled", reportCardUpdate.PhysicalED, updateCard.PhysicalED)
		th.AssertEqual(t, "lunch", reportCardUpdate.Lunch, updateCard.Lunch)
	}

	rr.Body.Reset()

	reqDELETE, err := http.NewRequest("DELETE", requestString, nil)
	if err != nil {
		t.Error(err)
	}
	r.ServeHTTP(rr, reqDELETE)
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	} else {
		_, err := students.GetReportCard(newStudent.StudentID)
		if err == nil {
			t.Error("Get succeeded after delete")
		}
	}
}
