package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/dgmann/document-manager/pkg/api"
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
	fs := flag.NewFlagSet("download", flag.ExitOnError)
	all := fs.Bool("a", false, "download all files")
	fs.Parse(args)
	if *all {
		return r.downloadAll(args[:1])
	}
	return r.downloadSingle(args[:1])
}

func (r *record) downloadAll(args []string) error {
	subargs := args[:1]
	fs := flag.NewFlagSet("download-all", flag.ExitOnError)
	output := fs.String("d", "", "output folder")
	if err := fs.Parse(subargs); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}
	if *output == "" {
		return errors.New("-d parameter is required")
	}
	records, err := dmClient.Records.List()
	if err != nil {
		return err
	}
	slog.Info("record list fetched", "count", len(records))
	downloadList := make(chan api.Record)
	go func() {
		for record := range downloadList {
			var data io.ReadCloser
			for data, err = dmClient.Records.Download(record.Id); err != nil; {
				slog.Error("error downloading PDF", "recordId", record.Id)
			}
			var destPath string
			if *record.Status == api.StatusInbox {
				destPath = path.Join(*output, "inbox", record.ReceivedAt.Format("2006-01-02")+".pdf")
			} else {
				destPath = path.Join(*output, *record.PatientId, *record.Category, record.Date.Format("2006-01-02")+".pdf")
			}
			if err := os.MkdirAll(path.Dir(destPath), os.ModePerm); err != nil {
				slog.Error("error creating directory", "error", err)
				os.Exit(1)
			}
			out, err := os.Create(destPath)
			if err != nil {
				slog.Error("error creating file", "error", err)
				os.Exit(1)
			}
			if _, err := io.Copy(out, data); err != nil {
				slog.Error("error writing file", "error", err)
				os.Exit(1)
			}
		}
	}()
	for _, record := range records {
		downloadList <- record
	}
	close(downloadList)
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
