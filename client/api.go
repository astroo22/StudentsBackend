package client

import "time"

type Student struct {
	StudentID      string    `json:"student_id"`
	Name           string    `json:"name"`
	CurrentYear    int       `json:"current_year"`
	GraduationYear int       `json:"graduation_year"`
	AvgGPA         int       `json:"gpa,omitempty"`
	Age            int       `json:"age"`
	Dob            time.Time `json:"dob"`
	Enrolled       bool      `json:"enrolled"`
}

// TODO: Also redo this for grades
// type ClassRoster struct {
// 	ID             int       `json:"ID"`
// 	GraduationYear int       `json:"classGraduationYear"`
// 	ClassRoster    []Student `json:"classRoster,omitempty"`
// }

// TODO: fix this
type StudentGrades struct {
	StudentID  int `json:"studentID"`
	Math       int `json:"math"`
	Science    int `json:"science"`
	English    int `json:"english"`
	PhysicalED int `json:"physicalED"`
	Lunch      int `json:"lunch"`
}

type Class struct {
	ClassID       string   `json:"class_id"`
	TeachingGrade int      `json:"teaching_grade"`
	ProfessorID   string   `json:"professor_id,omitempty"`
	Subject       string   `json:"subject"`
	Roster        []string `json:"roster,omitempty"`
	ClassAvg      int      `json:"class_avg,omitempty"`
}
