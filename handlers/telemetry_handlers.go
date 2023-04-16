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

// TODO: I would like this to also update the scoreboard file whenever I create that.
// update function for avgs
func UpdateDerivedData(w http.ResponseWriter, r *http.Request) {
	err := telemetry.FigureDerivedData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
}
