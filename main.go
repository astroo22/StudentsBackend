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
	// Create the HTTP server and set the router
	//router := http.NewServeMux()
	router := mux.NewRouter()

	// Set up the routes for Students
	router.HandleFunc("/students", handlers.CreateStudentHandler).Methods("POST")
	router.HandleFunc("/students/{id}", handlers.GetStudentHandler).Methods("GET")
	router.HandleFunc("/students/{id}", handlers.UpdateStudentHandler).Methods("PUT")
	router.HandleFunc("/students/{id}", handlers.DeleteStudentHandler).Methods("DELETE")

	// Set up the routes for Grades
	// router.HandleFunc("/grades", handlers.CreateGradeHandler).Methods("POST")
	// router.HandleFunc("/grades/{id}", handlers.GetGradeHandler).Methods("GET")
	// router.HandleFunc("/grades/{id}", handlers.UpdateGradeHandler).Methods("PUT")
	// router.HandleFunc("/grades/{id}", handlers.DeleteGradeHandler).Methods("DELETE")

	// Set up the routes for Classes
	router.HandleFunc("/classes", handlers.CreateClassHandler).Methods("POST")
	router.HandleFunc("/classes/{id}", handlers.GetClassHandler).Methods("GET")
	router.HandleFunc("/classes/{id}", handlers.UpdateClassHandler).Methods("PUT")
	router.HandleFunc("/classes/{id}", handlers.DeleteClassHandler).Methods("DELETE")

	// Set up the routes for Professors
	// router.HandleFunc("/professors", handlers.CreateProfessorHandler).Methods("POST")
	// router.HandleFunc("/professors/{id}", handlers.GetProfessorHandler).Methods("GET")
	// router.HandleFunc("/professors/{id}", handlers.UpdateProfessorHandler).Methods("PUT")
	// router.HandleFunc("/professors/{id}", handlers.DeleteProfessorHandler).Methods("DELETE")

	// Start the HTTP server
	log.Printf("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}
