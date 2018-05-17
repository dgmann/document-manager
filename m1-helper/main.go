package main

import (
	"net/http"
	"encoding/json"
	"github.com/dgmann/document-manager/m1-helper/m1"
	"io/ioutil"
	"flag"
	"github.com/dgmann/document-manager/m1-helper/hotkey"
	"os/exec"
	"fmt"
	"log"
)

func main() {
	fileName := flag.String("f", "", "BDT file containing current patient")
	serverUrl := flag.String("s", "http://localhost", "Document-Manager URL")
	flag.Parse()
	if *fileName == "" {
		panic("no file provided")
	}

	registerHotKey(*fileName, *serverUrl)
	startHttpServer(*fileName)
}

func toUtf8(iso88591Buf []byte) string {
	buf := make([]rune, len(iso88591Buf))
	for i, b := range iso88591Buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func registerHotKey(fileName, serverUrl string) {
	manager := hotkey.NewManager()
	manager.Register(hotkey.Hotkey{1, hotkey.ModAlt + hotkey.ModCtrl, 'P'})
	keyPresses := manager.Listen()
	go func() {
		for range keyPresses {
			f, err := ioutil.ReadFile(fileName)
			if err != nil {
				println("error reading patient file")
			}

			patient, err := m1.Parse(toUtf8(f))

			cmd := exec.Command(fmt.Sprintf("%s/patient/%s", serverUrl, patient.Id))
			err = cmd.Run()

			if err != nil {
				fmt.Printf("an error occurred: %s\n", err)
				log.Fatal(err)
			}
		}
	}()
}

func startHttpServer(fileName string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("BDT file path: " + fileName))
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

		patient, err := m1.Parse(toUtf8(f))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			j.Encode(err)
		}
		j.Encode(patient)
	})
	http.ListenAndServe(":3000", nil)
}
