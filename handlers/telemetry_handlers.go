package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"students/students"
	"students/students/telemetry"

	"github.com/gorilla/mux"
)

func CreateNewSchoolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner_id := vars["owner_id"]
	numPerGrade, err := strconv.Atoi(r.FormValue("num_per_grade"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	} else if numPerGrade > 100 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	name := r.FormValue("name")
	school, err := telemetry.NewSchool(numPerGrade, owner_id, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
	ret, err := json.Marshal(school)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling professor")
		return
	}
	w.Write(ret)
}

func GetGradeAvgForSchoolHandler(w http.ResponseWriter, r *http.Request) {
	// Get the schoolID from the URL parameter
	vars := mux.Vars(r)
	schoolID := vars["school_id"]
	if len(schoolID) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "no given id")
		return
	}
	// check if exist
	_, err := students.GetSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "no school with that id")
		return
	}
	// attempt update  TODO: fix this function to auto check if the values are 0.0 and run a class avg update if so
	gradeAvgList, err := telemetry.GetGradeAvgForSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting grade averages: %v", err)
		return
	}
	// means update
	if gradeAvgList[0].AvgGPA == 0.0 {
		classList, err := students.GetClassesForSchool(schoolID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "error on class get")
			return
		}
		err = telemetry.UpdateClassAvgs(classList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "error on class avg update")
			return
		}
		gradeAvgList, err = telemetry.GetGradeAvgForSchool(schoolID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error getting grade averages: %v", err)
			return
		}
	}
	ret, err := json.Marshal(gradeAvgList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling professor")
		return
	}
	w.Write(ret)
}

func GetBestProfessorsHandler(w http.ResponseWriter, r *http.Request) {
	professorIDs := r.URL.Query().Get("professor_ids")
	if len(professorIDs) == 0 {
		http.Error(w, "Missing professor_ids query parameter", http.StatusBadRequest)
		return
	}
	professorList := strings.Split(professorIDs, ",")
	if len(professorList) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error  professorList")
		return
	}
	bestProfessors, err := telemetry.GetBestProfessors(professorList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("logging error at professorget %v", err)
		return
	}
	if bestProfessors[0].StudentAvg == 0 {
		fmt.Println("Found values needing update for professors. Running updates!")
		err = telemetry.UpdateProfessorsStudentAvgs(professorList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Unexpected error during value updates")
			fmt.Printf("logging error at professorupdate %v", err)
			return
		}
		bestProfessors, err = telemetry.GetBestProfessors(professorList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	ret, err := json.Marshal(bestProfessors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling professor")
		return
	}
	w.Write(ret)
}
