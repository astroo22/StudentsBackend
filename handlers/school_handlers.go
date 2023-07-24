package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"students/client"
	"students/logger"
	"students/students"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	userID := r.Context().Value("user_id")
	if userID == nil {
		// Handle error: no user ID in context
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Please log in")
		return
	}
	if ownerID != userID {
		logger.Log.WithFields(logrus.Fields{
			"owner_id": ownerID,
			"user_id":  userID,
		}).Warn("attempted get of non owned object")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "non authorized get attempt of unowned objects")
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
		fmt.Println("no id")
		return
	}
	type SchoolData struct {
		OwnerID               string   `json:"owner_id"`
		Enrollment_change_ids []string `json:"enrollment_change_ids"`
	}
	var data SchoolData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		fmt.Println("hit")
		return
	}
	fmt.Println(data.OwnerID)

	//disable this while dev
	// ownerID := r.FormValue("owner_id")
	// if len(ownerID) == 0 {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprint(w, "Invalid request payload")
	// 	fmt.Println("hit")
	// 	return
	// }
	userID := r.Context().Value("user_id")
	if userID == nil {
		// Handle error: no user ID in context
		fmt.Println(userID)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Please log in")
		return
	}
	if data.OwnerID != userID {
		logger.Log.WithFields(logrus.Fields{
			"owner_id": data.OwnerID,
			"user_id":  userID,
		}).Warn("attempted update of non owned object")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "non authorized get attempt of unowned objects")
		return
	}

	opts := students.UpdateSchoolOptions{
		SchoolID: schoolID,
	}
	if len(data.Enrollment_change_ids) > 0 {
		err = students.FlipEnrollmentStatus(data.Enrollment_change_ids)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("flip")
			fmt.Println(err)
			return
		}
		go func() {
			_, err = students.UpdateSchoolAvg(schoolID)
			if err != nil {
				fmt.Println("error updating school avg: not normal error")
				fmt.Println(err)
				return
			}
			fmt.Println("finished avgs")
		}()

	}
	if len(opts.SchoolName) > 0 {
		fmt.Println(opts)
		err = opts.UpdateSchool()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server error")
			fmt.Println(err)
			return
		}
	}
	fmt.Println("made it through update~")
	w.WriteHeader(http.StatusOK)

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// everything from this point on could be future features
	// to be included in an update. over coded here

	// add to prof list
	// addToProf := r.PostFormValue("add_to_professor")
	// if addToProf != "" {
	// 	addToProfList := strings.Split(addToProf, ",")
	// 	fmt.Println(addToProfList)
	// 	opts.AddToProfessorList = addToProfList
	// }
	// // remove from prof list
	// removeFromProf := r.PostFormValue("remove_from_professor")
	// if removeFromProf != "" {
	// 	removeFromProfList := strings.Split(removeFromProf, ",")
	// 	fmt.Println(removeFromProfList)
	// 	opts.RemoveFromProfessorList = removeFromProfList
	// }
	// // add to class list
	// addToClass := r.PostFormValue("add_to_class")
	// if addToClass != "" {
	// 	addToClassList := strings.Split(addToClass, ",")
	// 	fmt.Println(addToClassList)
	// 	opts.AddToClassList = addToClassList
	// }
	// // remove from class list
	// removeFromClass := r.PostFormValue("remove_from_class")
	// if removeFromClass != "" {
	// 	removeFromClassList := strings.Split(removeFromClass, ",")
	// 	fmt.Println(removeFromClassList)
	// 	opts.RemoveFromClassList = removeFromClassList
	// }
	// // students
	// addToRoster := r.PostFormValue("add_to_roster")
	// if addToRoster != "" {
	// 	addToRosterList := strings.Split(addToRoster, ",")
	// 	fmt.Println(addToRosterList)
	// 	opts.AddToStudentList = addToRosterList
	// }
	// removeFromRoster := r.PostFormValue("remove_from_roster")
	// if removeFromRoster != "" {
	// 	removeFromRosterList := strings.Split(removeFromRoster, ",")
	// 	fmt.Println(removeFromRosterList)
	// 	opts.RemoveFromStudentList = removeFromRosterList
	// }

}
func UpdateSchoolAvgHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schoolID := vars["school_id"]
	if len(schoolID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}

	// need to spin off go routine
	go func() {
		_, err := students.UpdateSchoolAvg(schoolID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "internal server error")
			return
		}
	}()

	// ret, err := json.Marshal(newgpa)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprint(w, "Unexpected error mashalling school")
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
}

func DeleteSchoolHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit in delete")
	type deleteSchoolRequest struct {
		OwnerID string `json:"owner_id"`
	}
	vars := mux.Vars(r)
	schoolID := vars["school_id"]
	if len(schoolID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	reqBody := &deleteSchoolRequest{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to parse request body")
		return
	}
	ownerID := reqBody.OwnerID
	if len(ownerID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	userID := r.Context().Value("user_id")
	if userID == nil {
		// Handle error: no user ID in context
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Please log in")
		return
	}
	if ownerID != userID {
		logger.Log.WithFields(logrus.Fields{
			"owner_id": ownerID,
			"user_id":  userID,
		}).Warn("attempted update of non owned object")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "non authorized get attempt of unowned objects")
		return
	}
	err = students.DeleteSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "school not found")
		return
	}
	fmt.Println("actually deleted?")
	w.WriteHeader(http.StatusOK)
}