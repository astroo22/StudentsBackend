package client

import "time"

// need to create a section of errors
// need to label errors throughout

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

type Professor_API struct {
	ProfessorID string   `json:"professorid"`
	Name        string   `json:"name"`
	StudentAvg  float64  `json:"studentavg"`
	ClassList   []string `json:"classlist,omitempty"`
}

type ReportCard_API struct {
	StudentID  string   `json:"studentid"`
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

type School_API struct {
	SchoolID      string   `json:"schoolid"`
	SchoolOwnerID string   `json:"schoolownerid"`
	SchoolName    string   `json:"schoolname"`
	ProfessorList []string `json:"professorlist,omitempty"`
	ClassList     []string `json:"classlist,omitempty"`
	StudentList   []string `json:"studentlist,omitempty"`
}
