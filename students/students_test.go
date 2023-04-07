package students

import (
	"log"
	"time"

	th "students/testhelpers"
	"testing"
)

func TestHello(t *testing.T) {
	want := "Hello, world."
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

func Test_StudentsCrud(t *testing.T) {
	var (
		name1           = "Mittens"
		currentYear     = 11
		currentYear2    = 10
		graduationYear  = 2024
		graduationYear2 = 2025
		avgGPA          = 3.8
		avgGPA2         = 1.7
		age             = 16
		//age2 	= 17
		dob = time.Date(2007, time.February, 8, 4, 5, 5, 5, time.Local)
	)
	// Create
	student, err := CreateNewStudent(name1, currentYear, graduationYear, avgGPA, age, dob, true)
	if err != nil {
		log.Println(" err : ", err)
	}
	// Get
	studentMittens, err := GetStudent(student.StudentID)
	if err != nil {
		log.Println(" err : ", err)
	}
	// value checks
	if len(studentMittens.StudentID) > 0 {
		th.AssertEqual(t, "sudent name : ", studentMittens.Name, name1)
		th.AssertEqual(t, "Current Year : ", studentMittens.CurrentYear, currentYear)
		th.AssertEqual(t, "Graduation Year", studentMittens.GraduationYear, graduationYear)
		th.AssertEqual(t, "Avg Gpa: ", studentMittens.AvgGPA, avgGPA)
		th.AssertEqual(t, "Age: ", studentMittens.Age, age)

	} else {
		log.Println("Get did not return a value skipping get tests. (db on?)")
	}
	// #demotion
	opts := UpdateStudentOptions{
		StudentID:      student.StudentID,
		CurrentYear:    currentYear2,
		GraduationYear: graduationYear2,
		AvgGPA:         avgGPA2,
	}

	err = opts.UpdateStudent()
	if err != nil {
		log.Println(" err : ", err)
	}
	studentMittens, err = GetStudent(student.StudentID)
	if err != nil {
		log.Println(" err : ", err)
	}
	if len(studentMittens.StudentID) > 0 {
		th.AssertEqual(t, "sudent name : ", studentMittens.Name, name1)
		th.AssertEqual(t, "Current Year : ", studentMittens.CurrentYear, currentYear2)
		th.AssertEqual(t, "Graduation Year", studentMittens.GraduationYear, graduationYear2)
		th.AssertEqual(t, "Avg Gpa: ", studentMittens.AvgGPA, avgGPA2)
		th.AssertEqual(t, "Age: ", studentMittens.Age, age)

	} else {
		log.Println("Get did not return a value skipping update tests. (db on?)")
	}
	err = DeleteStudent(studentMittens.StudentID)
	if err != nil {
		log.Println(" err : ", err)
	}
	_, err = GetStudent(studentMittens.StudentID)
	if err == nil {
		log.Fatal("get succeeded. It should not have.")
	}

}
