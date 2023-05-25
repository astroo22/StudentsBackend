package telemetry

// import (
// 	"fmt"
// 	"students/students"
// 	th "students/testhelpers"
// 	"testing"
// )

// func Test_Telemetry(t *testing.T) {
// 	var (
// 		studentNUM = 30
// 	)
// 	studentList, professors, classes, reportcards, err := GenerateData(studentNUM, 1)
// 	if err != nil {
// 		t.Errorf("generation failed: %v", err)
// 	}
// 	err = BatchUploadTestData(studentList, professors, classes, reportcards, nil)
// 	if err != nil {
// 		t.Errorf("upload failed: %v", err)
// 	}
// 	if !th.AssertEqual(t, "students len", len(studentList), studentNUM) {
// 		t.Fatal("not enough students will create panic if continue")
// 	}
// 	//err = FigureDerivedData()
// 	if err != nil {
// 		t.Errorf("update of derived data failed")
// 	}
// 	th.AssertEqual(t, "reportcards len", len(reportcards), studentNUM)
// 	th.AssertEqual(t, "professors len", len(professors), 5)
// 	student, err := students.GetStudent(studentList[0].StudentID)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		th.AssertEqual(t, "student name", student.Name, studentList[0].Name)
// 		th.AssertEqual(t, "student age", student.Age, studentList[0].Age)
// 		th.AssertEqual(t, "student avg gpa", student.AvgGPA, studentList[0].AvgGPA)
// 		if student.AvgGPA == 0.0 {
// 			t.Errorf("avggpa didn't compile")
// 		}
// 		th.AssertEqual(t, "student current year", student.CurrentYear, studentList[0].CurrentYear)
// 		th.AssertEqual(t, "student name", student.GraduationYear, studentList[0].GraduationYear)
// 	}
// 	student, err = students.GetStudent(studentList[(studentNUM - 1)].StudentID)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		th.AssertEqual(t, "student name", student.Name, studentList[(studentNUM-1)].Name)
// 		th.AssertEqual(t, "student age", student.Age, studentList[(studentNUM-1)].Age)
// 		th.AssertEqual(t, "student avg gpa", student.AvgGPA, studentList[(studentNUM-1)].AvgGPA)
// 		th.AssertEqual(t, "student current year", student.CurrentYear, studentList[(studentNUM-1)].CurrentYear)
// 		th.AssertEqual(t, "student name", student.GraduationYear, studentList[(studentNUM-1)].GraduationYear)
// 	}

// 	prof, err := students.GetProfessor(professors[0].ProfessorID)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		th.AssertEqual(t, "prof name", prof.Name, professors[0].Name)
// 	}
// 	_, err = UpdateProfessorStudentAvg(professors[0].ProfessorID)
// 	if err != nil {
// 		t.Errorf("prof get student avg failed: %v", err)
// 	}
// 	//fmt.Println(prof)
// 	//fmt.Println(avg)
// }

// // test of updates and some more complicated SQL
// // This a test that is much easier to see the results of from the db view
// // I might connect this to a pipeline later on so I'm commenting out this test
// func Test_AvgUpdates(t *testing.T) {

// 	var (
// 		//profList  []students.Professor
// 		//profListn []students.Professor
// 		// schoolID so dont have to generate data
// 		schoolID = "78e046f3-9220-4390-b6af-5f236afc24c7"
// 		// numStdsGrade = 5
// 		// owner_id     = uuid.New().String()
// 		// schoolName   = "Rick and Morty Vindicators 4 "
// 	)

// 	// will use all 30 seconds of timeout if uncommented
// 	// err := UpdateAllSchoolAvgGpa()
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }
// 	// t.Fatal("completed hopefully?")

// 	// err := students.UpdateSchoolRankings()
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }
// 	// t.Fatal("completed hopefully?")

// 	// get school
// 	school, err := students.GetSchool(schoolID)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// get the classes
// 	classes, err := students.GetClasses(school.ClassList)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	fmt.Println(school.ProfessorList)
// 	// run an update
// 	err = UpdateClassAvgs(classes)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// get the grade avgs.
// 	_, err = GetGradeAvgForSchool(schoolID)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	for _, id := range school.ProfessorList {
// 		prof, err := students.GetProfessor(id)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		fmt.Println(prof.StudentAvg)
// 		//profList = append(profList, prof)
// 	}
// 	err = UpdateProfessorsStudentAvgs(school.ProfessorList)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	for _, id := range school.ProfessorList {
// 		prof, err := students.GetProfessor(id)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		fmt.Println(prof.StudentAvg)
// 		//profListn = append(profList, prof)
// 	}
// 	// TODO: These currently are manual tests where I go and look at the records updated because of timeout issues.
// 	// this needs to be updated with a create however the problem is that the standard test function timeout is only 30 seconds
// 	// and this throws import errors if I try to run only the tests in this package. So it forces to me run all the tests or it fails?
// 	// Theres alot of tests so this causes excess down time just to be able to set the timeout flag therefor I will redo this test later on.
// 	// NOTES: I read a bit about import errors in the testing package comments while looking through it.
// 	// it didn't really give much information on it. Current theory is nesting that folder within the students package has caused a
// 	// structure issue with how the testing package works. Possible solution is pulling the telemetry folder out which I can do but
// 	// there was a reason I moved it in the first place I just forgot that reason at the moment. However since it is functional and not critial
// 	// this issue will be ignored until I do a code review once everything is up and running as with many of my other TODO tags.
// }
