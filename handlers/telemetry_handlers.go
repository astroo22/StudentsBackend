package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	// Call the GetGradeAvgForSchool function to retrieve the data
	gradeAvgList, err := telemetry.GetGradeAvgForSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting grade averages: %v", err)
		return
	}
	ret, err := json.Marshal(gradeAvgList)
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
