package cli

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/dgmann/document-manager/pkg/api"
	"github.com/dgmann/document-manager/pkg/log"
	"github.com/schollz/progressbar/v3"
)

type check struct{}

func CheckCmd() *check {
	return &check{}
}

func (r *check) Execute(args []string) error {
	subArgs := args[1:]
	switch args[0] {
	case "patients":
		return r.patients(subArgs)
	case "pages":
		return r.pages(subArgs)
	default:
		l.Fatalf("error: unknown command - %q\n", args[0])
		return ErrUnknownCommand
	}
}

func (r *check) patients(args []string) error {
	slog.Info("fetching list of all patients. This may take some time...")
	patients, err := dmClient.Patients.All()
	if err != nil {
		return fmt.Errorf("error fetching patients: %w", err)
	}
	patientMap := make(map[string]api.Patient, len(patients))
	for _, patient := range patients {
		patientMap[patient.Id] = patient
	}

	slog.Info("fetching list of all records. This may take some time...")
	records, err := dmClient.Records.List()
	if err != nil {
		return err
	}
	slog.Info("record list fetched", "count", len(records))
	slog.Info("checking patients")

	var unknownPatientRecords []api.Record
	for _, record := range records {
		if _, ok := patientMap[*record.PatientId]; !ok {
			unknownPatientRecords = append(unknownPatientRecords, record)
		}
	}

	sort.Slice(unknownPatientRecords, func(i, j int) bool {
		return *unknownPatientRecords[i].PatientId < *unknownPatientRecords[j].PatientId
	})

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Patient\tRecord\t")
	for _, record := range unknownPatientRecords {
		fmt.Fprintf(w, "%s\t%s\t\n", *record.PatientId, record.Id)
	}
	return nil
}

func (r *check) pages(args []string) error {
	fs := flag.NewFlagSet("check-pages", flag.ContinueOnError)
	parallel := fs.Int("p", 4, "number of parallel workers")
	slog.Info("fetching list of all records. This may take some time...")
	records, err := dmClient.Records.List()
	if err != nil {
		return err
	}
	slog.Info("record list fetched", "count", len(records))

	bar := progressbar.NewOptions(len(records),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("checking records"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(true),
	)

	recordChan := make(chan api.Record)
	invalidPagesChan := make(chan api.Record)

	var wg sync.WaitGroup
	for range *parallel {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range recordChan {
				for _, page := range record.Pages {
					var resp *http.Response
					for {
						resp, err = http.Get(page.Url)
						if err == nil {
							break
						}
						slog.Debug("error fetching page", "recordId", record.Id, log.ErrAttr(err))
					}
					if resp.StatusCode != 200 {
						invalidPagesChan <- record
						break
					}
				}
				bar.Add(1)
			}
		}()
	}

	go func() {
		for _, record := range records {
			recordChan <- record
		}
		close(recordChan)
		wg.Wait()
		close(invalidPagesChan)
	}()

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "Patient\tRecord\t")
	for record := range invalidPagesChan {
		fmt.Fprintf(w, "%s\t%s\t\n", *record.PatientId, record.Id)
	}

	return nil
}
