package telemetry

// func Test_BatchStudentsCreate(t *testing.T) {
// 	// this tests CreateStudent and CreateStudents so creation in mass also works.
// 	students := GenerateTestData()
// 	err := CreateNewStudents(students)
// 	if err != nil {
// 		log.Println(" err : ", err)
// 	}
// 	//fmt.Printf("where are my students: %+v", students)
// 	th.AssertEqual(t, "generated: ", len(students), 10)
// 	for _, student := range students {
// 		fmt.Printf("StudentID: %s, Age: %d,Current Year:%d, Graduation Year: %d", student.StudentID, student.Age, student.CurrentYear, student.GraduationYear)
// 		fmt.Println()
// 	}
// }

// func Test_DataGen(t *testing.T) {
// 	students, professors, classes, reportcards, err := GenerateTestData(30, 1)
// 	if err != nil {
// 		t.Errorf("generation failed: %v", err)
// 	}
// 	err = BatchUploadTestData(students, professors, classes, reportcards, nil)
// 	if err != nil {
// 		t.Errorf("upload failed: %v", err)
// 	}
// 	// fmt.Printf("STUDENTS: %v", students)
// 	// fmt.Println("")
// 	// fmt.Printf("PROFESSORS: %v", professors)
// 	// fmt.Println("")
// 	// fmt.Printf("CLASSES: %v", classes)
// 	// fmt.Println("")
// 	// fmt.Printf("REPORTCARDS: %v", reportcards)
// }
