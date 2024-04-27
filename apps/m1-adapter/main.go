package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgmann/document-manager/m1-adapter/m1"
	"log"
	"net/http"
	"os"
	"path"
)

var username string
var password string
var host string
var port string
var dbname string

func init() {
	username = getEnv("DB_USERNAME", "")
	if len(username) == 0 {
		panic("invalid username: " + username)
	}

	password = getEnv("DB_PASSWORD", "")
	if len(password) == 0 {
		panic("invalid password: " + password)
	}

	host = getEnv("DB_HOST", "")
	if len(host) == 0 {
		panic("invalid host: " + host)
	}

	port = getEnv("DB_PORT", "1521")
	if len(port) == 0 {
		panic("invalid port: " + port)
	}

	dbname = getEnv("DB_NAME", "M1DB")
	if len(dbname) == 0 {
		panic("invalid host: " + dbname)
	}
}

func main() {
	dsn := fmt.Sprintf("%s/%s@%s:%s/%s", username, password, host, port, dbname)
	log.Println("Connecting to ", dsn)
	adapter := m1.NewDatabaseAdapter(dsn)
	err := adapter.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	defer adapter.Close()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			writeJSON(writer, response{"error": "page not found"}, http.StatusNotFound)
			return
		}
		writeJSON(writer, response{"message": "M1-Adapter"}, 200)
	})
	http.HandleFunc("/patients", func(writer http.ResponseWriter, request *http.Request) {
		firstname := request.URL.Query().Get("firstname")
		lastname := request.URL.Query().Get("lastname")
		var pats []*m1.Patient
		var err error

		if firstname != "" || lastname != "" {
			pats, err = adapter.FindPatientsByName(firstname, lastname)
		} else {
			pats, err = adapter.GetAllPatients()
		}
		if err != nil {
			writeJSON(writer, response{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSON(writer, pats, 200)
	})
	http.HandleFunc("/patients/", func(writer http.ResponseWriter, request *http.Request) {
		patId := path.Base(request.URL.Path)
		pat, err := adapter.GetPatient(patId)
		if err != nil {
			writeJSON(writer, response{
				"error":   err.Error(),
				"message": "Patient not found",
			}, http.StatusNotFound)
			return
		}
		writeJSON(writer, pat, 200)
	})
	log.Println("m1-adapter running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type response map[string]interface{}

func writeJSON(writer http.ResponseWriter, data interface{}, code int) {
	writer.WriteHeader(code)
	writer.Header().Add("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		writer.WriteHeader(500)
	}
}
