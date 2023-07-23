package students

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"students/sqlgeneric"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type SchoolCreationStatus_API struct {
	Status      string `json:"status"`
	OperationID string `json:"operation_id"`
	School      School `json:"school,omitempty"`
	Error       error  `json:"error,omitempty"`
}

var schoolCreationStatuses = make(map[string]*SchoolCreationStatus_API)
var mutex = &sync.Mutex{}

func CreateOperationEntry(operationID, statusMessage string) {
	status := &SchoolCreationStatus_API{
		Status: statusMessage,
	}
	mutex.Lock()
	schoolCreationStatuses[operationID] = status
	mutex.Unlock()
}

func UpdateOperationStatus(operationID, statusMessage string, err error) bool {
	// checks if exists
	if _, ok := schoolCreationStatuses[operationID]; !ok {
		// operationID does not exist in the map, handle the error or return
		fmt.Println("Error: Invalid operationID")
		fmt.Println(operationID)
		fmt.Println(schoolCreationStatuses)
		return false
	}
	// means function was hit to check if exists
	if len(statusMessage) == 0 {
		return true
	}
	if err != nil {
		status := &SchoolCreationStatus_API{
			Status: statusMessage,
			Error:  err,
		}
		mutex.Lock()
		schoolCreationStatuses[operationID] = status
		mutex.Unlock()
	} else {
		status := &SchoolCreationStatus_API{
			Status: statusMessage,
		}
		mutex.Lock()
		schoolCreationStatuses[operationID] = status
		mutex.Unlock()
	}
	return true
}

func GetOperationStatus(operationID string) (*SchoolCreationStatus_API, error) {
	mutex.Lock()
	status, exists := schoolCreationStatuses[operationID]
	mutex.Unlock()
	if !exists {
		return nil, fmt.Errorf("no such operation")
	}
	return status, nil
}

// CreateSchool: creates a school.
func NewSchool(operationID string, studentsPerGrade int, ownerID, schoolName string) (School, error) {
	return newSchool(operationID, studentsPerGrade, ownerID, schoolName)
}

func newSchool(operationID string, studentsPerGrade int, ownerID, schoolName string) (School, error) {
	var (
		// struct holders
		roster         []Student
		profList       []Professor
		classList      []Class
		reportCardList []ReportCard
		// empty for returns until last line
		school School
		// list for school
		studentList   []string
		professorList []string
		classesList   []string
		updateMessage string
	)

	for i := 1; i <= 12; i++ {
		stus, profs, classes, rcs, err := GenerateData(studentsPerGrade, i)
		if err != nil {
			return school, err
		}
		roster = append(roster, stus...)
		profList = append(profList, profs...)
		classList = append(classList, classes...)
		reportCardList = append(reportCardList, rcs...)
		updateMessage = fmt.Sprintf("Generating grade: %d", i)
		UpdateOperationStatus(operationID, updateMessage, nil)
	}
	UpdateOperationStatus(operationID, "Uploading Data...", nil)
	err := BatchUploadData(roster, profList, classList, reportCardList, nil)
	if err != nil {
		return school, err
	}
	for _, stu := range roster {
		studentList = append(studentList, stu.StudentID)
	}
	for _, prof := range profList {
		professorList = append(professorList, prof.ProfessorID)
	}
	for _, class := range classList {
		classesList = append(classesList, class.ClassID)
	}
	UpdateOperationStatus(operationID, "Compiling Data...", nil)
	school, err = CreateSchool(schoolName, ownerID, professorList, classesList, studentList)
	if err != nil {
		return school, err
	}
	return school, nil
}

