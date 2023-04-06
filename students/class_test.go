package students

import (
	"fmt"
	th "students/testhelpers"
	"testing"
)

// TODO: everything lol
func Test_ClassCrud(t *testing.T) {

	roster := []string{"mittens", "vesemir", "Bob"}
	finalRoster := []string{"mittens", "Jack", "Tigger"}
	createdClass, err := CreateClass(12, "prof123", "Math", roster)
	if err != nil {
		t.Errorf("Unexpected error creating class: %v", err)
	}

	retrievedClass, err := GetClass(createdClass.ClassID)
	if err != nil {
		t.Errorf("Unexpected error retrieving class: %v", err)
	}
	th.AssertEqual(t, "created vs retrieved professor_id", createdClass.ProfessorID, retrievedClass.ProfessorID)
	th.AssertEqual(t, "created vs retrieved subject", createdClass.Subject, retrievedClass.Subject)
	//th.AssertEqual(t, "created class vs retrieved class", createdClass.Roster[0], retrievedClass.Roster[0])

	updateOpts := UpdateClassOptions{
		ClassID:      createdClass.ClassID,
		ProfessorID:  "prof777",
		AddRoster:    []string{"Jack", "Tigger"},
		RemoveRoster: []string{"Bob", "vesemir"},
		ClassAvg:     90.5,
	}
	err = updateOpts.UpdateClass()
	if err != nil {
		t.Errorf("Unexpected error updating class: %v", err)
	}

	// Retrieve the updated class from the database
	updatedClass, err := GetClass(createdClass.ClassID)
	if err != nil {
		t.Errorf("Unexpected error retrieving updated class: %v", err)
	}
	if th.AssertNotEqual(t, "updated vs retrieved class", updatedClass, retrievedClass) {
		th.AssertEqual(t, "avg check", updatedClass.ClassAvg, 90.5)
		if len(updatedClass.Roster) > 0 {
			fmt.Println(updatedClass.Roster)
		}
		if th.AssertEqual(t, "len roster", len(updatedClass.Roster), len(finalRoster)) {
			th.AssertEqual(t, "roster check", updatedClass.Roster, finalRoster)
		}

	} else {
		t.Errorf("class shouldn't be equal update probably failed")
	}
	updateOpts = UpdateClassOptions{
		ClassID:   createdClass.ClassID,
		AddRoster: []string{"Jack1", "Tigger2"},
	}
	updateRoster := append(finalRoster, "Jack1", "Tigger2")

	err = updateOpts.UpdateClass()
	if err != nil {
		t.Errorf("Unexpected error updating class: %v", err)
	}

	updatedClass, err = GetClass(createdClass.ClassID)
	if err != nil {
		t.Errorf("Unexpected error retrieving updated class: %v", err)
	}
	if th.AssertEqual(t, "len roster", len(updatedClass.Roster), len(updateRoster)) {
		th.AssertEqual(t, "roster check", updatedClass.Roster, updateRoster)
	}
	updateOpts = UpdateClassOptions{
		ClassID:      createdClass.ClassID,
		RemoveRoster: []string{"Jack1", "Tigger2"},
	}
	err = updateOpts.UpdateClass()
	if err != nil {
		t.Errorf("Unexpected error updating class: %v", err)
	}
	updatedClass, err = GetClass(createdClass.ClassID)
	if err != nil {
		t.Errorf("Unexpected error retrieving updated class: %v", err)
	}
	fmt.Println(updatedClass.Roster)
	if th.AssertEqual(t, "len roster", len(updatedClass.Roster), len(finalRoster)) {
		th.AssertEqual(t, "roster check", updatedClass.Roster, finalRoster)
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
