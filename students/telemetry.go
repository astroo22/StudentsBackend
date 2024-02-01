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

type GradeAvg_API struct {
	Grade  int     `json:"grade"`
	AvgGPA float64 `json:"avg_gpa"`
}

func GetGradeAvgForSchool(schoolID string) ([]GradeAvg_API, error) {
	return getGradeAvgForSchool(schoolID)
}

func getGradeAvgForSchool(schoolID string) ([]GradeAvg_API, error) {
	var (
		grdAvg sql.NullFloat64
	)
	query := `
			SELECT c.teaching_grade as grade, COALESCE(AVG(c.class_avg), 0.0) as avg_gpa
			FROM classes c
			WHERE c.class_id IN (
				SELECT UNNEST(class_list) 
				FROM schools 
				WHERE school_id = $1
				)
		GROUP BY grade
		ORDER BY avg_gpa DESC
	`
	db, err := sqlgeneric.Init()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gradeAvgs := []GradeAvg_API{}
	for rows.Next() {
		gradeAvg := GradeAvg_API{}
		err := rows.Scan(
			&gradeAvg.Grade,
			&grdAvg,
		)
		if err != nil {
			return nil, err
		}
		if grdAvg.Valid {
			var value interface{}
			value, err = grdAvg.Value()
			if err == nil {
				//classAvg = math.Round(classAvg*100) / 100
				gradeAvg.AvgGPA = math.Round(value.(float64)*100) / 100
			} else {
				fmt.Println(err)
			}
		}
		gradeAvgs = append(gradeAvgs, gradeAvg)
	}

	return gradeAvgs, nil
}
func GetBestProfessors(schoolID string) ([]Professor, error) {
	return getBestProfessors(schoolID)
}
func getBestProfessors(schoolID string) ([]Professor, error) {
	school, err := GetSchool(schoolID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()

	placeholders := make([]string, 0, len(school.ProfessorList))
	batchVals := make([]interface{}, 0, len(school.ProfessorList))
	for i := 0; i < len(school.ProfessorList); i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		batchVals = append(batchVals, school.ProfessorList[i])
	}

	getStatement := fmt.Sprintf(`SELECT professor_id, name, student_avg, class_list FROM Professors WHERE professor_id IN (%s) ORDER BY student_avg DESC`, strings.Join(placeholders, ","))
	rows, err := db.Query(getStatement, batchVals...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	bestProfessors, err := ScanProfessors(rows)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bestProfessors, nil
}

// UpdateProfessorStudentAvg: updates specific professors student avg
func UpdateProfessorStudentAvg(profID string) (float64, error) {
	return updateProfessorStudentAvg(profID)
}
func updateProfessorStudentAvg(profID string) (float64, error) {
	db, err := sqlgeneric.Init()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var studentAvg float64
	query := `UPDATE Professors SET student_avg = (
		SELECT COALESCE(AVG(class_avg),0) 
		FROM Classes 
		WHERE class_id = ANY(SELECT UNNEST(class_list) FROM Professors WHERE professor_id=$1)
	)
	WHERE professor_id = $1
	RETURNING student_avg;
	`

	err = db.QueryRow(query, profID).Scan(&studentAvg)
	if err != nil {
		return 0, err
	}
	return studentAvg, nil
}

// might need to spin off thread here ~
// this is pretty fast now not worth a new thread since schools are limited in size
func UpdateProfessorsStudentAvgs(professors []string) error {
	return updateProfessorsStudentAvgs(professors)
}
func updateProfessorsStudentAvgs(professors []string) error {
	fmt.Printf("hit update professors for: %v", professors)
	fmt.Println()
	db, err := sqlgeneric.Init()
	if err != nil {
		return err
	}
	defer db.Close()
	result := make([]string, len(professors))
	for i, val := range professors {
		result[i] = strings.ReplaceAll(val, ",", " ")
	}
	//fmt.Println(len(result))
	query := `UPDATE Professors 
	SET student_avg = (
		SELECT COALESCE(AVG(class_avg), 0)
		FROM Classes 
		WHERE class_id = ANY(Professors.class_list) 
	)
	WHERE professor_id = ANY($1)`

	_, err = db.Exec(query, pq.Array(result))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("finished update for professors")
	return nil
}
func UpdateSchoolAvg(schoolID string) (float64, error) {
	return updateSchoolAvg(schoolID)
}

func updateSchoolAvg(schoolID string) (float64, error) {
	var (
		schoolAvg float64
	)
	err := UpdateClassAvgs(schoolID)
	if err != nil {
		return 0, err
	}
	classes, err := GetClassesForSchool(schoolID)
	if err != nil {
		return 0, err
	}
	for _, class := range classes {
		schoolAvg += class.ClassAvg
	}
	schoolAvg = math.Round(schoolAvg/float64(len(classes))*100) / 100
	//fmt.Println(schoolAvg)
	updateQuery := `UPDATE Schools SET avg_gpa = $1 WHERE school_id = $2`
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
		return 0, err
	}
	defer db.Close()
	_, err = db.Exec(updateQuery, schoolAvg, schoolID)
	if err != nil {
		return 0, err
	}
	return schoolAvg, nil

}

