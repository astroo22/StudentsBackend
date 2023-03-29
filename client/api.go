package client

import "time"

type Student_API struct {
	StudentID      string    `json:"studentid"`
	Name           string    `json:"name"`
	CurrentYear    int       `json:"currentyear"`
	GraduationYear int       `json:"graduationyear"`
	AvgGPA         float64   `json:"avggpa,omitempty"`
	Age            int       `json:"age"`
	Dob            time.Time `json:"dob"`
	Enrolled       bool      `json:"enrolled"`
}

// need one for professors

type ReportCard_API struct {
	StudentID  string   `json:"studentID"`
	Math       float64  `json:"math"`
	Science    float64  `json:"science"`
	English    float64  `json:"english"`
	PhysicalED float64  `json:"physicaled"`
	Lunch      float64  `json:"lunch"`
	ClassList  []string `json:"classlist,omitempty"`
}

type Class_API struct {
	ClassID       string   `json:"classid"`
	TeachingGrade int      `json:"teachinggrade"`
	ProfessorID   string   `json:"professorid,omitempty"`
	Subject       string   `json:"subject"`
	Roster        []string `json:"roster,omitempty"`
	ClassAvg      int      `json:"class_avg,omitempty"`
}
