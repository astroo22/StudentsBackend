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

type User struct {
	OwnerID        string
	UserName       string
	Email          string
	HashedPassword string
	SchoolList     []string
}
type CreateNewUserOptions struct {
	UserName       string
	Email          string
	HashedPassword string
}

type UpdateUserOptions struct {
	// will not change
	OwnerID string
	// will change
	UserName         string
	Email            string
	HashedPassword   string
	AddSchoolList    []string
	RemoveSchoolList []string
}

func (opts CreateNewUserOptions) CreateNewUser() (User, error) {
	return opts.createNewUser()
}

func (opts CreateNewUserOptions) createNewUser() (User, error) {
	var (
		NewID       = uuid.New().String()
		SQL         = `INSERT INTO Users("owner_id","user_name",`
		values      []interface{}
		placeholder = "$1,$2,$3"
		user        = User{}
	)
	ret := User{}
	values = append(values, NewID, opts.UserName)

	if len(opts.Email) != 0 {
		SQL += `"email",`
		values = append(values, opts.Email)
		placeholder += ",$4"
		ret.Email = opts.Email
	}

	SQL += fmt.Sprintf(`"hashed_password") VALUES (%s)`, placeholder)
	values = append(values, opts.HashedPassword)

	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}
	defer db.Close()
	_, err = db.Exec(SQL, values...)
	if err != nil {
		return user, err
	}
	ret.OwnerID = NewID
	ret.UserName = opts.UserName

	return ret, nil
}

func GetUser(ownerID string) (User, error) {
	return getUser(ownerID)
}
func getUser(ownerID string) (User, error) {
	getStatement := `SELECT * FROM Users WHERE owner_id = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
		return User{}, err
	}
	defer db.Close()
	user, err := ScanUser(db.QueryRow(getStatement, ownerID))
	if err != nil {
		return User{}, err
	}
	return user, nil
}
func GetUserByUserName(userName string) (User, error) {
	return getUserByUserName(userName)
}
func getUserByUserName(userName string) (User, error) {
	getStatement := `SELECT * FROM Users WHERE user_name = $1`
	db, err := sqlgeneric.Init()
	if err != nil {
		fmt.Println(err)
		log.Printf(" err : %v", err)
		return User{}, err
	}
	defer db.Close()
	user, err := ScanUser(db.QueryRow(getStatement, userName))
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (opts UpdateUserOptions) UpdateUser() error {
	return opts.updateUser()
}
func (opts UpdateUserOptions) updateUser() error {
	var (
		SQL    = `UPDATE Users SET`
		values []interface{}
		i      = 2
	)
	values = append(values, opts.OwnerID)
	if len(opts.UserName) != 0 {
		SQL += fmt.Sprintf(" user_name = $%d,", i)
		values = append(values, opts.UserName)
		i++
	}
	if len(opts.Email) != 0 {
		if opts.Email == "-1" {
			SQL += fmt.Sprintf(" email = $%d,", i)
			values = append(values, "")
			i++
		} else {
			SQL += fmt.Sprintf(" email = $%d,", i)
			values = append(values, opts.Email)
			i++
		}
	}
	if len(opts.HashedPassword) != 0 {
		SQL += fmt.Sprintf(" hashed_password = $%d,", i)
		values = append(values, opts.HashedPassword)
		i++
	}
	if len(opts.AddSchoolList) != 0 || len(opts.RemoveSchoolList) != 0 {
		schoolList, err := opts.prepUserSchoolListUpdate()
		if err != nil {
			return err
		}
		SQL += fmt.Sprintf(" school_list = $%d", i)
		values = append(values, pq.Array(schoolList))
		i++
	}
	if SQL[len(SQL)-1] == ',' {
		SQL = SQL[:len(SQL)-1]
	}
	SQL += " WHERE owner_id = $1"
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println(" err : ", err)
	}

	defer db.Close()
	_, err = db.Exec(SQL, values...)
	if err != nil {
		return err
	}
	return nil
}

func (opts UpdateUserOptions) prepUserSchoolListUpdate() ([]string, error) {
	var ret []string
	schools, err := GetAllSchoolsForUser(opts.OwnerID)
	if err != nil {
		return nil, err
	}
	for _, school := range schools {
		if !slices.Contains(opts.RemoveSchoolList, school.SchoolID) {
			ret = append(ret, school.SchoolID)
		}
	}
	ret = append(ret, opts.AddSchoolList...)
	return ret, nil
}

func DeleteUser(ownerID string) error {
	return deleteUser(ownerID)
}

func deleteUser(ownerID string) error {
	db, err := sqlgeneric.Init()
	if err != nil {
		log.Println("Error initializing database: ", err)
		return err
	}
	defer db.Close()

	// transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return err
	}

	// SQL statements
	deleteSchools := `DELETE FROM Schools WHERE owner_id = $1`

	deleteStudents := `DELETE FROM Students WHERE school_id IN 
					(SELECT school_id FROM Schools WHERE owner_id = $1)`

	deleteReportCards := `DELETE FROM ReportCards WHERE student_id IN 
					(SELECT student_id FROM Students WHERE school_id IN 
					(SELECT school_id FROM Schools WHERE owner_id = $1))`

	deleteProfessors := `DELETE FROM Professors WHERE school_id IN 
					(SELECT school_id FROM Schools WHERE owner_id = $1)`

	deleteClasses := `DELETE FROM Classes WHERE professor_id IN 
					(SELECT professor_id FROM Professors WHERE school_id IN 
					(SELECT school_id FROM Schools WHERE owner_id = $1))`

	deleteUser := `DELETE FROM Users WHERE owner_id = $1`

	// Execute deletion queries
	if _, err = tx.Exec(deleteReportCards, ownerID); err != nil {
		tx.Rollback()
		log.Println("Error deleting report cards: ", err)
		return err
	}
	if _, err = tx.Exec(deleteStudents, ownerID); err != nil {
		tx.Rollback()
		log.Println("Error deleting students: ", err)
		return err
	}
	if _, err = tx.Exec(deleteClasses, ownerID); err != nil {
		tx.Rollback()
		log.Println("Error deleting classes: ", err)
		return err
	}
	if _, err = tx.Exec(deleteProfessors, ownerID); err != nil {
		tx.Rollback()
		log.Println("Error deleting professors: ", err)
		return err
	}
	if _, err = tx.Exec(deleteSchools, ownerID); err != nil {
		tx.Rollback()
		log.Println("Error deleting schools: ", err)
		return err
	}
	if _, err = tx.Exec(deleteUser, ownerID); err != nil {
		tx.Rollback()
		log.Println("Error deleting user: ", err)
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction: ", err)
		return err
	}

	return nil
}

// func deleteUser(ownerID string) error {
// 	SQL := `DELETE FROM Users WHERE owner_id = $1`
// 	db, err := sqlgeneric.Init()
// 	if err != nil {
// 		log.Println(" err : ", err)
// 	}
// 	defer db.Close()
// 	_, err = db.Exec(SQL, ownerID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// SCANS
func ScanUser(row *sql.Row) (User, error) {
	return scanUser(row)
}

func scanUser(row *sql.Row) (User, error) {
	var (
		user       = User{}
		schoolList sql.NullString
		email      sql.NullString
	)
	err := row.Scan(
		&user.OwnerID,
		&user.UserName,
		&email,
		&user.HashedPassword,
		&schoolList)
	if err != nil {
		return User{}, err
	}
	if schoolList.Valid {
		user.SchoolList = strings.Split(schoolList.String, ",")
	}
	if email.Valid {
		user.Email = email.String
	}
	return user, nil
}
