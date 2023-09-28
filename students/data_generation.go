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
	schoolID := uuid.New().String()
	fmt.Println("hit1")
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
	fmt.Println("hit2")
	fmt.Println("hit3")
	for _, stu := range roster {
		studentList = append(studentList, stu.StudentID)
	}
	for _, prof := range profList {
		professorList = append(professorList, prof.ProfessorID)
	}
	for _, class := range classList {
		classesList = append(classesList, class.ClassID)
	}
	fmt.Println("hit4")
	UpdateOperationStatus(operationID, "Compiling Data...", nil)
	school, err := CreateSchool(schoolID, schoolName, ownerID, professorList, classesList, studentList)
	if err != nil {
		return school, err
	}
	UpdateOperationStatus(operationID, "Uploading Data...", nil)
	avg_gpa, err := BatchUploadData(schoolID, operationID, roster, profList, classList, reportCardList, nil)
	if err != nil {
		return school, err
	}
	school.AvgGPA = avg_gpa
	fmt.Println("hit5")
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

func BatchUploadData(schoolID, operationID string, studentList []Student, profs []Professor, classes []Class, reportCards []ReportCard, err error) (float64, error) {
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	UpdateOperationStatus(operationID, "Compiling report card data...", nil)
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
	fmt.Println("batch1")
	UpdateOperationStatus(operationID, "Uploading Student Data...", nil)
	// upload as batch
	err = CreateNewStudents(schoolID, studentList)
	if err != nil {
		return 0, err
	}
	fmt.Println("batch2")
	fmt.Println("students successfully created")
	UpdateOperationStatus(operationID, "Uploading Report Card Data...", nil)
	err = CreateReportCards(reportCards)
	if err != nil {
		return 0, err
	}
	fmt.Println("batch3")
	fmt.Println("report cards successfully created")
	UpdateOperationStatus(operationID, "Uploading Professor Data...", nil)
	err = CreateProfessors(schoolID, profs)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	fmt.Println("batch4")
	fmt.Println("professors successfully created")
	UpdateOperationStatus(operationID, "Uploading Class Data...", nil)
	err = CreateNewClasses(classes)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// UPDATING DATA SHOULD BE MOVED
	// update classlists of professors
	// ~~~~~~~~~
	// THIS was probably removed due to performance issues however I don't think the classlists are being set
	// so this breaks everything the classlists seem to not be connected in the db. AND if the classlist
	// cant produce the related stuents in the list then it cant grab the db information needed to update the avg
	// ~~~~~~~~~
	UpdateOperationStatus(operationID, "Compiling school averages...", nil)
	err = UpdateProfessorsClassList(classes)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// update class Avgs
	// err = UpdateClassAvgs(classes)
	// if err != nil {
	// 	return err
	// }
	// return nil
	schoolAvg, err := UpdateSchoolAvg(schoolID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	fmt.Println(schoolAvg)

	// I dont yet know how performant this I may not want it here.
	err = UpdateSchoolRankings()
	if err != nil {
		return 0, err
	}

	fmt.Println("professors successfully created")
	fmt.Printf("Generated: %d students, %d report cards, %d professors and %d classes with an Average GPA of: %d", len(studentList), len(reportCards), len(profs), len(classes), schoolAvg)
	fmt.Println("")
	return schoolAvg, nil
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
	fmt.Printf("generating data for grade %v", grade)
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
	classes, err := GenerateClasses(profs, studentIDs, grade)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// generate report cards
	reportCards, err := GenerateReportCards(studentList, classes)
	if err != nil {
		fmt.Printf("generation in report cards %v", err)
		return nil, nil, nil, nil, err
	}
	fmt.Printf("generated data for grade %v", grade)
	return studentList, profs, classes, reportCards, nil
}

func GenerateClasses(professors []Professor, studentIDs []string, grade int) ([]Class, error) {
	var classes []Class
	classNames := [5]string{"math", "science", "english", "physicaled", "lunch"}
	for k, prof := range professors {
		class := Class{
			ClassID:       uuid.New().String(),
			TeachingGrade: grade,
			ProfessorID:   prof.ProfessorID,
			Subject:       classNames[k],
			Roster:        studentIDs,
		}
		classes = append(classes, class)
	}
	return classes, nil
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
		if i == 1 {
			fmt.Printf("dob check: %v", dob)
			fmt.Println()
		}
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

func CreateNewClasses(classes []Class) error {
	placeholders := make([]string, 0, len(classes))
	batchVals := make([]interface{}, 0, len(classes)*5)
	for n, class := range classes {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", n*5+1, n*5+2, n*5+3, n*5+4, n*5+5))
		batchVals = append(batchVals, class.ClassID)
		batchVals = append(batchVals, class.TeachingGrade)
		batchVals = append(batchVals, class.ProfessorID)
		batchVals = append(batchVals, class.Subject)
		batchVals = append(batchVals, pq.Array(class.Roster))
	}
	insertStatement := fmt.Sprintf(`INSERT INTO Classes("class_id","teaching_grade","professor_id","subject","roster") values %s`, strings.Join(placeholders, ","))
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
		fmt.Println(err)
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, batchVals...)
	if err != nil {
		return err
	}
	return nil
}

// CreateNewStudents uses batch processing to commit all in one db hit to be more performant
func CreateNewStudents(schoolID string, students []Student) error {
	return createNewStudents(schoolID, students)
}
func createNewStudents(schoolID string, students []Student) error {
	placeholders := make([]string, 0, len(students))
	batchVals := make([]interface{}, 0, len(students)*9)
	for n, student := range students {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", n*9+1, n*9+2, n*9+3, n*9+4, n*9+5, n*9+6, n*9+7, n*9+8, n*9+9))
		batchVals = append(batchVals, student.StudentID)
		batchVals = append(batchVals, schoolID)
		batchVals = append(batchVals, student.Name)
		batchVals = append(batchVals, student.CurrentYear)
		batchVals = append(batchVals, student.GraduationYear)
		batchVals = append(batchVals, student.AvgGPA)
		batchVals = append(batchVals, student.Age)
		batchVals = append(batchVals, student.Dob)
		batchVals = append(batchVals, student.Enrolled)
	}
	insertStatement := fmt.Sprintf(`INSERT into STUDENTS("student_id","school_id","name","current_year","graduation_year","avg_gpa","age","dob","enrolled") values %s`, strings.Join(placeholders, ","))
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Println("err %w", err)
	}
	defer db.Close()
	_, err = db.Exec(insertStatement, batchVals...)
	if err != nil {
		fmt.Println(err)
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

func CreateProfessors(schoolID string, profs []Professor) error {
	return createProfessors(schoolID, profs)
}
func createProfessors(schoolID string, profs []Professor) error {
	var placeholders []string
	var values []interface{}
	for i, prof := range profs {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", (3*i)+1, (3*i)+2, (3*i)+3))
		values = append(values, prof.ProfessorID, schoolID, prof.Name)
	}

	insertStatement := fmt.Sprintf(`INSERT INTO Professors("professor_id","school_id","name") VALUES %s`, strings.Join(placeholders, ", "))
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
