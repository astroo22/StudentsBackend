package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"students/client"
	"students/students"

	"github.com/gorilla/mux"
)

func CreateClassHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	classGrade, err := strconv.Atoi(r.FormValue("teaching_grade"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	profID := r.FormValue("professor_id")
	subject := r.FormValue("subject")
	roster := r.PostFormValue("roster")
	rosterList := strings.Split(roster, ",")
	class, err := students.CreateClass(classGrade, profID, subject, rosterList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error creating class")
		return
	}
	ret, err := json.Marshal(client.ClassToAPI(class))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling class")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func GetClassHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classID := vars["class_id"]
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
	ret, err := json.Marshal(client.ClassToAPI(class))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling class")
		return
	}
	w.Write(ret)
}

func UpdateClassHandler(w http.ResponseWriter, r *http.Request) {
	opts := students.UpdateClassOptions{}
	// class id
	vars := mux.Vars(r)
	classID := vars["class_id"]
	if len(classID) <= 30 {
		//http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Class not found")
		return
	} else {
		// do a get check if exists else returnz
		_, err := students.GetClass(classID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "No class with that ID")
			return
		}
		opts.ClassID = classID
	}
	// prof id
	profID := r.FormValue("professor_id")
	if len(profID) < 30 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "professor id")
	} else {
		opts.ProfessorID = profID
	}
	classAvg, err := strconv.ParseFloat(r.FormValue("class_avg"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid value")
	} else if classAvg > 0.0 && classAvg <= 4.00 {
		opts.ClassAvg = classAvg
	}

	// roster add
	addRoster := r.PostFormValue("add_roster")
	if addRoster != "" {
		addRosterList := strings.Split(addRoster, ",")
		//fmt.Println(len(addRosterList))
		opts.AddRoster = addRosterList
	}

	// roster remove
	removeRoster := r.PostFormValue("remove_roster")
	if removeRoster != "" {
		removeRosterList := strings.Split(removeRoster, ",")
		opts.RemoveRoster = removeRosterList
	}
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
