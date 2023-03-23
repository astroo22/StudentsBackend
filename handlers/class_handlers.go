package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"students/students"

	"github.com/gorilla/mux"
)

func CreateClassHandler(w http.ResponseWriter, r *http.Request) {
	classGrade, err := strconv.Atoi(r.FormValue("teaching_grade"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	profID := r.FormValue("professor_id")
	subject := r.FormValue("subject")
	roster := r.PostForm["roster"]

	class, err := students.CreateClass(classGrade, profID, subject, roster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error creating class")
		return
	}
	ret, err := json.Marshal(class)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling class")
		return
	}
	w.Write(ret)
}

func GetClassHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classID := vars["class_id"]
	//classID := r.FormValue("class_id")
	if len(classID) <= 30 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Class not found")
		return
	}

	class, err := students.GetClass(classID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error retrieving class")
		return
	}

	if len(class.ClassID) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Class not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}

func UpdateClassHandler(w http.ResponseWriter, r *http.Request) {
	opts := students.UpdateClassOptions{}
	// class id
	vars := mux.Vars(r)
	classID := vars["class_id"]
	//classID := r.FormValue("class_id")
	if len(classID) <= 30 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Class not found")
		return
	} else {
		// do a get check if exists else return
		_, err := students.GetClass(classID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "No class with that ID")
			return
		}
		opts.ClassID = classID
	}
	// prof id
	profID := r.FormValue("prof_id")
	if len(profID) <= 30 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Class not found")
		return
	} else {
		opts.ProfessorID = profID
	}

	classAvg, err := strconv.ParseFloat(r.FormValue("class_avg"), 64)
	if err != nil {
		fmt.Fprint(w, "invalid value")
	} else if classAvg > 0.0 && classAvg <= 4.00 {
		opts.ClassAvg = classAvg
	}
	// TODO: come up with a good value check for roster is uuid of student_id
	opts.Roster = r.PostForm["roster"]

	err = opts.UpdateClass()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteClassHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classID := vars["class_id"]
	err := students.DeleteClass(classID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error deleting class")
		return
	}

	w.WriteHeader(http.StatusOK)
}
