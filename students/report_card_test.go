package students

import (
	"fmt"
	"log"
	th "students/testhelpers"
	"testing"
)

func Test_ReportCardCrud(t *testing.T) {
	finalClassList := []string{"math", "science", "english", "lunch"}
	students := GenerateTestData()
	err := CreateNewStudents(students)
	if err != nil {
		log.Println(" err : ", err)
	}
	reportCardsM := []ReportCard{}
	for _, student := range students {
		reportCard, err := CreateReportCard(student.StudentID)
		if err != nil {
			log.Println(" err : ", err)
		}
		reportCardsM = append(reportCardsM, reportCard)
	}

	th.AssertEqual(t, "len reportcards generated", len(reportCardsM), 10)
	newReportCard, err := GetReportCard(reportCardsM[0].StudentID)
	if err != nil {
		log.Println(" err : ", err)
	}
	if len(newReportCard.StudentID) > 0 {
		th.AssertEqual(t, "Math", newReportCard.Math, reportCardsM[0].Math)
		th.AssertEqual(t, "Science", newReportCard.Science, reportCardsM[0].Science)
		th.AssertEqual(t, "English", newReportCard.English, reportCardsM[0].English)
		th.AssertEqual(t, "Physical_ed", newReportCard.PhysicalED, reportCardsM[0].PhysicalED)
		th.AssertEqual(t, "Lunch", newReportCard.Lunch, reportCardsM[0].Lunch)
	} else {
		fmt.Println("didn't get report card skipping tests")
	}

	// dat 4.0 tho
	tempClassList := []string{"math", "science", "english", "football"}
	opts := UpdateReportCardOptions{
		StudentID:    students[0].StudentID,
		Math:         4.0,
		Science:      4.0,
		English:      4.0,
		PhysicalED:   4.0,
		Lunch:        4.0,
		AddClassList: tempClassList,
	}
	err = opts.UpdateReportCard()
	if err != nil {
		log.Println(" err : ", err)
	}
	updatedReport, err := GetReportCard(students[0].StudentID)
	if err != nil {
		log.Println(" err : ", err)
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
		StudentID:       students[0].StudentID,
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
		log.Println(" err : ", err)
	}
	updatedReport, err = GetReportCard(students[0].StudentID)
	if err != nil {
		log.Println(" err : ", err)
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

	err = DeleteBatchReportCard(students)
	if err != nil {
		log.Println(" err : ", err)
	}
	_, err = GetReportCard(students[0].StudentID)
	if err == nil {
		log.Fatal("get succeeded shouldn't have. sad")
	}
}
