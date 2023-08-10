package main

import (
	"fmt"
	"os"

	"net/http"
	"students/auth"
	"students/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {

	// HEY... LISTEN!: comment out this section if not on prod
	logFile, err := os.OpenFile("/var/log/backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(logFile)

	fmt.Println("HEREEE WE GOOOOOO")
	err = auth.LoadSecretKey()
	if err != nil {
		return
	}
	// var(
	// 	lastUpdate = time.Now()
	// 	interval = 10 * time.Minute
	// )

	// DATA GENERATION
	// err := telemetry.AdminGenerateTestSchools()
	// if err != nil {
	// 	fmt.Println("data generation failed in main")
	// }

	// Create the HTTP server and set the router
	//router := http.NewServeMux()
	router := mux.NewRouter()
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "prod" {
		router.Use(addCorsHeadersProd)
	} else {
		router.Use(addCorsHeaders)
	}
	// backend status
	router.HandleFunc("/status", handlers.BackendStatus).Methods("GET")

	// Set up the routes for Students
	//router.HandleFunc("/students", handlers.CreateStudentHandler).Methods("POST")
	router.HandleFunc("/students", handlers.GetAllStudentsHandler).Methods("GET")
	router.HandleFunc("/students/{student_id}", handlers.GetStudentHandler).Methods("GET")
	// router.HandleFunc("/students/{student_id}", handlers.UpdateStudentHandler).Methods("PUT")
	// router.HandleFunc("/students/{student_id}", handlers.DeleteStudentHandler).Methods("DELETE")

	// Set up the routes for ReportCards
	// router.HandleFunc("/reportcard", handlers.CreateReportCardHandler).Methods("POST")
	router.HandleFunc("/reportcard/{student_id}", handlers.GetReportCardHandler).Methods("GET")
	// router.HandleFunc("/reportcard/{student_id}", handlers.UpdateReportCardHandler).Methods("PUT")
	// router.HandleFunc("/reportcard/{student_id}", handlers.DeleteReportCardHandler).Methods("DELETE")

	// Set up the routes for Classes
	// router.HandleFunc("/classes", handlers.CreateClassHandler).Methods("POST")
	router.HandleFunc("/classes/{class_id}", handlers.GetClassHandler).Methods("GET")
	// router.HandleFunc("/classes/{class_id}", handlers.UpdateClassHandler).Methods("PUT")
	// router.HandleFunc("/classes/{class_id}", handlers.DeleteClassHandler).Methods("DELETE")

	// Set up the routes for Professors
	// router.HandleFunc("/professors", handlers.CreateProfessorHandler).Methods("POST")
	router.HandleFunc("/professors/{prof_id}", handlers.GetProfessorHandler).Methods("GET")
	// router.HandleFunc("/professors/{prof_id}", handlers.UpdateProfessorHandler).Methods("PUT")
	// router.HandleFunc("/professors/{prof_id}", handlers.DeleteProfessorHandler).Methods("DELETE")

	// Set up routes for the school
	router.HandleFunc("/schools", handlers.GetAllSchools).Methods("GET")
	router.Handle("/schools/{owner_id}", auth.AuthRequired(http.HandlerFunc(handlers.GetAllSchoolsForUser))).Methods("GET")
	router.HandleFunc("/schools/school/{school_id}", handlers.GetSchoolHandler).Methods("GET")
	router.Handle("/schools/school/{school_id}", auth.AuthRequired(http.HandlerFunc(handlers.UpdateSchoolHandler))).Methods("PUT")
	router.Handle("/schools/school/{school_id}/delete", auth.AuthRequired(http.HandlerFunc(handlers.DeleteSchoolHandler))).Methods("PUT")
	router.HandleFunc("/schools/school/{school_id}/students", handlers.GetStudentsForSchoolHandler).Methods("GET")
	router.HandleFunc("/schools/school/{school_id}/classes", handlers.GetClassesForSchoolHandler).Methods("GET")

	// telemetry
	// data generation handlers included under telemetry
	router.Handle("/telemetry/{owner_id}", auth.AuthRequired(http.HandlerFunc(handlers.CreateNewSchoolHandler))).Methods("POST")
	router.HandleFunc("/telemetry/creation_status/{operation_id}", handlers.SchoolCreationStatusHandler).Methods("GET")
	router.HandleFunc("/telemetry/best-professors", handlers.GetBestProfessorsHandler).Methods("GET")
	router.HandleFunc("/telemetry/{school_id}/classes/avg_gpa", handlers.GetGradeAvgForSchoolHandler).Methods("GET")
	router.HandleFunc("/telemetry/{school_id}/update/avg_gpa", handlers.UpdateSchoolAvgHandler).Methods("GET")

	// AUTH
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	// might not need logout
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	router.HandleFunc("/users/create_user", handlers.CreateUserHandler).Methods("POST")
	router.Handle("/users/update_user/{owner_id}", auth.AuthRequired(http.HandlerFunc(handlers.UpdateUserHandler))).Methods("PUT")
	router.Handle("/users/delete_user/{owner_id}", auth.AuthRequired(http.HandlerFunc(handlers.DeleteUserHandler))).Methods("DELETE")

	// Create a config file for prod,dev
	// create a loading function below start server
	// no loading func created might make one when things finished tho.
	startServer(router)

}

func addCorsHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		handler.ServeHTTP(w, r)
	})
}
func addCorsHeadersProd(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://www.hiremeresume.com")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			fmt.Println("options branch hit")
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func startServer(handler http.Handler) {

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		// when done this will become dev and prod will only build through pipe
		appEnv = "dev"
		//appEnv = "prod"
	}
	if appEnv == "prod" {
		log.Println("Starting server prod on :5000")
		log.Fatal(http.ListenAndServe(":5000", handler))

	} else {
		log.Println("Starting server local on :3000")
		log.Fatal(http.ListenAndServe(":3000", handler))

	}
}