// func AdminGenerateTestSchools() error {
// 	var (
// 		owernerID   = "The New Vibe"
// 		schoolName1 = "Busta Rhymes Academy"
// 		schoolName2 = "Rick and Morty Vindicators 4"
// 		schoolName3 = "PLUS ULTRA ACADEMY"
// 		schoolName4 = "Xavier Institue for Higher Learning"
// 		schoolName5 = "Institue for UnderWater Basket Weaving"
// 		stdPerGrade = 20
// 	)
// 	_, err := NewSchool(stdPerGrade, owernerID, schoolName1)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = NewSchool(stdPerGrade, owernerID, schoolName2)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = NewSchool(stdPerGrade, owernerID, schoolName3)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = NewSchool(stdPerGrade, owernerID, schoolName4)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = NewSchool(stdPerGrade, owernerID, schoolName5)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func BatchUploadData(studentList []Student, profs []Professor, classes []Class, reportCards []ReportCard, err error) error {
	if err != nil {
		return err
	}
	// set avg gpas also grab studentids
	for i, student := range studentList {
		//roster = append(roster, student.StudentID)
		for _, rc := range reportCards {
			if student.StudentID == rc.StudentID {
				studentList[i].AvgGPA = math.Round(((rc.English+rc.Lunch+rc.Math+rc.PhysicalED+rc.Science)/5)*100) / 100
				continue
			}
		}
	}

	// upload as batch
	err = CreateNewStudents(studentList)
	if err != nil {
		return err
	}
	fmt.Println("students successfully created")
	err = CreateReportCards(reportCards)
	if err != nil {
		return err
	}
	fmt.Println("report cards successfully created")
	err = CreateProfessors(profs)
	if err != nil {
		return err
	}

	// UPDATING DATA SHOULD BE MOVED
	// update classlists of professors

	fmt.Println("professors successfully created")
	fmt.Printf("Generated: %d students, %d report cards, %d professors and %d classes", len(studentList), len(reportCards), len(profs), len(classes))
	fmt.Println("")
	return nil
}

// func RunTelemetry(classes []Class) error {
// 	err := UpdateProfessorsClassList(classes)
// 	if err != nil {
// 		return err
// 	}
// 	// update class Avgs
// 	err = UpdateClassAvgs(classes)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// GenerateTestData:
func GenerateData(numStutotal int, grade int) ([]Student, []Professor, []Class, []ReportCard, error) {
	studentList, err := GenerateStudents(numStutotal, grade)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	profs, err := GenerateProfessors(5)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	studentIDs := []string{}
	for _, student := range studentList {
		studentIDs = append(studentIDs, student.StudentID)
	}
	var classes []Class
	classNames := [5]string{"math", "science", "english", "physicaled", "lunch"}

	for k, prof := range profs {
		class, err := CreateClass(grade, prof.ProfessorID, classNames[k], studentIDs)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		classes = append(classes, class)
	}
	// generate report cards
	reportCards, err := GenerateReportCards(studentList, classes)
	if err != nil {
		fmt.Printf("generation in report cards %v", err)
		return nil, nil, nil, nil, err
	}

	return studentList, profs, classes, reportCards, nil
}

// STUDENTS
func GenerateStudents(numStudents int, grade int) ([]Student, error) {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	studentList := make([]Student, numStudents)

	for i := 0; i < numStudents; i++ {
		currentYear := grade
		graduationYear := time.Now().Year() + (12 - grade)

		age := rng.Intn(2) + 7 + grade

		// Generate a random date of birth
		currentYearInCalendar := time.Now().Year()
		dob := time.Date((currentYearInCalendar - age), time.Month(rng.Intn(12)+1), rng.Intn(28)+1, 0, 0, 0, 0, time.UTC)
		enrolled := rng.Intn(2) == 1

		student := Student{
			StudentID:      uuid.New().String(),
			Name:           randomdata.FullName(randomdata.RandomGender),
			CurrentYear:    currentYear,
			GraduationYear: graduationYear,
			Age:            age,
			Dob:            dob,
			Enrolled:       enrolled,
		}
		// Add the new student record to the array
		studentList[i] = student
	}
	return studentList, nil
}

// CreateNewStudents uses batch processing to commit all in one db hit to be more performant
func CreateNewStudents(students []Student) error {
	return createNewStudents(students)
}
func createNewStudents(students []Student) error {
	placeholders := make([]string, 0, len(students))
	batchVals := make([]interface{}, 0, len(students)*8)
	for n, student := range students {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", n*8+1, n*8+2, n*8+3, n*8+4, n*8+5, n*8+6, n*8+7, n*8+8))
		batchVals = append(batchVals, student.StudentID)
		batchVals = append(batchVals, student.Name)
		batchVals = append(batchVals, student.CurrentYear)
		batchVals = append(batchVals, student.GraduationYear)
		batchVals = append(batchVals, student.AvgGPA)
		batchVals = append(batchVals, student.Age)
		batchVals = append(batchVals, student.Dob)
		batchVals = append(batchVals, student.Enrolled)
	}
	insertStatement := fmt.Sprintf(`INSERT into STUDENTS("student_id","name","current_year","graduation_year","avg_gpa","age","dob","enrolled") values %s`, strings.Join(placeholders, ","))
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println("err %w", err)
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, batchVals...)
	if err != nil {
		return err
	}
	return nil
}

