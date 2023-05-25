package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"students/client"
	"students/students"

	"github.com/gorilla/mux"
)

//decided no create school handler

// GetAllSchools
func GetAllSchools(w http.ResponseWriter, r *http.Request) {

	schools, err := students.GetAllSchools()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		fmt.Fprint(w, "Internal server error")
		return
	}
	fmt.Println("GetALLSchools: all schools returned")
	ret, err := json.Marshal(client.SchoolsToAPI(schools))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling school")
		return
	}
	w.Write(ret)
}

// GetClassesForSchoolHandler
func GetClassesForSchoolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schoolID := vars["school_id"]

	classes, err := students.GetClassesForSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error retrieving classes")
		return
	}
	ret, err := json.Marshal(client.ClassesToAPI(classes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling classes")
		return
	}
	w.Write(ret)
}
func GetStudentsForSchoolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schoolID := vars["school_id"]

	students, err := students.GetStudentsForSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		fmt.Fprint(w, "Unexpected error retrieving classes")
		return
	}
	ret, err := json.Marshal(client.StudentsToAPI(students))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling students")
		return
	}
	w.Write(ret)
}

// GetSchoolHandler
func GetSchoolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schoolID := vars["school_id"]
	if len(schoolID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	school, err := students.GetSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "school not found")
		return
	}
	ret, err := json.Marshal(client.SchoolToAPI(school))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling school")
		return
	}
	w.Write(ret)
}

// GetSchoolForUserHandler
func GetAllSchoolsForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerID := vars["owner_id"]
	if len(ownerID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	schools, err := students.GetAllSchoolsForUser(ownerID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("GetAllSchoolsForUser: %v", err)
		fmt.Println()
		fmt.Fprint(w, "ownerid not found")
		return
	}
	ret, err := json.Marshal(client.SchoolsToAPI(schools))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling school")
		return
	}
	w.Write(ret)
}

// UpdateSchoolHandler
func UpdateSchoolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schoolID := vars["school_id"]
	if len(schoolID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	//disable this while dev
	// ownerID := r.FormValue("owner_id")
	// school, err := students.GetSchool(schoolID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	fmt.Fprint(w, "school not found")
	// 	return
	// }
	// if school.OwnerID != ownerID{
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	fmt.Fprint(w,"invalid permission")
	// 	return
	// }
	opts := students.UpdateSchoolOptions{
		SchoolID:   schoolID,
		SchoolName: r.FormValue("name"),
	}

	// add to prof list
	addToProf := r.PostFormValue("add_to_professor")
	if addToProf != "" {
		addToProfList := strings.Split(addToProf, ",")
		fmt.Println(addToProfList)
		opts.AddToProfessorList = addToProfList
	}
	// remove from prof list
	removeFromProf := r.PostFormValue("remove_from_professor")
	if removeFromProf != "" {
		removeFromProfList := strings.Split(removeFromProf, ",")
		fmt.Println(removeFromProfList)
		opts.RemoveFromProfessorList = removeFromProfList
	}
	// add to class list
	addToClass := r.PostFormValue("add_to_class")
	if addToClass != "" {
		addToClassList := strings.Split(addToClass, ",")
		fmt.Println(addToClassList)
		opts.AddToClassList = addToClassList
	}
	// remove from class list
	removeFromClass := r.PostFormValue("remove_from_class")
	if removeFromClass != "" {
		removeFromClassList := strings.Split(removeFromClass, ",")
		fmt.Println(removeFromClassList)
		opts.RemoveFromClassList = removeFromClassList
	}
	// students
	addToRoster := r.PostFormValue("add_to_roster")
	if addToRoster != "" {
		addToRosterList := strings.Split(addToRoster, ",")
		fmt.Println(addToRosterList)
		opts.AddToStudentList = addToRosterList
	}
	removeFromRoster := r.PostFormValue("remove_from_roster")
	if removeFromRoster != "" {
		removeFromRosterList := strings.Split(removeFromRoster, ",")
		fmt.Println(removeFromRosterList)
		opts.RemoveFromStudentList = removeFromRosterList
	}
	err := opts.UpdateSchool()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteSchoolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schoolID := vars["school_id"]
	if len(schoolID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	err := students.DeleteSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "school not found")
		return
	}
	w.WriteHeader(http.StatusOK)
}
