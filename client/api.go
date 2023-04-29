package client

import (
	"students/students"
	"time"
)

// The following code is to seperate the API layer from the internal structure which in a business
// setting allows more controll over the front facing api structure. Therefore this would often exist so I've included it.
// However for the purposes of this application it is not needed.
// I've included it to display my familiarity with the concepts and its basic implementation.

// TODO: need to create a section of errors
// TODO: distribute said errors

type Student_API struct {
	StudentID      string    `json:"student_id"`
	Name           string    `json:"name"`
	CurrentYear    int       `json:"current_year"`
	GraduationYear int       `json:"graduation_year"`
	AvgGPA         float64   `json:"avg_gpa,omitempty"`
	Age            int       `json:"age"`
	Dob            time.Time `json:"dob"`
	Enrolled       bool      `json:"enrolled"`
}

type Professor_API struct {
	ProfessorID string   `json:"professor_id"`
	Name        string   `json:"name"`
	StudentAvg  float64  `json:"student_avg"`
	ClassList   []string `json:"class_list,omitempty"`
}

type ReportCard_API struct {
	StudentID  string   `json:"student_id"`
	Math       float64  `json:"math"`
	Science    float64  `json:"science"`
	English    float64  `json:"english"`
	PhysicalED float64  `json:"physical_ed"`
	Lunch      float64  `json:"lunch"`
	ClassList  []string `json:"class_list,omitempty"`
}

type Class_API struct {
	ClassID       string   `json:"class_id"`
	TeachingGrade int      `json:"teaching_grade"`
	ProfessorID   string   `json:"professor_id,omitempty"`
	Subject       string   `json:"subject"`
	Roster        []string `json:"roster,omitempty"`
	ClassAvg      float64  `json:"class_avg,omitempty"`
}

type School_API struct {
	SchoolID      string   `json:"school_id"`
	SchoolOwnerID string   `json:"school_owner_id"`
	SchoolName    string   `json:"school_name"`
	AvgGPA        float64  `json:"avg_gpa"`
	Ranking       int      `json:"ranking"`
	ProfessorList []string `json:"professor_list,omitempty"`
	ClassList     []string `json:"class_list,omitempty"`
	StudentList   []string `json:"student_list,omitempty"`
}

// unsure if I agree with these structures existance but ordering is needed in front end to be more performant
type GradeAvg_API struct {
	Grade  int     `json:"grade"`
	AvgGPA float64 `json:"avg_gpa"`
}

// TODO: convert some of the handlers to use the API
// converters
// Student to api
func StudentToAPI(stu students.Student) Student_API {
	return Student_API{
		StudentID:      stu.StudentID,
		Name:           stu.Name,
		CurrentYear:    stu.CurrentYear,
		GraduationYear: stu.GraduationYear,
		AvgGPA:         stu.AvgGPA,
		Age:            stu.Age,
		Dob:            stu.Dob,
		Enrolled:       stu.Enrolled,
	}
}

// Students to API
func StudentsToAPI(students []students.Student) []Student_API {
	apiStudents := make([]Student_API, len(students))
	for i, student := range students {
		apiStudents[i] = StudentToAPI(student)
	}
	return apiStudents
}

// Professor to API
func ProfessorToAPI(prof students.Professor) Professor_API {
	return Professor_API{
		ProfessorID: prof.ProfessorID,
		Name:        prof.Name,
		StudentAvg:  prof.StudentAvg,
		ClassList:   prof.ClassList,
	}
}

// Professors to API
func ProfessorsToAPI(professors []students.Professor) []Professor_API {
	apiProfessors := make([]Professor_API, len(professors))
	for i, prof := range professors {
		apiProfessors[i] = ProfessorToAPI(prof)
	}
	return apiProfessors
}

// ReportCard to API
func ReportCardToAPI(rc students.ReportCard) ReportCard_API {
	return ReportCard_API{
		StudentID:  rc.StudentID,
		Math:       rc.Math,
		Science:    rc.Science,
		English:    rc.English,
		PhysicalED: rc.PhysicalED,
		Lunch:      rc.Lunch,
		ClassList:  rc.ClassList,
	}
}

// ReportCards to API
func ReportCardsToAPI(reportCards []students.ReportCard) []ReportCard_API {
	apiReportCards := make([]ReportCard_API, len(reportCards))
	for i, rc := range reportCards {
		apiReportCards[i] = ReportCardToAPI(rc)
	}
	return apiReportCards
}

// Class to API
func ClassToAPI(class students.Class) Class_API {
	return Class_API{
		ClassID:       class.ClassID,
		TeachingGrade: class.TeachingGrade,
		ProfessorID:   class.ProfessorID,
		Subject:       class.Subject,
		Roster:        class.Roster,
		ClassAvg:      class.ClassAvg,
	}
}

// Classes to API
func ClassesToAPI(classes []students.Class) []Class_API {
	apiClasses := make([]Class_API, len(classes))
	for i, class := range classes {
		apiClasses[i] = ClassToAPI(class)
	}
	return apiClasses
}

// School to API
func SchoolToAPI(school students.School) School_API {
	return School_API{
		SchoolID:      school.SchoolID,
		SchoolOwnerID: school.OwnerID,
		SchoolName:    school.SchoolName,
		AvgGPA:        school.AvgGPA,
		Ranking:       school.Ranking,
		ProfessorList: school.ProfessorList,
		ClassList:     school.ClassList,
		StudentList:   school.StudentList,
	}
}
func SchoolsToAPI(schools []students.School) []School_API {
	apiSchools := make([]School_API, len(schools))
	for i, school := range schools {
		apiSchools[i] = SchoolToAPI(school)
	}
	return apiSchools
}
