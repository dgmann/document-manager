package server

import (
	"context"
	"encoding/json"
	"github.com/dgmann/document-manager/m1-helper/m1"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func Run(ctx context.Context, fileName string) {
	srv := &http.Server{Addr: ":3000"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("BDT file path: " + fileName)); err != nil {
			logrus.Error(err)
		}
	})

	http.HandleFunc("/patient", func(w http.ResponseWriter, r *http.Request) {
		j := json.NewEncoder(w)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		f, err := ioutil.ReadFile(fileName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			j.Encode(err)
		}

		patient, err := m1.Parse(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			j.Encode(err)
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
