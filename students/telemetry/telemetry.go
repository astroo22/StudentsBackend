package telemetry

import (
	"database/sql"
	"fmt"
	"math"
	"students/client"
	"students/sqlgeneric"
	"students/students"

	"github.com/lib/pq"
)

func GetGradeAvgForSchool(schoolID string) ([]client.GradeAvg_API, error) {
	var (
		grdAvg sql.NullFloat64
	)
	// query := `
	//     SELECT c.teaching_grade as grade, coalesce(AVG(c.class_avg),0.0) as avg_gpa
	//     FROM schools s, unnest(s.class_list) cl, classes c
	//     WHERE s.school_id = $1 AND cl = c.class_id
	//     GROUP BY grade
	//     ORDER BY avg_gpa DESC
	// `
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

	gradeAvgs := []client.GradeAvg_API{}
	for rows.Next() {
		gradeAvg := client.GradeAvg_API{}
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
				gradeAvg.AvgGPA = value.(float64)
			} else {
				fmt.Println(err)
			}
		}
		fmt.Printf("GradeAvg: %v", gradeAvg.AvgGPA)
		fmt.Println("")
		gradeAvgs = append(gradeAvgs, gradeAvg)
	}

	return gradeAvgs, nil
}

// FigureDerivedDate: WIll run a db update on avg's should be used after a generation
// will probably tie to a handler to get updates whenever
// func FigureDerivedData() error {
// 	// err := UpdateClassAvgs()
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	err := UpdateAllProfessorStudentAvgs()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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

func UpdateProfessorsStudentAvgs(professors []string) error {
	return updateProfessorsStudentAvgs(professors)
}
func updateProfessorsStudentAvgs(professors []string) error {
	db, err := sqlgeneric.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE Professors SET student_avg = (
        SELECT COALESCE(AVG(class_avg),0) 
        FROM Classes 
        WHERE class_id = ANY(SELECT UNNEST(class_list) FROM Professors WHERE professor_id=Professors.professor_id)
    )
	WHERE professor_id = ANY($1)`

	_, err = db.Exec(query, professors)
	if err != nil {
		return err
	}
	return nil
}

// This is more performant and useful now but I can probably find a way to batch update the class updates later
// Updates all class avgs in the table
func UpdateClassAvgs(classList []students.Class) error {
	return updateClassAvgs(classList)
}

func updateClassAvgs(classList []students.Class) error {
	db, err := sqlgeneric.Init()
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Println(classList[0].Roster)
	for _, class := range classList {
		//fmt.Printf("class",class.)
		// Get the report card grades for each student in the class
		// query := `SELECT math, science, english, physical_ed, lunch FROM ReportCards WHERE student_id = ANY($1)`
		// rows, err := db.Query(query, pq.Array(class.Roster))
		// if err != nil {
		// 	return err
		// }
		var (
			reportCards []students.ReportCard
		)
		for _, studentID := range class.Roster {
			rc, err := students.GetReportCard(studentID)
			if err != nil {
				return err
			}
			reportCards = append(reportCards, rc)
		}

		// reportCards, err := students.GetReportCards(class.Roster)
		// if err != nil {
		// 	return err
		// }
		fmt.Println(reportCards)
		if len(reportCards) == 0 {
			return fmt.Errorf("no report cards found")
		}
		// Calculate the class average
		var totalGrade, numStudents int
		for _, reportCard := range reportCards {
			//var math, science, english, physicalEd, lunch float64
			//err = rows.Scan(&math, &science, &english, &physicalEd, &lunch)
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
		fmt.Println(totalGrade)
		fmt.Println(numStudents)
		classAvg := math.Round(float64(totalGrade) / float64(numStudents) / 100.0)

		fmt.Println(classAvg)
		// Update the class record with the calculated class average
		query := `UPDATE Classes SET class_avg = $1 WHERE class_id = $2`
		_, err = db.Exec(query, classAvg, class.ClassID)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateProfessorsClassList(classList []students.Class) error {
	return updateProfessorsClassList(classList)
}
func updateProfessorsClassList(classList []students.Class) error {
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
