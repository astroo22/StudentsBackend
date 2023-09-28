package students

import (
	"fmt"
	"strings"
	th "students/testhelpers"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_SchoolCrud(t *testing.T) {
	var (
		roster    []string
		classList []string
		profList  []string
		ownerID   = uuid.New().String()
		schoolID  = uuid.New().String()
	)

	//students going to update add rick then remove him
	stu1, err := CreateNewStudent("vesemir", 2, 2034, 3.5, 9, time.Date(2014, time.September, 7, 3, 2, 1, 0, time.UTC), true)
	if err != nil {
		t.Error(err)
	}
	stu2, err := CreateNewStudent("mittens", 2, 2034, 3.5, 9, time.Date(2014, time.September, 6, 3, 2, 1, 0, time.UTC), true)
	if err != nil {
		t.Error(err)
	}
	stu3, err := CreateNewStudent("rick", 2, 1976, 4.0, 9, time.Date(1976, time.September, 5, 3, 2, 1, 0, time.UTC), true)
	if err != nil {
		t.Error(err)
	}
	stu4, err := CreateNewStudent("morty", 2, 2003, 1.7, 9, time.Date(2014, time.September, 4, 3, 2, 1, 0, time.UTC), true)
	if err != nil {
		t.Error(err)
	}
	stu5, err := CreateNewStudent("squidward", 2, 2034, 3.2, 9, time.Date(2014, time.September, 3, 3, 2, 1, 0, time.UTC), true)
	if err != nil {
		t.Error(err)
	}
	//roster
	roster = append(roster, stu1.StudentID, stu2.StudentID, stu4.StudentID, stu5.StudentID)

	//professors guna kill off bob then maybe reincarnate him
	prof1, err := CreateProfessor("pvt pickles")
	if err != nil {
		t.Error(err)
	}
	prof2, err := CreateProfessor("the chosen one")
	if err != nil {
		t.Error(err)
	}
	prof3, err := CreateProfessor("bob")
	if err != nil {
		t.Error(err)
	}
	prof4, err := CreateProfessor("deadpool")
	if err != nil {
		t.Error(err)
	}
	prof5, err := CreateProfessor("the silver surfer")
	if err != nil {
		t.Error(err)
	}
	//proflist
	profList = append(profList, prof1.ProfessorID, prof2.ProfessorID, prof3.ProfessorID, prof4.ProfessorID, prof5.ProfessorID)

	//classes
	class1, err := CreateClass(2, prof1.ProfessorID, "math", roster)
	if err != nil {
		t.Error(err)
	}
	class2, err := CreateClass(2, prof2.ProfessorID, "science", roster)
	if err != nil {
		t.Error(err)
	}
	class3, err := CreateClass(2, prof3.ProfessorID, "english", roster)
	if err != nil {
		t.Error(err)
	}
	class4, err := CreateClass(2, prof4.ProfessorID, "physical_ed", roster)
	if err != nil {
		t.Error(err)
	}
	class5, err := CreateClass(2, prof5.ProfessorID, "lunch", roster)
	if err != nil {
		t.Error(err)
	}
	//classlist
	classList = append(classList, class1.ClassID, class2.ClassID, class4.ClassID, class5.ClassID)
	school, err := CreateSchool(schoolID, "PLUS ULTRA ACADEMY test1", ownerID, profList, classList, roster)
	if err != nil {
		t.Error(err)
	}
	// check get
	schoolGet, err := GetSchool(school.SchoolID)
	if err != nil {
		t.Error(err)
	} else {
		if th.AssertEqual(t, "classlist len", len(schoolGet.ClassList), 4) {
			th.AssertEqual(t, "classlist", schoolGet.ClassList, classList)
		}
		if th.AssertEqual(t, "proflist len", len(schoolGet.ProfessorList), 5) {
			th.AssertEqual(t, "proflist", schoolGet.ProfessorList, profList)
		}
		if th.AssertEqual(t, "roster len", len(schoolGet.StudentList), 4) {
			th.AssertEqual(t, "roster", schoolGet.StudentList, roster)
		}
	}
	opts := UpdateSchoolOptions{
		SchoolID: school.SchoolID,
		// adding rick
		AddToStudentList: strings.Split(stu3.StudentID, " "),
		// killing bob
		RemoveFromProfessorList: strings.Split(prof3.ProfessorID, " "),
		// adding to classlist
		AddToClassList: strings.Split(class3.ClassID, " "),
	}
	err = opts.UpdateSchool()
	if err != nil {
		t.Error(err)
	}
	//adding rick to list
	roster = append(roster, stu3.StudentID)
	fmt.Println(roster)
	//removing bob
	profList = remove(profList, prof3.ProfessorID)
	fmt.Println(profList)
	// adding class
	classList = append(classList, class3.ClassID)
	fmt.Println(classList)

	//check if rick exist and if bob is alive
	schoolGet, err = GetSchool(school.SchoolID)
	if err != nil {
		t.Error(err)
	} else {

		if th.AssertEqual(t, "StudentList len", len(schoolGet.StudentList), 5) {
			th.AssertEqual(t, "StudentList", schoolGet.StudentList, roster)
		}
		if th.AssertEqual(t, "proflist len", len(schoolGet.ProfessorList), 4) {
			th.AssertEqual(t, "proflist", schoolGet.ProfessorList, profList)
		}
		if th.AssertEqual(t, "classlist len", len(schoolGet.ClassList), 5) {
			th.AssertEqual(t, "classlist", schoolGet.ClassList, classList)
		}
	}
	opts = UpdateSchoolOptions{
		SchoolID: school.SchoolID,
		// rick got bored
		RemoveFromStudentList: strings.Split(stu3.StudentID, " "),
		// revive bob
		AddToProfessorList: strings.Split(prof3.ProfessorID, " "),
		// remove boring class
		RemoveFromClassList: strings.Split(class3.ClassID, " "),
	}
	err = opts.UpdateSchool()
	if err != nil {
		t.Error(err)
	}
	//adding rick to list
	roster = remove(roster, stu3.StudentID)

	//removing bob
	profList = append(profList, prof3.ProfessorID)

	// adding class
	classList = remove(classList, class3.ClassID)

	schoolGet, err = GetSchool(school.SchoolID)
	if err != nil {
		t.Error(err)
	} else {

		if th.AssertEqual(t, "classlist len", len(schoolGet.ClassList), 4) {
			th.AssertEqual(t, "classlist", schoolGet.ClassList, classList)
		}
		if th.AssertEqual(t, "proflist len", len(schoolGet.ProfessorList), 5) {
			th.AssertEqual(t, "proflist", schoolGet.ProfessorList, profList)
		}
		if th.AssertEqual(t, "roster len", len(schoolGet.StudentList), 4) {
			th.AssertEqual(t, "roster", schoolGet.StudentList, roster)
		}
	}
	err = DeleteSchool(school.SchoolID)
	if err != nil {
		t.Error(err)
	}
	schoolGet, err = GetSchool(school.SchoolID)
	if err == nil {
		t.Error("get failed after delete")
	}

}
