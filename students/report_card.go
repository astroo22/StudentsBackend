package students

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"students/sqlgeneric"
	"time"

	"github.com/lib/pq"
	"golang.org/x/exp/slices"
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
	Math            float64
	Science         float64
	English         float64
	PhysicalED      float64
	Lunch           float64
	AddClassList    []string
	RemoveClassList []string
}

func CreateReportCard(studentID string) (ReportCard, error) {
	return createReportCard(studentID)
}

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
	_, err = db.Exec(insertStatement, studentID, reportCard.Math, reportCard.Science, reportCard.English, reportCard.PhysicalED, reportCard.Lunch) // pq.Array(reportCard.ClassList))
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
func GetReportCards(studentIDs []string) ([]ReportCard, error) {
	return getReportCards(studentIDs)
}

func getReportCards(studentIDs []string) ([]ReportCard, error) {
	if len(studentIDs) == 0 {
		return nil, fmt.Errorf("no data")
	}
	placeholders := make([]string, 0, len(studentIDs))
	batchVals := make([]interface{}, 0, len(studentIDs))
	for i := 0; i < len(studentIDs); i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		batchVals = append(batchVals, studentIDs[i])
	}

	getStatement := fmt.Sprintf(`SELECT * FROM ReportCards WHERE student_id IN (%s)`, strings.Join(placeholders, ","))
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	ret, err := db.Query(getStatement, batchVals...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	reportCards, err := ScanReportCards(ret)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return reportCards, nil
}
func GetReportCardsOfEnrolled(studentIDs []string) ([]ReportCard, error) {
	return getReportCardsOfEnrolled(studentIDs)
}

func getReportCardsOfEnrolled(studentIDs []string) ([]ReportCard, error) {
	if len(studentIDs) == 0 {
		return nil, fmt.Errorf("no data")
	}
	placeholders := make([]string, 0, len(studentIDs))
	batchVals := make([]interface{}, 0, len(studentIDs))
	for i := 0; i < len(studentIDs); i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		batchVals = append(batchVals, studentIDs[i])
	}

	getStatement := fmt.Sprintf(`SELECT ReportCards.* FROM ReportCards 
		JOIN Students ON ReportCards.student_id = Students.student_id 
		WHERE ReportCards.student_id IN (%s) AND Students.enrolled = true`, strings.Join(placeholders, ","))
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()
	ret, err := db.Query(getStatement, batchVals...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	reportCards, err := ScanReportCards(ret)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return reportCards, nil
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
	if len(opts.RemoveClassList) > 0 || len(opts.AddClassList) > 0 {
		classList, err := opts.prepClassListUpdate()
		if err != nil {
			return err
		}
		SQL += fmt.Sprintf(" class_list = $%d", i)
		values = append(values, pq.Array(classList))
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
func (opts UpdateReportCardOptions) prepClassListUpdate() ([]string, error) {
	reportCard, err := getReportCard(opts.StudentID)
	if err != nil {
		return nil, err
	}
	var classList []string
	for _, v := range reportCard.ClassList {
		if !slices.Contains(opts.RemoveClassList, v) {
			classList = append(classList, v)
		}
	}
	classList = append(classList, opts.AddClassList...)
	if len(classList) > 7 {
		return nil, fmt.Errorf("too many classes %v", classList)
	}
	return classList, err
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

func ScanReportCard(row *sql.Row) (ReportCard, error) {
	var (
		reportCard = ReportCard{}
		classList  []byte
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
	temp := removeBrackets(strings.Split(string(classList), ","))
	if len(temp) > 0 {
		reportCard.ClassList = temp
	}
	return reportCard, nil
}
func ScanReportCards(rows *sql.Rows) ([]ReportCard, error) {
	defer rows.Close()
	var (
		reportCards []ReportCard
		classList   []byte
	)
	for rows.Next() {
		reportCard := ReportCard{}
		err := rows.Scan(
			&reportCard.StudentID,
			&reportCard.Math,
			&reportCard.Science,
			&reportCard.English,
			&reportCard.PhysicalED,
			&reportCard.Lunch,
			&classList,
		)
		if err != nil {
			return nil, err
		}
		temp := removeBrackets(strings.Split(string(classList), ","))
		if len(temp) > 0 {
			reportCard.ClassList = temp
		}
		reportCards = append(reportCards, reportCard)
	}
	return reportCards, nil
}