// This is more performant and useful now but I can probably find a way to batch update the class updates at the end later
// Updates all class avgs in the table
func UpdateClassAvgs(schoolID string) error {
	return updateClassAvgs(schoolID)
}

func updateClassAvgs(schoolID string) error {

	classList, err := GetClassesForSchool(schoolID)
	if err != nil {
		return err
	}
	db, err := sqlgeneric.Init()
	if err != nil {
		return err
	}
	defer db.Close()
	for _, class := range classList {
		var (
			reportCards []ReportCard
		)
		reportCards, err := GetReportCardsOfEnrolled(class.Roster)
		if err != nil {
			return err
		}
		//fmt.Println(reportCards)
		if len(reportCards) == 0 {
			return fmt.Errorf("no report cards found")
		}
		// Calculate the class average
		var totalGrade, numStudents int
		for _, reportCard := range reportCards {
			// math, science, english, physicalEd, lunch
			switch class.Subject {
			case "math":
				totalGrade += int(reportCard.Math * 100)
			case "science":
				totalGrade += int(reportCard.Science * 100)
			case "english":
				totalGrade += int(reportCard.English * 100)
			case "physicaled":
				totalGrade += int(reportCard.PhysicalED * 100)
			case "lunch":
				totalGrade += int(reportCard.Lunch * 100)
			default:
				fmt.Printf("class miss: %s", class.Subject)
			}
			numStudents++
		}
		classAvg := math.Round(float64(totalGrade)/float64(numStudents)) / 100.0
		//fmt.Println(classAvg)
		//classAvg = math.Round(classAvg*100) / 100
		// Update the class record with the calculated class average
		query := `UPDATE Classes SET class_avg = $1 WHERE class_id = $2`
		_, err = db.Exec(query, classAvg, class.ClassID)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateProfessorsClassList(classList []Class) error {
	return updateProfessorsClassList(classList)
}
func updateProfessorsClassList(classList []Class) error {
	db, err := sqlgeneric.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	// Group the classes by professor_id
	classMap := make(map[string][]string)
	for _, class := range classList {
		classMap[class.ProfessorID] = append(classMap[class.ProfessorID], class.ClassID)
	}

	// Batch update the professors' class lists
	// other batch method which is just as performant.
	query := `UPDATE Professors SET class_list = $1 WHERE professor_id = $2`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for profID, classes := range classMap {
		_, err := stmt.Exec(pq.Array(classes), profID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Shouldn't really run this unless in go routine
// even then this takes awhile and is pretty resource intensive so if my db gets big its not ideal
func UpdateAllSchoolAvgGpa() error {
	return updateAllSchoolAvgGpa()
}

func updateAllSchoolAvgGpa() error {
	schools, err := GetAllSchools()
	if err != nil {
		return err
	}
	for _, school := range schools {
		avg, err := UpdateSchoolAvg(school.SchoolID)
		if err != nil {
			return err
		}
		fmt.Println(avg)
	}
	return nil
}

// TODO: MIGHT BE WORTH HAVING THIS RUN ON CYCLE EVERY LIL BIT
// func UpdateStudentAvgs() error {
// 	// Get a list of all students from the database
// 	studentList, err := students.GetAllStudents()
// 	if err != nil {
// 		return err
// 	}
// 	reportCardList, err := students.GetAllReportCards()
// 	if err != nil {
// 		return err
// 	}
// 	// Update the average GPA for each student
// 	for _, student := range studentList {
// 		for _, reportCard := range reportCardList {
// 			if student.StudentID == reportCard.StudentID {
// 				avgGPA := (reportCard.Math + reportCard.Science + reportCard.English + reportCard.PhysicalED + reportCard.Lunch) / 5
// 				avgGPA = math.Round(avgGPA*100) / 100 // Round to two decimal places
// 				opts := students.UpdateStudentOptions{
// 					StudentID: student.StudentID,
// 					AvgGPA:    avgGPA,
// 					Enrolled:  true,
// 				}
// 				err = opts.UpdateStudent()
// 				if err != nil {
// 					return err
// 				}
// 				continue
// 			}
// 		}
// 	}
// 	return nil
// }
