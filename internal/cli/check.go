package cli

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/dgmann/document-manager/pkg/api"
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
	default:
		l.Fatalf("error: unknown command - %q\n", args[0])
		return ErrUnknownCommand
	}
}

func (r *check) patients(args []string) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	parallel := fs.Int("p", 4, "number of parallel workers")
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

	bar := progressbar.NewOptions(len(records),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("checking records"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(true),
	)

	recordChan := make(chan api.Record)
	unknownPatientRecordsChan := make(chan api.Record)

	var wg sync.WaitGroup
	for range *parallel {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range recordChan {
				if _, ok := patientMap[*record.PatientId]; !ok {
					unknownPatientRecordsChan <- record
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
		close(unknownPatientRecordsChan)
	}()

	var unknownPatientRecords []api.Record
	for record := range unknownPatientRecordsChan {
		unknownPatientRecords = append(unknownPatientRecords, record)
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
