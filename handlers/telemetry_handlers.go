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
	fmt.Println("REMINDER YOU NEED TO FIX ALL THE DATA IN THE DB FOR THIS TO WORK SO THAT IT UPDATES ON CREATE IN FUTURE")
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
	fmt.Println("attempting getBestProf")
	// professors := r.PostFormValue("professor_ids")

	// fmt.Printf("retriving best professors m1: %v", professorIDs)
	fmt.Println("")
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
	fmt.Println(len(bestProfessors))
	if bestProfessors[0].StudentAvg == 0 {
		fmt.Println("Found values needed update for professors")
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
		fmt.Println(len(bestProfessors))
	}
	ret, err := json.Marshal(bestProfessors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling professor")
		return
	}
	w.Write(ret)
}

// TODO: I would like this to also update the scoreboard file whenever I create that.
// update function for avgs
// func UpdateDerivedDataHandler(lastUpdate time.Time, interval time.Duration) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// time limit of
// 		if time.Since(lastUpdate) < interval {
// 			http.Error(w, "Updates too frequent", http.StatusTooManyRequests)
// 			return
// 		}

// 		go func() {
// 			err := telemetry.FigureDerivedData()
// 			if err != nil {
// 				log.Println("Error updating derived data:", err)
// 			}
// 		}()
// 		lastUpdate = time.Now()
// 		w.WriteHeader(http.StatusOK)
// 	}
// }
