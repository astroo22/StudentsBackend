package client

import "time"

type Student struct {
	StudentID      int       `json:"studentID"`
	Name           string    `json:"name"`
	CurrentYear    int       `json:"currentYear"`
	GraduationYear int       `json:"graduationYear"`
	AvgGPA         int       `json:"avgGPA,omitempty"`
	Age            int       `json:"age"`
	Dob            time.Time `json:"dob"`
	Enrolled       bool      `json:"enrolled"`
}
type ClassRoster struct {
	ID             int       `json:"ID"`
	GraduationYear int       `json:"classGraduationYear"`
	ClassRoster    []Student `json:"classRoster,omitempty"`
}
type StudentGrades struct {
	StudentID  int `json:"studentID"`
	Math       int `json:"math"`
	Science    int `json:"science"`
	English    int `json:"english"`
	PhysicalED int `json:"physicalED"`
	Lunch      int `json:"lunch"`
}

type Class struct {
	ClassID         int       `json:"classID"`
	ProfessorName   string    `json:"professorName,omitempty"`
	AvgClassGPA     int       `json:"avgClassGPA,omitempty"`
	CurrentlyActive bool      `json:"currentlyActive"`
	StudentRoster   []Student `json:"studentRoster,omitempty"`
}
