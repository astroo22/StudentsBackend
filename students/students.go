package students

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"students/sqlgeneric"
	"time"

	"github.com/google/uuid"
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
type StudentUpdateOptions struct {
	// will not update
	StudentID string
	// will update
	Name           string
	CurrentYear    int
	GraduationYear int
	AvgGPA         float64
	Age            int
	Enrolled       bool
}

func Hello() string {
	return "Hello, world."
}

func CreateNewStudent(name string, currentYear int, graduationYear int, avgGPA float64, age int, dob time.Time, enrolled bool) (string, error) {
	return createNewStudent(name, currentYear, graduationYear, avgGPA, age, dob, enrolled)
}

// maybe return entire student later?
func createNewStudent(name string, currentYear int, graduationYear int, avgGPA float64, age int, dob time.Time, enrolled bool) (string, error) {
	if avgGPA == 0.0 {
		rand.Seed(time.Now().UnixNano())
		randomFloat := rand.Float64() * 4.0
		avgGPA = randomFloat
	}
	insertStatement := `INSERT INTO STUDENTS("student_id","name","current_year","graduation_year","avg_gpa","age","dob","enrolled") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	//dobStr := dob.Format("2006-01-02")
	//fmt.Println(dobStr)
	studentID := uuid.New().String()
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, studentID, name, currentYear, graduationYear, avgGPA, age, dob, enrolled)
	if err != nil {
		return "", err
	}
	return studentID, nil
}

// CreateNewStudents uses batch processing to commit all in one db hit to be more performant
func CreateNewStudents(students []Student) error {
	return createNewStudents(students)
}

func createNewStudents(students []Student) error {
	batch := make([]string, 0, len(students))
	batchVals := make([]interface{}, 0, len(students)*8)
	for n, student := range students {
		batch = append(batch, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", n*8+1, n*8+2, n*8+3, n*8+4, n*8+5, n*8+6, n*8+7, n*8+8))
		batchVals = append(batchVals, student.StudentID)
		batchVals = append(batchVals, student.Name)
		batchVals = append(batchVals, student.CurrentYear)
		batchVals = append(batchVals, student.GraduationYear)
		batchVals = append(batchVals, student.AvgGPA)
		batchVals = append(batchVals, student.Age)
		batchVals = append(batchVals, student.Dob)
		batchVals = append(batchVals, student.Enrolled)
	}
	insertStatement := fmt.Sprintf(`INSERT into STUDENTS("student_id","name","current_year","graduation_year","avg_gpa","age","dob","enrolled") values %s`, strings.Join(batch, ","))
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, batchVals...)
	if err != nil {
		return err
	}
	return nil
}

func GetStudent(studentID string) (Student, error) {
	return getStudent(studentID)
}

func getStudent(studentID string) (Student, error) {
	getStatement := `SELECT * FROM STUDENTS WHERE student_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ret, err := ScanStudent(db.QueryRow(getStatement, studentID))
	if err != nil {
		return ret, err
	}
	return ret, nil
}
func GetAllStudents() ([]Student, error) {
	getStatement := `SELECT * FROM STUDENTS`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Fatal(err)
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
func UpdateStudent(opts StudentUpdateOptions) error {
	return updateStudent(opts)
}

// Should modify to check for sql no rows on a get here
func updateStudent(opts StudentUpdateOptions) error {
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
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(SQL, values...)
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
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(deleteStatement, studentID)
	if err != nil {
		return err
	}
	return nil
}

func ScanStudent(row *sql.Row) (Student, error) {
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
			return students, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return students, err
	}
	return students, nil
}

func GenerateTestData() []Student {
	rand.Seed(time.Now().UnixNano())
	students := make([]Student, 10)

	for i := 0; i < 10; i++ {
		name := "Student " + fmt.Sprint(i)
		currentYear := rand.Intn(12)
		graduationYear := time.Now().Year() + (12 - currentYear)

		avgGPA := float64(rand.Intn(400)) / 100.0

		age := rand.Intn(10) + 15

		// Generate a random date of birth
		dob := time.Date(rand.Intn(10)+1990, time.Month(rand.Intn(12)+1), rand.Intn(28)+1, 0, 0, 0, 0, time.UTC)
		enrolled := rand.Intn(2) == 1

		student := Student{
			StudentID:      uuid.New().String(),
			Name:           name,
			CurrentYear:    currentYear,
			GraduationYear: graduationYear,
			AvgGPA:         avgGPA,
			Age:            age,
			Dob:            dob,
			Enrolled:       enrolled,
		}

		// Add the new student record to the array
		students[i] = student
	}

	return students
}
