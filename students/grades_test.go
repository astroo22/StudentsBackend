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
		th.AssertEqual(t, "PhysicalED", newReportCard.PhysicalED, reportCardsM[0].PhysicalED)
		th.AssertEqual(t, "Lunch", newReportCard.Lunch, reportCardsM[0].Lunch)
	} else {
		fmt.Println("didn't get report card skipping tests")
	}
	fmt.Println("Create Read complete")
	err = DeleteBatchReportCard(students)
	if err != nil {
		log.Fatal(err)
	}
	_, err = GetReportCard(students[0].StudentID)
	if err == nil {
		log.Fatal("get succeeded shouldn't have. sad")
	}

}
