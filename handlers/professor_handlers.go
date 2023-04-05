package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"students/students"

	"github.com/gorilla/mux"
)

func CreateProfessorHandler(w http.ResponseWriter, r *http.Request) {
	//profID := r.FormValue("professor_id")
	name := r.FormValue("name")

	professor, err := students.CreateProfessor(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error creating professor")
		return
	}
	ret, err := json.Marshal(professor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling professor")
		return
	}
	w.Write(ret)
}

func GetProfessorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profID := vars["professor_id"]
	professor, err := students.GetProfessor(profID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error retrieving professor")
		return
	}
	if len(professor.ProfessorID) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "professor not found")
		return
	}
	ret, err := json.Marshal(professor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling professor")
		return
	}
	w.Write(ret)
}

func UpdateProfessorHandler(w http.ResponseWriter, r *http.Request) {
	opts := students.UpdateProfessorOptions{}

	vars := mux.Vars(r)
	profID := vars["professor_id"]
	if len(profID) <= 30 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "professor not found")
		return
	} else {
		// do a get check if exists else return
		_, err := students.GetProfessor(profID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "No class with that ID")
			return
		}
		opts.ProfessorID = profID
	}
	studentAvg, err := strconv.ParseFloat(r.FormValue("student_avg"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid value")
	} else if studentAvg > 0.0 && studentAvg <= 4.00 {
		fmt.Println(studentAvg)
		opts.StudentAvg = studentAvg
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid value")
	}
	classes := r.PostFormValue("class_list")
	if classes != "" {
		classList := strings.Split(classes, ",")
		if len(classList) > 0 {
			fmt.Println(classList)
			opts.ClassList = classList
		}
	}
	fmt.Println(opts)
	err = opts.UpdateProfessor()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteProfessorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["professor_id"]
	// test of uuidparse
	// profID, err := uuid.Parse(ID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprint(w, "incorrect ID format")
	// 	return
	// }
	err := students.DeleteProfessor(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error deleting class")
		return
	}
	w.WriteHeader(http.StatusOK)
}
