package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"students/telemetry"

	"github.com/gorilla/mux"
)

func CreateSchoolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	num := vars["num_per_grade"]
	numPerGrade, err := strconv.Atoi(num)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	} else if numPerGrade > 100 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	err = telemetry.NewSchool(numPerGrade)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
}
func UpdateDerivedData(w http.ResponseWriter, r *http.Request) {
	err := telemetry.FigureDerivedData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
}
