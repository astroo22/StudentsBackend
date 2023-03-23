package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"students/students"
	"time"

	"github.com/gorilla/mux"
)

func CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	currentYear, err := strconv.Atoi(r.FormValue("current_year"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	graduationYear, err := strconv.Atoi(r.FormValue("graduation_year"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	gpa, err := strconv.ParseFloat(r.FormValue("gpa"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}

	age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	dobStr := r.FormValue("dob")
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid dob parameter")
		return
	}
	enrolled, err := strconv.ParseBool(r.FormValue("enrolled"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}

	student, err := students.CreateNewStudent(name, currentYear, graduationYear, gpa, age, dob, enrolled)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal server error")
		return
	}
	fmt.Println(student)
	// TODO: this tomorrow. get handler tests working CreateNewStudent should probably return student object.
	// Need to incorperate the API file then create all the other handler files
	// then finally the telemetry package and we can do generation and stuff
	// ret, err := json.Marshal(Student{
	// 	StudentID:      student.StudentID,
	// 	Name:           student.Name,
	// 	CurrentYear:    student.CurrentYear,
	// 	GraduationYear: student.GraduationYear,
	// 	GPA:            student.GPA,
	// 	Age:            student.Age,
	// 	DOB:            student.DOB,
	// 	Enrolled:       student.Enrolled,
	// })
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling class")
		return
	}
	//w.Write(ret)
}
func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	//studentID := r.FormValue("student_id")
	if len(studentID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	student, err := students.GetStudent(studentID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Student not found")
		return
	}
	ret, err := json.Marshal(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unexpected error mashalling class")
		return
	}
	w.Write(ret)
}
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	opts := students.UpdateStudentOptions{}
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	if len(studentID) <= 30 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	} else {
		_, err := students.GetStudent(studentID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "No student with that ID")
			return
		}
		opts.StudentID = studentID
	}
	name := r.FormValue("name")
	if len(name) > 0 {
		opts.Name = name
	}
	currentYear, err := strconv.Atoi(r.FormValue("current_year"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	if currentYear > 0 && currentYear < 13 {
		opts.CurrentYear = currentYear
	}
	graduationYear, err := strconv.Atoi(r.FormValue("graduation_year"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
	}
	if graduationYear > 0 && graduationYear < 13 {
		opts.GraduationYear = graduationYear
	}
	gpa, err := strconv.ParseFloat(r.FormValue("avg_gpa"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
	}
	if gpa > 0.0 && gpa <= 4.00 {
		opts.AvgGPA = gpa
	}
	age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	if age > 0 && age < 110 {
		opts.Age = age
	}
	enrolled, err := strconv.ParseBool(r.FormValue("enrolled"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}
	opts.Enrolled = enrolled
	err = opts.UpdateStudent()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}
func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	err := students.DeleteStudent(studentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Invalid request payload")
		return
	}

}
