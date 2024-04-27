package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/dgmann/document-manager/m1-helper/bdt"
)

func Run(ctx context.Context, fileName string, port string) {
	srv := &http.Server{Addr: ":" + port}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("BDT file path: " + fileName + "\nCtrl + Alt + P")); err != nil {
			log.Printf(err.Error())
		}
	})

	http.HandleFunc("/patient", func(w http.ResponseWriter, r *http.Request) {
		j := json.NewEncoder(w)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		f, err := os.Open(fileName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			j.Encode(err)
			return
		}
		defer f.Close()

		patient, err := bdt.Parse(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			j.Encode(err)
			return
		}
		j.Encode(patient)
	})
	go func() {
		select {
		case <-ctx.Done():
			srv.Shutdown(context.Background())
			return
		}
	}()
	srv.ListenAndServe()
}
