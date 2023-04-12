package students

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"students/sqlgeneric"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/exp/slices"
)

type Professor struct {
	ProfessorID string
	Name        string
	StudentAvg  float64
	ClassList   []string
}
type UpdateProfessorOptions struct {
	// will not update
	ProfessorID string
	// will update
	StudentAvg      float64
	AddClassList    []string
	RemoveClassList []string
}

// CREATE
func CreateProfessor(name string) (Professor, error) {
	return createProfessor(name)
}
func createProfessor(name string) (Professor, error) {
	prof := Professor{
		ProfessorID: uuid.New().String(),
		Name:        name,
	}
	insertStatement := `INSERT INTO Professors("professor_id","name") VALUES ($1,$2)`

	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, prof.ProfessorID, prof.Name)
	if err != nil {
		return Professor{}, err
	}
	return prof, nil
}
func GetProfessor(professorID string) (Professor, error) {
	getStatement := `SELECT * FROM Professors WHERE professor_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Printf(" err : %v", err)
	}
	defer db.Close()
	prof, err := ScanProf(db.QueryRow(getStatement, professorID))
	if err != nil {
		return Professor{}, err
	}
	return prof, nil
}

// UPDATE
func (opts UpdateProfessorOptions) UpdateProfessor() error {
	return opts.updateProfessor()
}
func (opts UpdateProfessorOptions) updateProfessor() error {
	var (
		SQL    = `UPDATE Professors SET`
		values []interface{}
		i      = 2
	)
	values = append(values, opts.ProfessorID)

	if opts.StudentAvg != 0 {
		SQL += fmt.Sprintf(" student_avg = $%d,", i)
		values = append(values, opts.StudentAvg)
		i++
		fmt.Println(opts.StudentAvg)
	}
	if len(opts.AddClassList) != 0 || len(opts.RemoveClassList) != 0 {
		classList, err := opts.prepProfessorClassUpdate()
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
	SQL += " WHERE professor_id = $1"
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
func (opts UpdateProfessorOptions) prepProfessorClassUpdate() ([]string, error) {
	var ret []string
	prof, err := GetProfessor(opts.ProfessorID)
	if err != nil {
		return nil, err
	}
	for _, v := range prof.ClassList {
		if !slices.Contains(opts.RemoveClassList, v) {
			ret = append(ret, v)
		}
	}
	ret = append(ret, opts.AddClassList...)
	return ret, nil
}

// DELETE
func DeleteProfessor(professorID string) error {
	return deleteProfessor(professorID)
}

func deleteProfessor(professorID string) error {
	deleteStatement := `DELETE FROM Professors WHERE professor_id =$1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(deleteStatement, professorID)
	if err != nil {
		return err
	}
	return nil
}

// SCANS
func ScanProf(row *sql.Row) (Professor, error) {
	return scanProf(row)
}

func scanProf(row *sql.Row) (Professor, error) {
	var (
		prof      = Professor{}
		stdAvg    sql.NullFloat64
		classList sql.NullString
	)
	err := row.Scan(
		&prof.ProfessorID,
		&prof.Name,
		&stdAvg,
		&classList,
	)
	if err != nil {
		return Professor{}, err
	}
	if stdAvg.Valid {
		var value interface{}
		value, err = stdAvg.Value()
		if err == nil {
			prof.StudentAvg = value.(float64)
		} else {
			fmt.Println(err)
		}
	}
	if classList.Valid {
		prof.ClassList = strings.Split(classList.String, ",")
	}
	return prof, nil
}
