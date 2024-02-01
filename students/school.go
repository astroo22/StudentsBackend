package students

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"students/sqlgeneric"

	"github.com/lib/pq"
)

type School struct {
	SchoolID      string
	OwnerID       string
	SchoolName    string
	AvgGPA        float64
	Ranking       int
	ProfessorList []string
	ClassList     []string
	StudentList   []string
}

// will update other values as functionality is added
type UpdateSchoolOptions struct {
	// will not update
	SchoolID string

	// will update
	SchoolName string `json:"name"`
	Avg_gpa    float64

	AddToProfessorList      []string
	RemoveFromProfessorList []string

	AddToClassList      []string
	RemoveFromClassList []string

	AddToStudentList      []string
	RemoveFromStudentList []string

	// enrollment
	Enrollment_change_ids []string `json:"enrollment_change_ids"`
}

func CreateSchool(schoolID, name, ownerID string, professorList []string, classList []string, studentList []string) (School, error) {
	return createSchool(schoolID, name, ownerID, professorList, classList, studentList)
}
func createSchool(schoolID, name, ownerID string, professorList []string, classList []string, studentList []string) (School, error) {
	insertStatement := `INSERT INTO Schools("school_id","owner_id","school_name","professor_list","class_list","student_list") VALUES($1,$2,$3,$4,$5,$6)`
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Println(" err : ", err)
		return School{}, err
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, schoolID, ownerID, name, pq.Array(professorList), pq.Array(classList), pq.Array(studentList))
	if err != nil {
		fmt.Println(err)
		return School{}, err
	}

	schoolAvg, err := UpdateSchoolAvg(schoolID)
	if err != nil {
		fmt.Println(err)
		return School{}, err
	}

	err = UpdateSchoolRankings()
	if err != nil {
		return School{}, err
	}
	ret := School{
		SchoolID:      schoolID,
		OwnerID:       ownerID,
		SchoolName:    name,
		AvgGPA:        schoolAvg,
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
	getStatement := `SELECT school_id, school_name, avg_gpa, ranking, professor_list, class_list, student_list FROM Schools WHERE school_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
	}
	defer db.Close()
	school, err := ScanSchool(db.QueryRow(getStatement, schoolID))
	if err != nil {
		return School{}, err
	}
	return school, err
}

func GetAllSchools() ([]School, error) {
	return getAllSchools()
}
func getAllSchools() ([]School, error) {
	getStatement := `SELECT school_id, school_name, avg_gpa, ranking, professor_list,class_list,student_list FROM Schools`
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
	}
	defer db.Close()
	ret, err := db.Query(getStatement)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	schools, err := ScanSchools(ret)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return schools, err
}

// this is currently untested
func GetAllSchoolsForUser(ownerID string) ([]School, error) {
	return getAllSchoolsForUser(ownerID)
}
func getAllSchoolsForUser(ownerID string) ([]School, error) {
	getStatement := `SELECT school_id, school_name, avg_gpa, ranking, professor_list,class_list,student_list FROM Schools WHERE owner_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
		return nil, err
	}
	defer db.Close()
	ret, err := db.Query(getStatement, ownerID)
	if err != nil {
		return nil, err
	}
	schools, err := ScanSchools(ret)
	if err != nil {
		return nil, err
	}
	return schools, err
}

// this function needs to have updated scans before it will work
// func UpdateSchoolAvg(schoolID string) (float64, error) {
// 	var (
// 		schoolAvg float64
// 	)
// 	classes, err := GetClassesForSchool(schoolID)
// 	if err != nil {
// 		return 0, err
// 	}

// 	for _, class := range classes {
// 		schoolAvg += class.ClassAvg
// 	}
// 	schoolAvg = math.Round(schoolAvg/float64(len(classes))*100) / 100
// 	fmt.Println(schoolAvg)
// 	updateQuery := `UPDATE Schools SET avg_gpa = $1 WHERE school_id = $2`
// 	db, err := sqlgeneric.Init()
// 	if err != nil {
// 		log.Printf(" err : %v", err)
// 	}
// 	defer db.Close()
// 	_, err = db.Exec(updateQuery, schoolAvg, schoolID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return schoolAvg, nil

// }
func UpdateSchoolRankings() error {
	return updateSchoolRankings()
}

func updateSchoolRankings() error {
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
		return err
	}
	defer db.Close()

	// Query to update school rankings
	query := `WITH ranked_schools AS (
			SELECT school_id, avg_gpa, DENSE_RANK() OVER (ORDER BY avg_gpa DESC) AS ranking
			FROM Schools
		)
		UPDATE Schools
		SET ranking = ranked_schools.ranking
		FROM ranked_schools
		WHERE Schools.school_id = ranked_schools.school_id
	`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	fmt.Println("School rankings updated successfully")
	return nil
}

