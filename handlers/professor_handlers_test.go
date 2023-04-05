package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"students/client"
	"students/students"
	th "students/testhelpers"
)

func TestCrudProfessorHandler(t *testing.T) {
	name := "birdperson"

	form := url.Values{}
	form.Add("name", name)

	reqBody := strings.NewReader(form.Encode())

	reqPOST, err := http.NewRequest("POST", "/professors", reqBody)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}
	reqPOST.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqPOST.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/professors", CreateProfessorHandler).Methods("POST")
	r.HandleFunc("/professors/{professor_id}", GetProfessorHandler).Methods("GET")
	r.HandleFunc("/professors/{professor_id}", UpdateProfessorHandler).Methods("PUT")
	r.HandleFunc("/professors/{professor_id}", DeleteProfessorHandler).Methods("DELETE")

	r.ServeHTTP(rr, reqPOST)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	prof := client.Professor_API{}
	err = json.Unmarshal(rr.Body.Bytes(), &prof)
	if err != nil {
		t.Errorf("error unmarshalling response body: %s", err)
	} else {
		th.AssertEqual(t, "profid len", len(prof.ProfessorID), 36)
		th.AssertEqual(t, "name", prof.Name, name)
	}

	rr.Body.Reset()

	reqGET, err := http.NewRequest("GET", "/professors/"+prof.ProfessorID, nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	r.ServeHTTP(rr, reqGET)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	profGet := client.Professor_API{}
	err = json.Unmarshal(rr.Body.Bytes(), &profGet)
	if err != nil {
		t.Errorf("error unmarshalling response body: %s", err)
	} else {
		th.AssertEqual(t, "profid len", len(profGet.ProfessorID), 36)
		th.AssertEqual(t, "profid", profGet.ProfessorID, prof.ProfessorID)
		th.AssertEqual(t, "name", profGet.Name, name)
	}

	rr.Body.Reset()

	studentGpa := 3.5
	form = url.Values{}
	form.Add("student_avg", fmt.Sprintf("%f", studentGpa))
	reqBody = strings.NewReader(form.Encode())

	reqUPDATE, err := http.NewRequest("PUT", "/professors/"+prof.ProfessorID, reqBody)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}
	reqUPDATE.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqUPDATE.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

	r.ServeHTTP(rr, reqUPDATE)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
	profUpdate, err := students.GetProfessor(prof.ProfessorID)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(profUpdate)
		th.AssertEqual(t, "profid len", len(profUpdate.ProfessorID), 36)
		th.AssertEqual(t, "profid", profUpdate.ProfessorID, prof.ProfessorID)
		th.AssertEqual(t, "name", profUpdate.Name, prof.Name)
		th.AssertEqual(t, "student avg", profUpdate.StudentAvg, studentGpa)
	}

	rr.Body.Reset()

	reqDELETE, err := http.NewRequest("DELETE", "/professors/"+prof.ProfessorID, nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	r.ServeHTTP(rr, reqDELETE)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	} else {
		_, err := students.GetProfessor(prof.ProfessorID)
		if err == nil {
			t.Error("Get succeeded after delete")
		}
	}

}

// func TestGetProfessorHandler(t *testing.T) {
// 	name := "John Doe"
// 	professor, err := students.CreateProfessor(name)
// 	assert.NoError(t, err)

// 	req, err := http.NewRequest("GET", "/professors/"+professor.ProfessorID.String(), nil)
// 	if err != nil {
// 		t.Fatalf("error creating request: %s", err)
// 	}

// 	rr := httptest.NewRecorder()

// 	r := mux.NewRouter()
// 	r.HandleFunc("/professors/{professor_id}", GetProfessorHandler).Methods("GET")

// 	r.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	var result students.Professor
// 	err = json.Unmarshal(rr.Body.Bytes(), &result)
// 	assert.NoError(t, err)
// 	assert.Equal(t, professor.ProfessorID, result.ProfessorID)
// 	assert.Equal(t, professor.Name, result.Name)
// }

// func TestUpdateProfessorHandler(t *testing.T) {
// 	name := "John Doe"
// 	professor, err := students.CreateProfessor(name)
// 	assert.NoError(t, err)

// 	form := url.Values{}
// 	form.Add("student_avg", "3.5")

// 	reqBody := strings.NewReader(form.Encode())

// 	reqPOST, err := http.NewRequest("PUT", "/professors/"+professor.ProfessorID, reqBody)
// 	if err != nil {
// 		t.Fatalf("error creating request: %s", err)
// 	}
// 	reqPOST.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	reqPOST.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))

// 	rr := httptest.NewRecorder()

// 	r := mux.NewRouter()
// 	r.HandleFunc("/professors/{professor_id}", UpdateProfessorHandler).Methods("PUT")

// 	r.ServeHTTP(rr, reqPOST)

// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	result, err := students.GetProfessor(professor.ProfessorID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 3.5, result.StudentAvg)
// }

// func TestDeleteProfessorHandler(t *testing.T) {
// 	name := "John Doe"
// 	professor, err := students.CreateProfessor(name)
// 	assert.NoError(t, err)

// 	//req, err := http.NewRequest("DELETE", "/professors/"+professor.Professor
// }

/*
func Test_CrudProfHandlers(t *testing.T) {
	prof := client.Professor_API{
		ProfessorID:   "123456789012345678901234567890",
		Name:       "bob",
		StudentAvg: 3.51,
		ClassList:        make([]string, 2),
	}
	prof.ClassList[0] = "123456789012345678901234567890"
	prof.ClassList[1] = "234567890123456789012345678901"

	// create handler first
	form := url.Values{}
	form.Add("professor_id", prof.ProfessorID)
	form.Add("name", prof.Name)
	classList := []string{"123456789012345678901234567890", "234567890123456789012345678901"}
	form.Add("classlist", strings.Join(classList, ","))
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
		}
	}

	rr.Body.Reset()

	// create a new update request with new values
	form = url.Values{}
	form.Add("professor_id", "123456789012345678901234567895")
	roster = []string{"123456789012345678901234567892", "234567890123456789012345678903"}
	form.Add("class_avg", "3.8")
	form.Add("roster", strings.Join(roster, ","))
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
		if th.AssertEqual(t, "roster len", len(classUpdate.Roster), 2) {
			th.AssertEqual(t, "roster[0]", classUpdate.Roster[0], "123456789012345678901234567892")
			th.AssertEqual(t, "roster[1]", classUpdate.Roster[1], "234567890123456789012345678903")
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
}*/
