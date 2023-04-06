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

	"github.com/gorilla/mux"
)

func Test_CrudClassHandlers(t *testing.T) {
	class := client.Class_API{
		TeachingGrade: 5,
		ProfessorID:   "123456789012345678901234567890",
		Subject:       "math",
		Roster:        make([]string, 2),
	}
	class.Roster[0] = "123456789012345678901234567890"
	class.Roster[1] = "234567890123456789012345678901"

	// create handler first
	form := url.Values{}
	form.Add("teaching_grade", fmt.Sprint(class.TeachingGrade))
	form.Add("professor_id", class.ProfessorID)
	form.Add("subject", class.Subject)
	roster := []string{"123456789012345678901234567890", "234567890123456789012345678901"}
	form.Add("roster", strings.Join(roster, ","))
	reqBody := strings.NewReader(form.Encode())

	reqPOST, err := http.NewRequest("POST", "/classes", reqBody)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}
	reqPOST.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqPOST.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	// create a response recorder
	rr := httptest.NewRecorder()

	// create a router and register the handler
	r := mux.NewRouter()
	r.HandleFunc("/classes", CreateClassHandler).Methods("POST")
	r.HandleFunc("/classes/{class_id}", GetClassHandler).Methods("GET")
	r.HandleFunc("/classes/{class_id}", UpdateClassHandler).Methods("PUT")
	r.HandleFunc("/classes/{class_id}", DeleteClassHandler).Methods("DELETE")

	r.ServeHTTP(rr, reqPOST)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
	//students.GetClass()
	classApi := client.Class_API{}
	err = json.Unmarshal(rr.Body.Bytes(), &classApi)
	if err != nil {
		t.Errorf("error unmarshalling response body: %s", err)
	} else {
		fmt.Println(classApi)
		th.AssertEqual(t, "classID", len(classApi.ClassID), 36)
		th.AssertEqual(t, "teaching grade", classApi.TeachingGrade, 5)
		th.AssertEqual(t, "profID", classApi.ProfessorID, "123456789012345678901234567890")
		th.AssertEqual(t, "subject", classApi.Subject, "math")
		//fmt.Println(classApi.Roster[1])
		if th.AssertEqual(t, "roster len", len(classApi.Roster), 2) {
			th.AssertEqual(t, "roster[0]", classApi.Roster[0], "123456789012345678901234567890")
			th.AssertEqual(t, "roster[1]", classApi.Roster[1], "234567890123456789012345678901")
		}
	}

	// create a new GET request
	requestString := fmt.Sprintf("/classes/%s", classApi.ClassID)
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
	// checking for data loss or incorrect passing of values
	classGet := client.Class_API{}
	err = json.Unmarshal(rr.Body.Bytes(), &classGet)
	if err != nil {
		t.Errorf("error unmarshalling response body: %s", err)
	} else {
		fmt.Println(classGet)
		th.AssertEqual(t, "teaching grade", classGet.TeachingGrade, classApi.TeachingGrade)
		th.AssertEqual(t, "profID", classGet.ProfessorID, classApi.ProfessorID)
		th.AssertEqual(t, "subject", classGet.Subject, classApi.Subject)
		if th.AssertEqual(t, "roster len", len(classGet.Roster), 2) {
			th.AssertEqual(t, "roster[0]", classGet.Roster[0], classApi.Roster[0])
			th.AssertEqual(t, "roster[1]", classGet.Roster[1], classApi.Roster[1])
		} else {
			fmt.Println(classApi.Roster)
		}
	}

	rr.Body.Reset()

	// create a new update request with new values
	form = url.Values{}
	form.Add("professor_id", "123456789012345678901234567895")
	roster = []string{"123456789012345678901234567892", "234567890123456789012345678903"}
	form.Add("class_avg", "3.8")
	form.Add("add_roster", strings.Join(roster, ","))
	reqBody = strings.NewReader(form.Encode())

	// build request
	requestString = fmt.Sprintf("/classes/%s", classApi.ClassID)
	reqUPDATE, err := http.NewRequest("PUT", requestString, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	reqUPDATE.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqUPDATE.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	r.ServeHTTP(rr, reqUPDATE)
	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	classUpdate, err := students.GetClass(classApi.ClassID)
	if err != nil {
		t.Error(err)
	} else {
		th.AssertEqual(t, "professorid", classUpdate.ProfessorID, "123456789012345678901234567895")
		th.AssertEqual(t, "class_avg", classUpdate.ClassAvg, 3.8)
		if th.AssertEqual(t, "roster len", len(classUpdate.Roster), 4) {
			th.AssertEqual(t, "roster[0]", classUpdate.Roster[2], "123456789012345678901234567892")
			th.AssertEqual(t, "roster[1]", classUpdate.Roster[3], "234567890123456789012345678903")
		} else {
			fmt.Println(classApi.Roster)
		}
	}

	rr.Body.Reset()

	requestString = fmt.Sprintf("/classes/%s", classApi.ClassID)
	reqDELETE, err := http.NewRequest("DELETE", requestString, nil)
	if err != nil {
		t.Error(err)
	}
	r.ServeHTTP(rr, reqDELETE)
	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	} else {
		_, err := students.GetClass(classApi.ClassID)
		if err == nil {
			t.Error("Get succeeded after delete")
		}
	}
}
