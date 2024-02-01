package students

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"students/sqlgeneric"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Student struct {
	StudentID      string
	Name           string
	CurrentYear    int
	GraduationYear int
	AvgGPA         float64
	Age            int
	Dob            time.Time
	Enrolled       bool
}
type UpdateStudentOptions struct {
	// will not update
	StudentID string
	// will update
	CurrentYear    int
	GraduationYear int
	AvgGPA         float64
	Age            int
	Enrolled       bool
}

func CreateNewStudent(name string, currentYear int, graduationYear int, avgGPA float64, age int, dob time.Time, enrolled bool) (Student, error) {
	return createNewStudent(name, currentYear, graduationYear, avgGPA, age, dob, enrolled)
}

func createNewStudent(name string, currentYear int, graduationYear int, avgGPA float64, age int, dob time.Time, enrolled bool) (Student, error) {
	insertStatement := `INSERT INTO STUDENTS("student_id","name","current_year","graduation_year","avg_gpa","age","dob","enrolled") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	studentID := uuid.New().String()
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
		return Student{}, err
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, studentID, name, currentYear, graduationYear, avgGPA, age, dob, enrolled)
	if err != nil {
		return Student{}, err
	}
	ret := Student{
		StudentID:      studentID,
		Name:           name,
		CurrentYear:    currentYear,
		GraduationYear: graduationYear,
		AvgGPA:         avgGPA,
		Age:            age,
		Enrolled:       enrolled,
	}
	return ret, nil
}

func GetStudent(studentID string) (Student, error) {
	return getStudent(studentID)
}

func getStudent(studentID string) (Student, error) {
	getStatement := `SELECT "student_id","name","current_year","graduation_year","avg_gpa","age","dob","enrolled" FROM STUDENTS WHERE student_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
		return Student{}, err
	}
	defer db.Close()
	ret, err := ScanStudent(db.QueryRow(getStatement, studentID))
	if err != nil {
		return ret, err
	}
	return ret, nil
}
func GetAllStudents() ([]Student, error) {
	return getAllStudents()
}

func getAllStudents() ([]Student, error) {
	getStatement := `SELECT * FROM STUDENTS`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()
	ret, err := db.Query(getStatement)
	if err != nil {
		return nil, err
	}
	students, err := ScanStudents(ret)
	if err != nil {
		return nil, err
	}
	return students, nil
}

// for the sake of simplicity this will b a stomp
func (opts UpdateStudentOptions) UpdateStudent() error {
	return opts.updateStudent()
}

// Should modify to check for sql no rows on a get here
func (opts UpdateStudentOptions) updateStudent() error {
	var (
		SQL    = `UPDATE STUDENTS SET`
		values []interface{}
		i      = 2
	)
	values = append(values, opts.StudentID)
	if opts.CurrentYear != 0 {
		SQL += fmt.Sprintf(" current_year = $%d,", i)
		values = append(values, opts.CurrentYear)
		i++
	}
	if opts.GraduationYear != 0 {
		SQL += fmt.Sprintf(" graduation_year = $%d,", i)
		values = append(values, opts.GraduationYear)
		i++
	}
	if opts.AvgGPA != 0 {
		SQL += fmt.Sprintf(" avg_gpa = $%d,", i)
		values = append(values, opts.AvgGPA)
		i++
	}
	if opts.Age != 0 {
		SQL += fmt.Sprintf(" age = $%d,", i)
		values = append(values, opts.Age)
		i++
	}
	SQL += fmt.Sprintf(" enrolled = $%d WHERE student_id = $1", i)
	values = append(values, opts.Enrolled)
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(SQL, values...)
	if err != nil {
		return err
	}
	return nil
}
func FlipEnrollmentStatus(studentIDs []string) error {
	return flipEnrollmentStatus(studentIDs)
}

func flipEnrollmentStatus(studentIDs []string) error {
	updateStatement := `UPDATE STUDENTS SET "enrolled" = NOT "enrolled" WHERE "student_id" = ANY($1)`
	fmt.Println(studentIDs)
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
		return err
	}
	defer db.Close()
	_, err = db.Exec(updateStatement, pq.Array(studentIDs))
	if err != nil {
		return err
	}

	return nil
}

func DeleteStudent(studentID string) error {
	return deleteStudent(studentID)
}
func deleteStudent(studentID string) error {
	deleteStatement := `DELETE FROM STUDENTS WHERE student_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
		return err
	}
	defer db.Close()
	_, err = db.Exec(deleteStatement, studentID)
	if err != nil {
		return err
	}
	return nil
}

func ScanStudent(row *sql.Row) (Student, error) {
	return scanStudent(row)
}

func scanStudent(row *sql.Row) (Student, error) {
	student := Student{}
	err := row.Scan(
		&student.StudentID,
		&student.Name,
		&student.CurrentYear,
		&student.GraduationYear,
		&student.AvgGPA,
		&student.Age,
		&student.Dob,
		&student.Enrolled,
	)
	if err != nil {
		return student, err
	}
	return student, nil
}
func ScanStudents(rows *sql.Rows) ([]Student, error) {
	return scanStudents(rows)
}

func scanStudents(rows *sql.Rows) ([]Student, error) {
	defer rows.Close()
	var students []Student
	for rows.Next() {
		student := Student{}
		err := rows.Scan(
			&student.StudentID,
			&student.Name,
			&student.CurrentYear,
			&student.GraduationYear,
			&student.AvgGPA,
			&student.Age,
			&student.Dob,
			&student.Enrolled,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return students, err
	}
	return students, nil
}

func removeBrackets(roster []string) []string {
	var (
		ret []string
	)
	for _, v := range roster {
		removed := strings.ReplaceAll(strings.ReplaceAll(v, "{", ""), "}", "")
		if removed != "" {
			ret = append(ret, removed)
		}
	}

	return ret
}
