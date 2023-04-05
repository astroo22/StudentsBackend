package students

import (
	th "students/testhelpers"
	"testing"
)

func Test_ProfessorCrud(t *testing.T) {
	var (
		newAvg = 3.5
	)
	createdProf, err := CreateProfessor("chapstick")
	if err != nil {
		t.Errorf("Unexpected error creating professor: %v", err)
	}

	retrievedProf, err := GetProfessor(createdProf.ProfessorID)
	if err != nil {
		t.Errorf("Unexpected error retrieving professor: %v", err)
	}

	th.AssertEqual(t, "created professor vs retrieved professor", createdProf, retrievedProf)
	opts := UpdateProfessorOptions{
		ProfessorID: createdProf.ProfessorID,
		StudentAvg:  newAvg,
	}
	err = opts.UpdateProfessor()
	if err != nil {
		t.Errorf("Unexpected error updating class: %v", err)
	}

	// Retrieve the updated class from the database
	updatedProf, err := GetProfessor(createdProf.ProfessorID)
	if err != nil {
		t.Errorf("Unexpected error retrieving updated class: %v", err)
	}

	if th.AssertNotEqual(t, "updated class vs retrieved class", updatedProf, retrievedProf) {
		th.AssertEqual(t, "std avg comparison", updatedProf.StudentAvg, newAvg)
	}
}
