package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"students/client"
	"students/logger"
	"students/students"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func CreateNewSchoolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner_id := vars["owner_id"]
	numPerGrade, err := strconv.Atoi(r.FormValue("num_per_grade"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		fmt.Printf("err %v", err)
		return
	} else if numPerGrade > 100 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		fmt.Println(err)
		return
	}
	name := r.FormValue("name")
	logger.Log.WithFields(logrus.Fields{
		"owner_id": owner_id,
		"name":     name,
	}).Info("created a school")

	operationID := uuid.NewString()

	students.CreateOperationEntry(operationID, "in progress")

	go func() {
		new_school, err := students.NewSchool(operationID, numPerGrade, owner_id, name)
		if err != nil {
			students.UpdateOperationStatus(operationID, "error", err)
			return
		}
		// maybe issue here will just do a school for now see how goes
		students.UpdateOperationStatus(operationID, "Complete", err)
		school := new_school
		status := students.SchoolCreationStatus_API{
			Status: "Complete",
			School: school,
		}

		logger.Log.WithFields(logrus.Fields{
			"owner_id": owner_id,
			"name":     name,
		}).Info("Finished creating the school")

		ret, err := json.Marshal(status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Unexpected error marshalling response")
			return
		}
		w.Write(ret)
	}()

	status := students.SchoolCreationStatus_API{
		Status:      "School creation in progress",
		OperationID: operationID,
	}
	ret, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error marshalling response")
		return
	}
	w.Write(ret)
}

func SchoolCreationStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	operationID := vars["operation_id"]

	status, err := students.GetOperationStatus(operationID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ret, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error marshalling response")
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
	// check if exist
	_, err := students.GetSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "no school with that id")
		return
	}
	// attempt update  TODO: fix this function to auto check if the values are 0.0 and run a class avg update if so
	gradeAvgList, err := students.GetGradeAvgForSchool(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting grade averages: %v", err)
		return
	}
	// means update
	if gradeAvgList[0].AvgGPA == 0.0 {
		err = students.UpdateClassAvgs(schoolID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "error on class avg update")
			return
		}
		gradeAvgList, err = students.GetGradeAvgForSchool(schoolID)
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
	fmt.Println("hit GetBestProfs")
	vars := mux.Vars(r)
	schoolID := vars["school_id"]
	bestProfessors, err := students.GetBestProfessors(schoolID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("incorrect schoolID")
		return
	}
	if bestProfessors[0].StudentAvg == 0 {
		fmt.Println("Found values needing update for professors. Running updates!")
		var ids []string

		for _, prof := range bestProfessors {
			ids = append(ids, prof.ProfessorID)
		}

		err = students.UpdateProfessorsStudentAvgs(ids)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Unexpected error during value updates")
			fmt.Printf("logging error at professorupdate %v", err)
			return
		}
		bestProfessors, err = students.GetBestProfessors(schoolID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	ret, err := json.Marshal(client.ProfessorsToAPI(bestProfessors))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling professor")
		return
	}
	fmt.Println("completed GetBestProfs")
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}
