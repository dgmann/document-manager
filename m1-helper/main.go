package main

import (
	"net/http"
	"encoding/json"
	"github.com/dgmann/document-manager/m1-helper/m1"
	"io/ioutil"
	"flag"
)

func main() {
	fileName := flag.String("f", "", "BDT file containing current patient")
	flag.Parse()
	if *fileName == "" {
		panic("no file provided")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("BDT file path: " + *fileName))
	})
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
	})
	http.HandleFunc("/patient", func(w http.ResponseWriter, r *http.Request) {
		j := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		f, err := ioutil.ReadFile(*fileName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			j.Encode(err)
		}

		patient, err := m1.Parse(toUtf8(f))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			j.Encode(err)
		}
		j.Encode(patient)
	})
	http.ListenAndServe(":3000", nil)
}

func toUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}
