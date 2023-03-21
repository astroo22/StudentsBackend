package students

import (
	"fmt"
	"log"
	th "students/testhelpers"
	"testing"
)

func Test_GradesCrud(t *testing.T) {
	students := GenerateTestData()
	err := CreateNewStudents(students)
	if err != nil {
		log.Fatal(err)
	}
	reportCardsM := []ReportCard{}
	for _, student := range students {
		reportCard, err := CreateReportCard(student.StudentID)
		if err != nil {
			log.Fatal(err)
		}
		reportCardsM = append(reportCardsM, reportCard)
	}

	th.AssertEqual(t, "len reportcards generated", len(reportCardsM), 10)
	newReportCard, err := GetReportCard(reportCardsM[0].StudentID)
	if err != nil {
		log.Fatal(err)
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
	fmt.Println("Create Read complete")
	// dat 4.0 tho
	opts := UpdateReportCardOptions{
		StudentID:  students[0].StudentID,
		Math:       4.0,
		Science:    4.0,
		English:    4.0,
		PhysicalED: 4.0,
		Lunch:      4.0,
	}
	err = opts.UpdateReportCard()
	if err != nil {
		log.Fatal(err)
	}
	updatedReport, err := GetReportCard(students[0].StudentID)
	if err != nil {
		log.Fatal(err)
	}
	if len(updatedReport.StudentID) > 0 {
		th.AssertEqual(t, "Math", updatedReport.Math, 4.0)
		th.AssertEqual(t, "Science", updatedReport.Science, 4.0)
		th.AssertEqual(t, "English", updatedReport.English, 4.0)
		th.AssertEqual(t, "Physical_ed", updatedReport.PhysicalED, 4.0)
		th.AssertEqual(t, "Lunch", updatedReport.Lunch, 4.0)
		// TODO: add testing for the classlist array once we have classes and tests set up
		// also add some update variations in here
	}
	err = DeleteBatchReportCard(students)
	if err != nil {
		log.Fatal(err)
	}
	_, err = GetReportCard(students[0].StudentID)
	if err == nil {
		log.Fatal("get succeeded shouldn't have. sad")
	}
}
