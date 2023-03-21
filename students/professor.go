package students

import (
	"database/sql"
	"fmt"
	"log"
	"students/sqlgeneric"

	"github.com/google/uuid"
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
	StudentAvg float64
	ClassList  []string
}

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
		log.Fatal(err)
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
		log.Fatal(err)
	}
	defer db.Close()
	prof, err := ScanProf(db.QueryRow(getStatement, professorID))
	if err != nil {
		return Professor{}, err
	}
	return prof, nil
}

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
	}
	if len(opts.ClassList) != 0 {
		SQL += fmt.Sprintf(" class_list = $%d", i)
		values = append(values, opts.ClassList)
		i++
	}
	if SQL[len(SQL)-1] == ',' {
		SQL = SQL[:len(SQL)-1]
	}
	SQL += " WHERE professor_id = $1"
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

func DeleteProfessor(professorID string) error {
	deleteStatement := `DELETE * FROM Professors WHERE professor_id =$1`
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(deleteStatement, professorID)
	if err != nil {
		return err
	}
	return nil
}

func ScanProf(row *sql.Row) (Professor, error) {
	return scanProf(row)
}

// this is going to error on nulls will fix when I write testing. This will also need to be fixed in class and grades probably
func scanProf(row *sql.Row) (Professor, error) {
	prof := Professor{}
	err := row.Scan(
		&prof.ProfessorID,
		&prof.Name,
		&prof.ClassList,
		&prof.StudentAvg,
	)
	if err != nil {
		return Professor{}, err
	}
	return prof, nil
}
