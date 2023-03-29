package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"students/students"

	"github.com/gorilla/mux"
)

func CreateReportCardHandler(w http.ResponseWriter, r *http.Request) {
	//todo: add get check here
	studentID := r.FormValue("student_id")
	if len(studentID) <= 35 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	reportcard, err := students.CreateReportCard(studentID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}

	ret, err := json.Marshal(reportcard)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling card")
		return
	}
	w.Write(ret)
}

func GetReportCardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	if len(studentID) <= 35 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	grade, err := students.GetReportCard(studentID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Grade not found")
		return
	}
	ret, err := json.Marshal(grade)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling class")
		return
	}
	w.Write(ret)
}

func UpdateReportCardHandler(w http.ResponseWriter, r *http.Request) {
	opts := students.UpdateReportCardOptions{}
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	if len(studentID) <= 35 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	_, err := students.GetReportCard(studentID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No student with that ID")
		return
	}
	opts.StudentID = studentID

	math, err := strconv.ParseFloat(r.FormValue("math"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload math")
		return
	} else if math >= 0.00 && math <= 4.00 {
		opts.Math = math
	}
	science, err := strconv.ParseFloat(r.FormValue("science"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload math")
		return
	} else if science >= 0.00 && science <= 4.00 {
		opts.Science = science
	}
	english, err := strconv.ParseFloat(r.FormValue("english"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload math")
		return
	} else if english >= 0.00 && english <= 4.00 {
		opts.English = english
	}
	physicalED, err := strconv.ParseFloat(r.FormValue("physicaled"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload math")
		return
	} else if physicalED >= 0.00 && physicalED <= 4.00 {
		opts.PhysicalED = physicalED
	}
	lunch, err := strconv.ParseFloat(r.FormValue("lunch"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload math")
		return
	} else if lunch >= 0.00 && lunch <= 4.00 {
		opts.Lunch = lunch
	}

	err = opts.UpdateReportCard()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteReportCardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	err := students.DeleteReportCard(studentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
}
