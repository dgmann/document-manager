package cli

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"sync"
	"time"

	"github.com/dgmann/document-manager/pkg/api"
	"github.com/schollz/progressbar/v3"
)

type record struct{}

func Record() *record {
	return &record{}
}

var l = log.New(os.Stderr, "", 1)

func (r *record) Execute(args []string) error {
	subArgs := args[1:]
	switch args[0] {
	case "download":
		return r.download(subArgs)
	case "show":
		return r.show(subArgs)
	default:
		l.Fatalf("error: unknown command - %q\n", args[0])
		return errors.New("unkown command")
	}
}

func (r *record) download(args []string) error {
	fs := flag.NewFlagSet("download", flag.ContinueOnError)
	all := fs.Bool("a", false, "download all files")
	fs.Parse(args[:1])
	if *all {
		return r.downloadAll(args[1:])
	}
	return r.downloadSingle(args[1:])
}

type downloadedRecord struct {
	Path   string
	Record api.Record
}

type recordError struct {
	record api.Record
	err    error
}

func (r *record) downloadAll(args []string) error {
	fs := flag.NewFlagSet("download-all", flag.ExitOnError)
	output := fs.String("d", "", "output folder")
	parallel := fs.Int("p", 4, "number of parallel workers")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}
	if *output == "" {
		return errors.New("-d parameter is required")
	}

	if err := os.MkdirAll(*output, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	indexFile, err := os.Create(path.Join(*output, "index.csv"))
	if err != nil {
		return fmt.Errorf("error creating index file: %w", err)
	}

	slog.Info("fetching list of all records. This may take some time...")
	records, err := dmClient.Records.List()
	if err != nil {
		return err
	}
	slog.Info("record list fetched", "count", len(records))

	slog.Info("fetching list of all patients. This may take some time...")
	patients, err := dmClient.Patients.All()
	if err != nil {
		return fmt.Errorf("error fetching patients: %w", err)
	}
	patientMap := make(map[string]api.Patient, len(patients))
	for _, patient := range patients {
		patientMap[patient.Id] = patient
	}

	slog.Info("fetching list of all categories. This may take some time...")
	categories, err := dmClient.Categories.All()
	if err != nil {
		return fmt.Errorf("error fetching categories: %w", err)
	}
	categoryMap := make(map[string]api.Category, len(categories))
	for _, category := range categories {
		categoryMap[category.Id] = category
	}

	bar := progressbar.NewOptions(len(records),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("downloading records"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(true),
	)

	downloadList := make(chan api.Record)
	downloaded := make(chan downloadedRecord)
	errorChan := make(chan recordError)

	var wg sync.WaitGroup
	for range *parallel {
		wg.Add(1)
		go func() {
			for record := range downloadList {
				var data io.ReadCloser
				if data, err = dmClient.Records.Download(record.Id); err != nil {
					slog.Error("error downloading PDF", "recordId", record.Id)
					errorChan <- recordError{record, fmt.Errorf("error downloading PDF")}
					return
				}
				var destPath string
				if *record.Status == api.StatusInbox {
					destPath = path.Join(*output, "inbox", record.ReceivedAt.Format("2006-01-02")+".pdf")
				} else {
					destPath = path.Join(*output, *record.PatientId, *record.Category, record.Date.Format("2006-01-02")+".pdf")
				}
				if err := os.MkdirAll(path.Dir(destPath), os.ModePerm); err != nil {
					errorChan <- recordError{record, fmt.Errorf("error creating directory: %w", err)}
					return
				}
				out, err := os.Create(destPath)
				if err != nil {
					errorChan <- recordError{record, fmt.Errorf("error creating file: %w", err)}
					return
				}
				if _, err := io.Copy(out, data); err != nil {
					errorChan <- recordError{record, fmt.Errorf("error writing file: %w", err)}
					return
				}
				downloaded <- downloadedRecord{Record: record, Path: destPath}
			}
			wg.Done()
		}()
	}

	go func() {
		w := csv.NewWriter(indexFile)
		w.Write([]string{"patient_id", "patient_lastname", "patient_firstname", "patient_birthdate", "category", "path"})
		for d := range downloaded {
			patient, ok := patientMap[*d.Record.PatientId]
			if !ok {
				slog.Warn("unkown patient", "recordId", d.Record.Id, "patientId", *d.Record.PatientId)
				patient = api.Patient{Id: d.Record.Id, FirstName: "unkown", LastName: "unkown"}
			}
			category := categoryMap[*d.Record.Category]
			birthDate := ""
			if patient.BirthDate != nil {
				birthDate = patient.BirthDate.Format(time.RFC3339)
			}
			w.Write([]string{patient.Id, patient.LastName, patient.FirstName, birthDate, category.Name, d.Path})
			bar.Add(1)
		}
		bar.Finish()
	}()

	go func() {
		for err := range errorChan {
			slog.Error("error downloading records", "error", err.err.Error(), "record", err.record.Id)
			os.Exit(1)
		}
	}()

	for _, record := range records {
		downloadList <- record
	}
	wg.Wait()
	close(downloadList)
	close(downloaded)
	close(errorChan)

	slog.Info("download complete")
	return nil
}

func (r *record) downloadSingle(args []string) error {
	recordId := args[0]
	fs := flag.NewFlagSet("download-single", flag.ContinueOnError)
	output := fs.String("d", "./", "output folder")
	if err := fs.Parse(args[:1]); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	destPath := path.Join(*output, recordId+".pdf")
	data, err := dmClient.Records.Download(recordId)
	if err != nil {
		return fmt.Errorf("error downloading PDF: %w", err)
	}
	if err := os.MkdirAll(path.Dir(destPath), os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	if _, err := io.Copy(out, data); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	fmt.Printf("successfully downloaded to %s", destPath)
	return nil
}

func (r *record) show(args []string) error {
	id := args[0]
	record, err := dmClient.Records.Get(id)
	if err != nil {
		return fmt.Errorf("error fetching record %s: %w", id, err)
	}

	prettyJson, err := json.MarshalIndent(record, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling record %s: %w", id, err)
	}
	fmt.Println(string(prettyJson))
	return nil
}
