package telemetry

import (
	"fmt"
	"students/students"
	th "students/testhelpers"
	"testing"

	"github.com/google/uuid"
)

func Test_Telemetry(t *testing.T) {
	var (
		studentNUM = 30
	)
	studentList, professors, classes, reportcards, err := GenerateData(studentNUM, 1)
	if err != nil {
		t.Errorf("generation failed: %v", err)
	}
	err = BatchUploadTestData(studentList, professors, classes, reportcards, nil)
	if err != nil {
		t.Errorf("upload failed: %v", err)
	}
	if !th.AssertEqual(t, "students len", len(studentList), studentNUM) {
		t.Fatal("not enough students will create panic if continue")
	}
	//err = FigureDerivedData()
	if err != nil {
		t.Errorf("update of derived data failed")
	}
	th.AssertEqual(t, "reportcards len", len(reportcards), studentNUM)
	th.AssertEqual(t, "professors len", len(professors), 5)
	student, err := students.GetStudent(studentList[0].StudentID)
	if err != nil {
		t.Error(err)
	} else {
		th.AssertEqual(t, "student name", student.Name, studentList[0].Name)
		th.AssertEqual(t, "student age", student.Age, studentList[0].Age)
		th.AssertEqual(t, "student avg gpa", student.AvgGPA, studentList[0].AvgGPA)
		if student.AvgGPA == 0.0 {
			t.Errorf("avggpa didn't compile")
		}
		th.AssertEqual(t, "student current year", student.CurrentYear, studentList[0].CurrentYear)
		th.AssertEqual(t, "student name", student.GraduationYear, studentList[0].GraduationYear)
	}
	student, err = students.GetStudent(studentList[(studentNUM - 1)].StudentID)
	if err != nil {
		t.Error(err)
	} else {
		th.AssertEqual(t, "student name", student.Name, studentList[(studentNUM-1)].Name)
		th.AssertEqual(t, "student age", student.Age, studentList[(studentNUM-1)].Age)
		th.AssertEqual(t, "student avg gpa", student.AvgGPA, studentList[(studentNUM-1)].AvgGPA)
		th.AssertEqual(t, "student current year", student.CurrentYear, studentList[(studentNUM-1)].CurrentYear)
		th.AssertEqual(t, "student name", student.GraduationYear, studentList[(studentNUM-1)].GraduationYear)
	}

	prof, err := students.GetProfessor(professors[0].ProfessorID)
	if err != nil {
		t.Error(err)
	} else {
		th.AssertEqual(t, "prof name", prof.Name, professors[0].Name)
	}
	_, err = UpdateProfessorStudentAvg(professors[0].ProfessorID)
	if err != nil {
		t.Errorf("prof get student avg failed: %v", err)
	}
	//fmt.Println(prof)
	//fmt.Println(avg)
}

// test of derived data and getgradeavgforschool
func Test_AvgUpdates(t *testing.T) {

	var (
		// schoolID so dont have to generate data
		schoolID = "78e046f3-9220-4390-b6af-5f236afc24c7"
		// numStdsGrade = 5
		// owner_id     = uuid.New().String()
		// schoolName   = "Rick and Morty Vindicators 4 "
	)
	// school, err := NewSchool(numStdsGrade, owner_id, schoolName)
	// if err != nil {
	// 	t.Error(err)
	// }
	// avgs, err := GetGradeAvgForSchool(schoolID)
	// if err != nil {
	// 	t.Error(err)
	// }

	school, err := students.GetSchool(schoolID)
	if err != nil {
		t.Error(err)
	}
	classes, err := students.GetClasses(school.ClassList)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("studentID: %s", classes[0].Roster[0])
	fmt.Println()
	rc, err := students.GetReportCard(classes[0].Roster[0])
	if err != nil {
		t.Error(err)
	}
	fmt.Println("where am i?")
	fmt.Println(rc)
	fmt.Println(rc.English)
	err = UpdateClassAvgs(classes)
	if err != nil {
		t.Error(err)
	}
	// avgs, err := GetGradeAvgForSchool(schoolID)
	// if err != nil {
	// 	t.Error(err)
	// }
	// fmt.Println(avgs)
	t.Error(err)
}

func Test_TelemetryMassUpdates(t *testing.T) {
	// err := UpdateStudentAvgs()
	// if err != nil {
	// 	t.Error(err)
	// }
	// 	err := DeleteTables()
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	//school names might write into creating script later:
	// "PLUS ULTRA ACADEMY"
	// "The Xavier Institute for Higher Learning #8"
	// "Busta Rhymes Academy"
	// "Fosters Home for Imaginary Friends"
	//newschool

	_, err := NewSchool(10, uuid.New().String(), "Fosters Home for Imaginary Friends")
	if err != nil {
		t.Error(err)
	}
	// if err != nil {
	// 	t.Error(err)
	// }
}