// REPORT CARDS
func GenerateReportCards(studentList []Student, classList []Class) ([]ReportCard, error) {
	return generateReportCards(studentList, classList)
}
func generateReportCards(studentList []Student, classList []Class) ([]ReportCard, error) {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	var (
		reportCards []ReportCard
	)
	classIDs := []string{}
	for _, class := range classList {
		classIDs = append(classIDs, class.ClassID)
	}
	for _, student := range studentList {
		reportCard := ReportCard{
			StudentID:  student.StudentID,
			Math:       float64(rng.Intn(400)) / 100.0,
			Science:    float64(rng.Intn(400)) / 100.0,
			English:    float64(rng.Intn(400)) / 100.0,
			PhysicalED: float64(rng.Intn(400)) / 100.0,
			Lunch:      float64(rng.Intn(400)) / 100.0,
			ClassList:  classIDs,
		}
		reportCards = append(reportCards, reportCard)
	}
	return reportCards, nil
}

func CreateReportCards(reportCards []ReportCard) error {
	return createReportCards(reportCards)
}
func createReportCards(reportCards []ReportCard) error {
	var (
		values       []interface{}
		placeholders []string
		i            = 1
	)

	for _, reportCard := range reportCards {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i, i+1, i+2, i+3, i+4, i+5, i+6))
		values = append(values, reportCard.StudentID, reportCard.Math, reportCard.Science, reportCard.English, reportCard.PhysicalED, reportCard.Lunch, pq.Array(reportCard.ClassList))
		i += 7
	}
	query := fmt.Sprintf(`INSERT INTO ReportCards("student_id", "math", "science", "english", "physical_ed", "lunch","class_list") VALUES %s`, strings.Join(placeholders, ", "))
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println("err: ", err)
		return err
	}
	defer db.Close()
	_, err = db.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

// PROFESSORS
func GenerateProfessors(numProfs int) ([]Professor, error) {
	return generateProfessors(numProfs)
}
func generateProfessors(numProfs int) ([]Professor, error) {
	profList := make([]Professor, numProfs)
	for i := 0; i < numProfs; i++ {
		profID := uuid.New().String()
		name := randomdata.FullName(randomdata.RandomGender)

		prof := Professor{
			ProfessorID: profID,
			Name:        name,
		}
		profList[i] = prof
	}
	return profList, nil
}

func CreateProfessors(profs []Professor) error {
	return createProfessors(profs)
}
func createProfessors(profs []Professor) error {
	var placeholders []string
	var values []interface{}
	for i, prof := range profs {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", (2*i)+1, (2*i)+2))
		values = append(values, prof.ProfessorID, prof.Name)
	}

	insertStatement := fmt.Sprintf(`INSERT INTO Professors("professor_id","name") VALUES %s`, strings.Join(placeholders, ", "))
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(insertStatement, values...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteReportCards(students []Student) error {
	return deleteReportCards(students)
}
func deleteReportCards(students []Student) error {
	batch := make([]string, 0, len(students))
	batchVals := make([]interface{}, 0, len(students))
	for n, student := range students {
		batch = append(batch, fmt.Sprintf("$%d", n+1))
		batchVals = append(batchVals, student.StudentID)
	}
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

// func DeleteTables() error {
// 	return deleteTables()
// }
// func deleteTables() error {
// 	var (
// 		students    = `DELETE FROM Students`
// 		reportCards = `DELETE FROM ReportCards`
// 		classes     = `DELETE FROM Classes`
// 		professors  = `DELETE FROM Professors`
// 	)

// 	db, err := sqlgeneric.Init()
// 	if err != nil {
// 		log.Println(" err : ", err)
// 	}
// 	_, err = db.Exec(students)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = db.Exec(reportCards)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = db.Exec(classes)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = db.Exec(professors)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
