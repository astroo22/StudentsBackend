package telemetry

import (
	"students/sqlgeneric"
	"students/students"

	"github.com/lib/pq"
)

// FigureDerivedDate: WIll run a db update on avg's should be used after a generation
// will probably tie to a handler to get updates whenever
func FigureDerivedData() error {
	err := UpdateClassAvgs()
	if err != nil {
		return err
	}
	err = UpdateAllProfessorStudentAvgs()
	if err != nil {
		return err
	}
	return nil
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

func UpdateAllProfessorStudentAvgs() error {
	return updateAllProfessorStudentAvgs()
}
func updateAllProfessorStudentAvgs() error {
	db, err := sqlgeneric.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE Professors SET student_avg = (
        SELECT COALESCE(AVG(class_avg),0) 
        FROM Classes 
        WHERE class_id = ANY(SELECT UNNEST(class_list) FROM Professors WHERE professor_id=Professors.professor_id)
    )`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// Updates all class avgs in the table
func UpdateClassAvgs() error {
	return updateClassAvgs()
}

// I can Probably make this more performant later
func updateClassAvgs() error {
	db, err := sqlgeneric.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	// Get all class records
	query := `SELECT class_id, roster, subject FROM Classes`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var classID string
		var roster []string
		var subject string
		err = rows.Scan(&classID, pq.Array(&roster), &subject)
		if err != nil {
			return err
		}

		// Get the report card grades for each student in the class
		query = `SELECT math, science, english, physical_ed, lunch FROM ReportCards WHERE student_id = ANY($1)`
		rows2, err := db.Query(query, pq.Array(roster))
		if err != nil {
			return err
		}
		defer rows2.Close()

		// Calculate the class average
		var totalGrade, numStudents int
		for rows2.Next() {
			var math, science, english, physicalEd, lunch float64
			err = rows2.Scan(&math, &science, &english, &physicalEd, &lunch)
			if err != nil {
				return err
			}
			switch subject {
			case "math":
				totalGrade += int(math * 100)
			case "science":
				totalGrade += int(science * 100)
			case "english":
				totalGrade += int(english * 100)
			case "physical_ed":
				totalGrade += int(physicalEd * 100)
			case "lunch":
				totalGrade += int(lunch * 100)
			}
			numStudents++
		}
		classAvg := float64(totalGrade) / float64(numStudents) / 100.0

		// Update the class record with the calculated class average
		query = `UPDATE Classes SET class_avg = $1 WHERE class_id = $2`
		_, err = db.Exec(query, classAvg, classID)
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
