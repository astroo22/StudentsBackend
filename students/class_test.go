package students

import (
	"fmt"
	th "students/testhelpers"
	"testing"

	"golang.org/x/exp/slices"
)

func Test_ClassCrud(t *testing.T) {
	th.TestingInit()
	defer th.TestingEnvDif()
	studentRoster, err := GenerateStudents(5, 12)
	if err != nil {
		t.Fatalf("Unexpected error creating student roster: %v", err)
	}
	newStudentRoster, err := GenerateStudents(2, 12)
	if err != nil {
		t.Fatalf("Unexpected error creating newStudent roster: %v", err)
	}
	prof1, err := CreateProfessor("PVT PICKLES")
	if err != nil {
		t.Fatalf("Unexpected error creating Professor: %v", err)
	}
	prof2, err := CreateProfessor("RICK")
	if err != nil {
		t.Fatalf("Unexpected error creating Professor: %v", err)
	}
	roster := []string{}
	for _, v := range studentRoster {
		roster = append(roster, v.StudentID)
	}
	newStdRoster := []string{}
	for _, v := range newStudentRoster {
		newStdRoster = append(newStdRoster, v.StudentID)
	}

	class, err := CreateClass(12, prof1.ProfessorID, "math", roster)
	if err != nil {
		t.Errorf("unexpected error creating class: %v", err)
	}
	fmt.Println(class)

	retrievedClass, err := GetClass(class.ClassID)
	if err != nil {
		t.Errorf("Unexpected error retrieving class: %v", err)
	}
	th.AssertEqual(t, "created vs retrieved professor_id", class.ProfessorID, retrievedClass.ProfessorID)
	th.AssertEqual(t, "created vs retrieved subject", class.Subject, retrievedClass.Subject)
	//th.AssertEqual(t, "created class vs retrieved class", createdClass.Roster[0], retrievedClass.Roster[0])

	updateOpts := UpdateClassOptions{
		ClassID:      class.ClassID,
		ProfessorID:  prof2.ProfessorID,
		AddRoster:    newStdRoster,
		RemoveRoster: []string{roster[0], roster[1]},
		ClassAvg:     90.5,
	}
	err = updateOpts.UpdateClass()
	if err != nil {
		t.Errorf("Unexpected error updating class: %v", err)
	}
	roster[0] = newStdRoster[0]
	roster[1] = newStdRoster[1]
	// Retrieve the updated class from the database
	updatedClass, err := GetClass(class.ClassID)
	if err != nil {
		t.Fatalf("Unexpected error retrieving updated class: %v", err)
	}
	if th.AssertNotEqual(t, "updated vs retrieved class", updatedClass, retrievedClass) {
		th.AssertEqual(t, "avg check", updatedClass.ClassAvg, 90.5)
		if len(updatedClass.Roster) > 0 {
			fmt.Println(updatedClass.Roster)
		}
		if th.AssertEqual(t, "len roster", len(updatedClass.Roster), len(roster)) {
			if !slices.Contains(updatedClass.Roster, roster[0]) || !slices.Contains(updatedClass.Roster, roster[1]) {
				t.Error("Roster check failed did not update with new values")
			}
		}

	} else {
		t.Errorf("class shouldn't be equal update probably failed")
	}

	err = DeleteClass(updatedClass.ClassID)
	if err != nil {
		t.Errorf("unable to delete class")
	}
	_, err = GetClass(updatedClass.ClassID)
	if err == nil {
		t.Errorf("get succeeded delete failed")
	}
}
