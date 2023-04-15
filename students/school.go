package students

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"students/sqlgeneric"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type School struct {
	SchoolID      string
	OwnerID       string
	SchoolName    string
	ProfessorList []string
	ClassList     []string
	StudentList   []string
}

type UpdateSchoolOptions struct {
	// will not update
	SchoolID string

	// will update
	SchoolName string

	AddToProfessorList      []string
	RemoveFromProfessorList []string

	AddToClassList      []string
	RemoveFromClassList []string

	AddToStudentList      []string
	RemoveFromStudentList []string
}

func CreateSchool(name, ownerID string, professorList []string, classList []string, studentList []string) (School, error) {
	return createSchool(name, ownerID, professorList, classList, studentList)
}
func createSchool(name, ownerID string, professorList []string, classList []string, studentList []string) (School, error) {

	schoolID := uuid.New().String()
	insertStatement := `INSERT INTO Schools("school_id","owner_id","school_name","professor_list","class_list","student_list") VALUES($1,$2,$3,$4,$5,$6)`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
		return School{}, err
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, schoolID, ownerID, name, pq.Array(professorList), pq.Array(classList), pq.Array(studentList))
	if err != nil {
		return School{}, err
	}
	ret := School{
		SchoolID:      schoolID,
		OwnerID:       ownerID,
		SchoolName:    name,
		ProfessorList: professorList,
		ClassList:     classList,
		StudentList:   studentList,
	}
	return ret, nil
}

func GetSchool(schoolID string) (School, error) {
	return getSchool(schoolID)
}
func getSchool(schoolID string) (School, error) {
	getStatement := `SELECT * FROM Schools WHERE school_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Printf(" err : %v", err)
	}
	defer db.Close()
	school, err := ScanSchool(db.QueryRow(getStatement, schoolID))
	if err != nil {
		return School{}, err
	}
	return school, err
}

// UPDATE
func (opts UpdateSchoolOptions) UpdateSchool() error {
	return opts.updateSchool()
}
func (opts UpdateSchoolOptions) updateSchool() error {
	var (
		SQL    = `UPDATE School SET`
		values []interface{}
		i      = 2
	)
	values = append(values, opts.SchoolID)
	school, err := opts.UpdateHelper()
	if err != nil {
		return err
	}

	if len(opts.SchoolName) > 0 {
		SQL += fmt.Sprintf(" school_name = $%d,", i)
		values = append(values, opts.SchoolName)
		i++
	}
	if len(opts.AddToProfessorList) != 0 || len(opts.RemoveFromProfessorList) != 0 {
		SQL += fmt.Sprintf(" professor_list = $%d,", i)
		values = append(values, pq.Array(school.ProfessorList))
		i++
	}
	if len(opts.AddToClassList) != 0 || len(opts.RemoveFromClassList) != 0 {
		SQL += fmt.Sprintf(" class_list = $%d,", i)
		values = append(values, pq.Array(school.ClassList))
		i++
	}
	if len(opts.AddToStudentList) != 0 || len(opts.RemoveFromStudentList) != 0 {
		SQL += fmt.Sprintf(" student_list = $%d", i)
		values = append(values, pq.Array(school.StudentList))
		i++
	}
	if SQL[len(SQL)-1] == ',' {
		SQL = SQL[:len(SQL)-1]
	}
	SQL += " WHERE school_id = $1"
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	fmt.Println(SQL)
	defer db.Close()
	_, err = db.Exec(SQL, values...)
	if err != nil {
		return err
	}
	return nil

}

func (opts UpdateSchoolOptions) UpdateHelper() (School, error) {
	school, err := GetSchool(opts.SchoolID)
	if err != nil {
		return School{}, err
	}

	// professorlist
	if len(opts.AddToProfessorList) > 0 {
		for _, professor := range opts.AddToProfessorList {
			if !contains(school.ProfessorList, professor) {
				school.ProfessorList = append(school.ProfessorList, professor)
			}
		}
	}
	if len(opts.RemoveFromProfessorList) > 0 {
		for _, professor := range opts.RemoveFromProfessorList {
			school.ProfessorList = remove(school.ProfessorList, professor)
		}
	}
	// classlist
	if len(opts.AddToClassList) > 0 {
		for _, class := range opts.AddToClassList {
			if !contains(school.ClassList, class) {
				school.ClassList = append(school.ClassList, class)
			}
		}
	}
	if len(opts.RemoveFromClassList) > 0 {
		for _, class := range opts.RemoveFromClassList {
			school.ClassList = remove(school.ClassList, class)
		}
	}
	// studentlist
	if len(opts.AddToStudentList) > 0 {
		for _, student := range opts.AddToStudentList {
			if !contains(school.StudentList, student) {
				school.StudentList = append(school.StudentList, student)
			}
		}
	}
	if len(opts.RemoveFromStudentList) > 0 {
		for _, student := range opts.RemoveFromStudentList {
			school.StudentList = remove(school.StudentList, student)
		}
	}
	return school, nil
}

// Utility function to check if a string slice contains a given string.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Utility function to remove a string from a string slice.
func remove(slice []string, str string) []string {
	for i, s := range slice {
		if s == str {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// DELETE SCHOOL
func DeleteSchool(schoolID string) error {
	return deleteSchool(schoolID)
}
func deleteSchool(schoolID string) error {
	deleteStatement := `DELETE FROM Schools WHERE school_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(deleteStatement, schoolID)
	if err != nil {
		return err
	}
	return nil
}

// SCANS
func ScanSchool(row *sql.Row) (School, error) {
	return scanSchool(row)
}
func scanSchool(row *sql.Row) (School, error) {
	var (
		school        = School{}
		professorList sql.NullString
		classList     sql.NullString
		StudentList   sql.NullString
	)
	err := row.Scan(
		&school.SchoolID,
		&school.OwnerID,
		&school.SchoolName,
		&professorList,
		&classList,
		&StudentList,
	)
	if err != nil {
		return School{}, err
	}
	if professorList.Valid {
		school.ProfessorList = strings.Split(professorList.String, ",")
	}
	if classList.Valid {
		school.ClassList = strings.Split(classList.String, ",")
	}
	if StudentList.Valid {
		school.StudentList = strings.Split(StudentList.String, ",")
	}
	return school, nil

}
