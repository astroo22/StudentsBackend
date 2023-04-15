package students

import (
	"fmt"
	"time"

	th "students/testhelpers"
	"testing"
)

func Test_ReportCardCrud(t *testing.T) {
	finalClassList := []string{"math", "science", "english", "lunch"}
	var (
		name1          = "Mittens"
		currentYear    = 11
		graduationYear = 2024
		avgGPA         = 3.8
		age            = 16
		dob            = time.Date(2007, time.February, 8, 4, 5, 5, 5, time.Local)
	)
	// Create
	student, err := CreateNewStudent(name1, currentYear, graduationYear, avgGPA, age, dob, true)
	if err != nil {
		t.Errorf(" err : %v", err)
	}

	reportCard, err := CreateReportCard(student.StudentID)
	if err != nil {
		t.Errorf(" err : %v", err)
	}

	newReportCard, err := GetReportCard(reportCard.StudentID)
	if err != nil {
		t.Errorf(" err : %v", err)
	}
	if len(newReportCard.StudentID) > 0 {
		th.AssertEqual(t, "Math", newReportCard.Math, reportCard.Math)
		th.AssertEqual(t, "Science", newReportCard.Science, reportCard.Science)
		th.AssertEqual(t, "English", newReportCard.English, reportCard.English)
		th.AssertEqual(t, "Physical_ed", newReportCard.PhysicalED, reportCard.PhysicalED)
		th.AssertEqual(t, "Lunch", newReportCard.Lunch, reportCard.Lunch)
	} else {
		fmt.Println("didn't get report card skipping tests")
	}

	// dat 4.0 tho
	tempClassList := []string{"math", "science", "english", "football"}
	opts := UpdateReportCardOptions{
		StudentID:    student.StudentID,
		Math:         4.0,
		Science:      4.0,
		English:      4.0,
		PhysicalED:   4.0,
		Lunch:        4.0,
		AddClassList: tempClassList,
	}
	err = opts.UpdateReportCard()
	if err != nil {
		t.Errorf(" err updating report card: %v", err)
	}
	updatedReport, err := GetReportCard(student.StudentID)
	if err != nil {
		t.Errorf(" err : %v", err)
	}
	if len(updatedReport.StudentID) > 0 {
		th.AssertEqual(t, "Math", updatedReport.Math, opts.Math)
		th.AssertEqual(t, "Science", updatedReport.Science, opts.Science)
		th.AssertEqual(t, "English", updatedReport.English, opts.English)
		th.AssertEqual(t, "Physical_ed", updatedReport.PhysicalED, opts.PhysicalED)
		th.AssertEqual(t, "Lunch", updatedReport.Lunch, opts.Lunch)
		fmt.Println(updatedReport.ClassList)
		if th.AssertEqual(t, "classlist len: ", len(updatedReport.ClassList), 4) {
			th.AssertEqual(t, "list check: ", updatedReport.ClassList, tempClassList)
		}
	}

	// Joined a gang: ruined the 4.0
	opts = UpdateReportCardOptions{
		StudentID:       student.StudentID,
		Math:            2.3,
		Science:         3.1,
		English:         3.4,
		PhysicalED:      4.0,
		Lunch:           4.0,
		AddClassList:    []string{"lunch"},
		RemoveClassList: []string{"football"},
	}
	err = opts.UpdateReportCard()
	if err != nil {
		t.Errorf(" err updating report card: %v", err)
	}
	updatedReport, err = GetReportCard(student.StudentID)
	if err != nil {
		t.Errorf(" err getting report card: %v", err)
	}
	if len(updatedReport.StudentID) > 0 {
		th.AssertEqual(t, "Math", updatedReport.Math, opts.Math)
		th.AssertEqual(t, "Science", updatedReport.Science, opts.Science)
		th.AssertEqual(t, "English", updatedReport.English, opts.English)
		th.AssertEqual(t, "Physical_ed", updatedReport.PhysicalED, opts.PhysicalED)
		th.AssertEqual(t, "Lunch", updatedReport.Lunch, opts.Lunch)
		if th.AssertEqual(t, "classlist len: ", len(updatedReport.ClassList), len(finalClassList)) {
			th.AssertEqual(t, "list check: ", updatedReport.ClassList, finalClassList)
		}
	}
	defer DeleteStudent(student.StudentID)

	err = DeleteReportCard(student.StudentID)
	if err != nil {
		t.Errorf(" err deleting report card : %v", err)
	}
	_, err = GetReportCard(student.StudentID)
	if err == nil {
		t.Error("get succeeded shouldn't have. sad")
	}
}
