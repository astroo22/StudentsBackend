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

type Class struct {
	ClassID       string
	TeachingGrade int
	ProfessorID   string
	Subject       string
	Roster        []string
	ClassAvg      float64
}
type UpdateClassOptions struct {
	// will not update
	ClassID string
	// will update
	ProfessorID string
	Roster      []string
	ClassAvg    float64
}

func CreateClass(teachingGrade int, professorID string, subject string, roster []string) (Class, error) {
	class := Class{
		ClassID:       uuid.New().String(),
		TeachingGrade: teachingGrade,
		ProfessorID:   professorID,
		Subject:       subject,
		Roster:        roster,
	}
	insertStatement := `INSERT INTO Classes("class_id","teaching_grade","professor_id","subject","roster") VALUES ($1,$2,$3,$4,$5)`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, class.ClassID, class.TeachingGrade, class.ProfessorID, class.Subject, pq.Array(class.Roster))
	if err != nil {
		return Class{}, err
	}
	return class, nil
}
func GetClass(classID string) (Class, error) {
	return getClass(classID)
}

func getClass(classID string) (Class, error) {
	getStatement := `SELECT * FROM Classes WHERE class_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	class, err := ScanClass(db.QueryRow(getStatement, classID))
	if err != nil {
		return Class{}, err
	}
	return class, nil

}
func (opts UpdateClassOptions) UpdateClass() error {
	return opts.updateClass()
}
func (opts UpdateClassOptions) updateClass() error {
	var (
		SQL    = `UPDATE Classes SET`
		values []interface{}
		i      = 2
	)
	values = append(values, opts.ClassID)
	if len(opts.ProfessorID) != 0 {
		SQL += fmt.Sprintf(" professor_id = $%d,", i)
		values = append(values, opts.ProfessorID)
		i++
	}
	if len(opts.Roster) != 0 {
		SQL += fmt.Sprintf(" roster = $%d,", i)
		values = append(values, pq.Array(opts.Roster))
		i++
	}
	if opts.ClassAvg != 0 {
		SQL += fmt.Sprintf(" class_avg = $%d,", i)
		values = append(values, opts.ClassAvg)
		i++
	}
	if SQL[len(SQL)-1] == ',' {
		SQL = SQL[:len(SQL)-1]
	}
	SQL += " WHERE class_id = $1"
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
		return err
	}
	defer db.Close()
	_, err = db.Exec(SQL, values...)
	if err != nil {
		return err
	}
	return nil
}

func removeBrackets(roster []string) []string {
	var (
		ret []string
	)
	for _, v := range roster {
		removed := strings.ReplaceAll(strings.ReplaceAll(v, "{", ""), "}", "")
		ret = append(ret, removed)
	}
	return ret
}
func ScanClass(row *sql.Row) (Class, error) {
	return scanClass(row)
}
func scanClass(row *sql.Row) (Class, error) {
	class := Class{}
	var (
		classAvg sql.NullFloat64
		roster   []byte
	)
	err := row.Scan(
		&class.ClassID,
		&class.TeachingGrade,
		&class.ProfessorID,
		&class.Subject,
		&roster,
		&classAvg,
	)
	if err != nil {
		return Class{}, err
	}
	if classAvg.Valid {
		var value interface{}
		value, err = classAvg.Value()
		if err == nil {
			class.ClassAvg = value.(float64)
		}
	}
	// TODO: strange error here not entirely sure why {} show up. review this later to solve
	temp := removeBrackets(strings.Split(string(roster), ","))
	class.Roster = temp
	return class, nil
}
func DeleteClass(classID string) error {
	return deleteClass(classID)
}
func deleteClass(classID string) error {
	deleteStatement := `DELETE FROM Classes WHERE class_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(deleteStatement, classID)
	if err != nil {
		return err
	}
	return nil
}