func GetClassesForSchool(schoolID string) ([]Class, error) {
	return getClassesForSchool(schoolID)
}
func getClassesForSchool(schoolID string) ([]Class, error) {
	query := `
        SELECT c.*
        FROM schools s, unnest(s.class_list) cl, classes c
        WHERE s.school_id = $1 AND cl = c.class_id
		ORDER BY c.class_avg DESC
    `
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(query, schoolID)
	if err != nil {
		return nil, err
	}
	classes, err := ScanClasses(rows)
	if err != nil {
		return nil, err
	}
	return classes, nil
}

func GetStudentsForSchool(schoolID string) ([]Student, error) {
	return getStudentsForSchool(schoolID)
}
func getStudentsForSchool(schoolID string) ([]Student, error) {
	query := `
		SELECT st.student_id, st.name,st.current_year,st.graduation_year, st.avg_gpa, st.age,st.dob,st.enrolled
		FROM schools s, unnest(s.student_list) sl, students st
		WHERE s.school_id = $1 and sl = st.student_id
		ORDER BY st.avg_gpa DESC
	`
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(query, schoolID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	students, err := ScanStudents(rows)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return students, nil
}

// UPDATE
func (opts UpdateSchoolOptions) UpdateSchool() error {
	return opts.updateSchool()
}
func (opts UpdateSchoolOptions) updateSchool() error {
	var (
		SQL    = `UPDATE Schools SET`
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
	if opts.Avg_gpa != 0 {
		SQL += fmt.Sprintf(" avg_gpa = $%d,", i)
		values = append(values, opts.Avg_gpa)
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
		fmt.Println(err)
		log.Println(" err : ", err)
		return nil
	}
	defer db.Close()
	//fmt.Println(SQL)
	_, err = db.Exec(SQL, values...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}
func (opts UpdateSchoolOptions) UpdateHelper() (School, error) {
	return opts.updateHelper()
}

func (opts UpdateSchoolOptions) updateHelper() (School, error) {
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
			if contains(school.ProfessorList, professor) {
				fmt.Printf("removing : %v ", professor)
				school.ProfessorList = remove(school.ProfessorList, professor)
			}
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
			//fmt.Printf("removing : %v ", class)
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
			if contains(school.StudentList, student) {
				//fmt.Printf("removing : %v ", student)
				school.StudentList = remove(school.StudentList, student)
			}
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

func remove(slice []string, str string) []string {
	index := -1
	for i, s := range slice {
		if s == str {
			index = i
			break
		}
	}
	if index == -1 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
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
		school = School{}
		// ownerID       sql.NullString
		professorList sql.NullString
		classList     sql.NullString
		StudentList   sql.NullString
		schAvg        sql.NullFloat64
	)
	err := row.Scan(
		&school.SchoolID,
		// &ownerID,
		&school.SchoolName,
		&schAvg,
		&school.Ranking,
		&professorList,
		&classList,
		&StudentList,
	)
	if err != nil {
		return School{}, err
	}
	if schAvg.Valid {
		var value interface{}
		value, err = schAvg.Value()
		if err == nil {
			school.AvgGPA = value.(float64)
		} else {
			fmt.Println(err)
		}
	}
	// if ownerID.Valid {
	// 	school.OwnerID = ownerID.String
	// }
	if professorList.Valid {
		school.ProfessorList = removeBrackets(strings.Split(professorList.String, ","))
	}
	if classList.Valid {

		school.ClassList = removeBrackets(strings.Split(classList.String, ","))
	}
	if StudentList.Valid {
		school.StudentList = removeBrackets(strings.Split(StudentList.String, ","))
	}
	return school, nil

}
func ScanSchools(rows *sql.Rows) ([]School, error) {
	return scanSchools(rows)
}
func scanSchools(rows *sql.Rows) ([]School, error) {
	defer rows.Close()
	var (
		schools = []School{}
		//ownerID       sql.NullString
		professorList sql.NullString
		classList     sql.NullString
		StudentList   sql.NullString
		schAvg        sql.NullFloat64
	)
	for rows.Next() {
		school := School{}
		err := rows.Scan(
			&school.SchoolID,
			//&ownerID,
			&school.SchoolName,
			&schAvg,
			&school.Ranking,
			&professorList,
			&classList,
			&StudentList,
		)
		if err != nil {
			return nil, err
		}
		if schAvg.Valid {
			var value interface{}
			value, err = schAvg.Value()
			if err == nil {
				school.AvgGPA = math.Round(value.(float64)*100) / 100

			} else {
				fmt.Println(err)
			}
		}
		// if ownerID.Valid {
		// 	school.OwnerID = ownerID.String
		// }
		if professorList.Valid {
			school.ProfessorList = removeBrackets(strings.Split(professorList.String, ","))
		}
		if classList.Valid {
			school.ClassList = removeBrackets(strings.Split(classList.String, ","))
		}
		if StudentList.Valid {
			school.StudentList = removeBrackets(strings.Split(StudentList.String, ","))
		}
		schools = append(schools, school)
	}
	if err := rows.Err(); err != nil {
		return schools, err
	}
	return schools, nil
}
