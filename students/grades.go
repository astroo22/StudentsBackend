package students

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"students/sqlgeneric"
	"time"
)

type ReportCard struct {
	StudentID  string
	Math       float64
	Science    float64
	English    float64
	PhysicalED float64
	Lunch      float64
	ClassList  []string
}
type UpdateReportCardOptions struct {
	// will not update
	StudentID string
	// will update
	Math       float64
	Science    float64
	English    float64
	PhysicalED float64
	Lunch      float64
	ClassList  []string
}

func CreateReportCard(studentID string) (ReportCard, error) {
	return createReportCard(studentID)
}

// TODO: Create a Batch Create version of this to be used with generate test data
func createReportCard(studentID string) (ReportCard, error) {
	rand.Seed(time.Now().UnixNano())
	reportCard := ReportCard{
		StudentID:  studentID,
		Math:       float64(rand.Intn(400)) / 100.0,
		Science:    float64(rand.Intn(400)) / 100.0,
		English:    float64(rand.Intn(400)) / 100.0,
		PhysicalED: float64(rand.Intn(400)) / 100.0,
		Lunch:      float64(rand.Intn(400)) / 100.0,
	}
	insertStatement := `INSERT INTO ReportCards("student_id","math","science","english","physical_ed","lunch") values ($1,$2,$3,$4,$5,$6)`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, studentID, reportCard.Math, reportCard.Science, reportCard.English, reportCard.PhysicalED, reportCard.Lunch)
	if err != nil {
		return ReportCard{}, err
	}
	return reportCard, nil
}

func GetReportCard(studentID string) (ReportCard, error) {
	return getReportCard(studentID)
}
func getReportCard(studentID string) (ReportCard, error) {
	getStatement := `SELECT * FROM ReportCards WHERE student_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	reportCard, err := ScanReportCard(db.QueryRow(getStatement, studentID))
	if err != nil {
		return ReportCard{}, err
	}
	return reportCard, nil
}
func (opts UpdateReportCardOptions) UpdateReportCard() error {
	return opts.updateReportCard()
}

func (opts UpdateReportCardOptions) updateReportCard() error {
	var (
		SQL    = `UPDATE ReportCards SET`
		values []interface{}
		i      = 2
	)
	values = append(values, opts.StudentID)
	if opts.Math != 0 {
		SQL += fmt.Sprintf(" math = $%d,", i)
		values = append(values, opts.Math)
		i++
	}
	if opts.Science != 0 {
		SQL += fmt.Sprintf(" science = $%d,", i)
		values = append(values, opts.Science)
		i++
	}
	if opts.English != 0 {
		SQL += fmt.Sprintf(" english = $%d,", i)
		values = append(values, opts.English)
		i++
	}
	if opts.PhysicalED != 0 {
		SQL += fmt.Sprintf(" physical_ed = $%d,", i)
		values = append(values, opts.PhysicalED)
		i++
	}
	if opts.Lunch != 0 {
		SQL += fmt.Sprintf(" lunch = $%d,", i)
		values = append(values, opts.Lunch)
		i++
	}
	if len(opts.ClassList) != 0 {
		SQL += fmt.Sprintf(" class_list = $%d", i)
		values = append(values, strings.Join(opts.ClassList, ","))
		i++
	}
	if SQL[len(SQL)-1] == ',' {
		SQL = SQL[:len(SQL)-1]
	}
	SQL += " WHERE student_id = $1"
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println("update reportcards : ", err)
	}
	defer db.Close()
	_, err = db.Exec(SQL, values...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteReportCard(studentID string) error {
	return deleteReportCard(studentID)
}
func deleteReportCard(studentID string) error {
	SQL := `DELETE FROM ReportCards WHERE student_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(SQL, studentID)
	if err != nil {
		return err
	}
	return nil
}
func DeleteBatchReportCard(students []Student) error {
	return deleteBatchReportCard(students)
}
func deleteBatchReportCard(students []Student) error {
	batch := make([]string, 0, len(students))
	batchVals := make([]interface{}, 0, len(students))
	for n, student := range students {
		batch = append(batch, fmt.Sprintf("$%d", n+1))
		batchVals = append(batchVals, student.StudentID)
	}
	//batch = append(batch, ")")
	SQL := fmt.Sprintf(`DELETE FROM ReportCards WHERE student_id IN (%s)`, strings.Join(batch, ","))
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(SQL, batchVals...)
	if err != nil {
		return err
	}
	return nil
}

func ScanReportCard(row *sql.Row) (ReportCard, error) {
	var (
		reportCard = ReportCard{}
		classList  sql.NullString
	)
	err := row.Scan(
		&reportCard.StudentID,
		&reportCard.Math,
		&reportCard.Science,
		&reportCard.English,
		&reportCard.PhysicalED,
		&reportCard.Lunch,
		&classList,
	)
	if err != nil {
		return reportCard, err
	}
	if classList.Valid {
		reportCard.ClassList = strings.Split(classList.String, ",")
	}
	return reportCard, nil
}
