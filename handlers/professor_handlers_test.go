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
	th.TestingInit()
	defer th.TestingEnvDif()
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
