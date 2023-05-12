package main

import (
	"fmt"
	"log"
	"net/http"
	"students/handlers"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("HEREEE WE GOOOOOO")
	// var(
	// 	lastUpdate = time.Now()
	// 	interval = 10 * time.Minute
	// )

	// Create the HTTP server and set the router
	//router := http.NewServeMux()
	router := mux.NewRouter()

	// Set up the routes for Students
	router.HandleFunc("/students", handlers.CreateStudentHandler).Methods("POST")
	router.HandleFunc("/students", handlers.GetAllStudentsHandler).Methods("GET")
	router.HandleFunc("/students/{student_id}", handlers.GetStudentHandler).Methods("GET")
	router.HandleFunc("/students/{student_id}", handlers.UpdateStudentHandler).Methods("PUT")
	router.HandleFunc("/students/{student_id}", handlers.DeleteStudentHandler).Methods("DELETE")

	// Set up the routes for ReportCards
	router.HandleFunc("/reportcard", handlers.CreateReportCardHandler).Methods("POST")
	router.HandleFunc("/reportcard/{student_id}", handlers.GetReportCardHandler).Methods("GET")
	router.HandleFunc("/reportcard/{student_id}", handlers.UpdateReportCardHandler).Methods("PUT")
	router.HandleFunc("/reportcard/{student_id}", handlers.DeleteReportCardHandler).Methods("DELETE")

	// Set up the routes for Classes
	router.HandleFunc("/classes", handlers.CreateClassHandler).Methods("POST")
	router.HandleFunc("/classes/{class_id}", handlers.GetClassHandler).Methods("GET")
	router.HandleFunc("/classes/{class_id}", handlers.UpdateClassHandler).Methods("PUT")
	router.HandleFunc("/classes/{class_id}", handlers.DeleteClassHandler).Methods("DELETE")

	// Set up the routes for Professors
	router.HandleFunc("/professors", handlers.CreateProfessorHandler).Methods("POST")
	router.HandleFunc("/professors/{prof_id}", handlers.GetProfessorHandler).Methods("GET")
	router.HandleFunc("/professors/{prof_id}", handlers.UpdateProfessorHandler).Methods("PUT")
	router.HandleFunc("/professors/{prof_id}", handlers.DeleteProfessorHandler).Methods("DELETE")

	// Set up routes for the school
	router.HandleFunc("/schools", handlers.GetAllSchools).Methods("GET")
	router.HandleFunc("/schools/{school_id}", handlers.GetSchoolHandler).Methods("GET")
	router.HandleFunc("/schools/{school_id}", handlers.UpdateSchoolHandler).Methods("PUT")
	router.HandleFunc("/schools/{school_id}", handlers.DeleteSchoolHandler).Methods("DELETE")
	router.HandleFunc("/schools/{school_id}/classes", handlers.GetClassesForSchoolHandler).Methods("GET")

	// data generation handlers included under telemetry
	router.HandleFunc("/telemetry/{owner_id}", handlers.CreateNewSchoolHandler).Methods("POST")

	// telemetry
	router.HandleFunc("/telemetry/best-professors", handlers.GetBestProfessorsHandler).Methods("GET")
	router.HandleFunc("/telemetry/{school_id}/classes/avg_gpa", handlers.GetGradeAvgForSchoolHandler).Methods("GET")

	// AUTH
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	// Start the HTTP server on port 3000
	log.Printf("Starting server on port 3000")
	if err := http.ListenAndServe(":3000", addCorsHeaders(router)); err != nil {
		log.Fatal(err)
	}

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
