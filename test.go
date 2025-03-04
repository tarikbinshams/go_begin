package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func test() {
	r := mux.NewRouter()

	r.HandleFunc("/hello", HandleHello)

	r.HandleFunc("/register", HandleRegister).Methods("POST")

	http.ListenAndServe(":4000", r)

	log.Println("Server started on :4000")
}

func HandleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

type RegisterForm struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (f *RegisterForm) Validate() error {
	if f.Name == "" || f.Email == "" || f.Password == "" {
		return fmt.Errorf("all fields are required")
	}
	return nil
}
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var form RegisterForm

	if contentType == "application/json" {
		// Handle JSON
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &form); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		log.Println("Parsed JSON:", form)
	} else {
		// Handle form data (multipart/form-data or application/x-www-form-urlencoded)
		if strings.HasPrefix(contentType, "multipart/form-data") {
			// Parse multipart form with a reasonable memory limit (e.g., 10MB)
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				log.Println("Error parsing multipart form:", err)
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}
		} else {
			// Parse URL-encoded form
			if err := r.ParseForm(); err != nil {
				log.Println("Error parsing form:", err)
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}
		}

		// Retrieve form values
		form.Name = r.PostFormValue("name")
		form.Email = r.PostFormValue("email")
		form.Password = r.PostFormValue("password")

		fmt.Println("Parsed Form Data---------------------------------")
		fmt.Println("Name:", form.Name)
		fmt.Println("Email:", form.Email)
		fmt.Println("Password:", form.Password)
	}

	if err := form.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registration successful"))
}
