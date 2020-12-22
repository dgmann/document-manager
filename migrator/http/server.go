package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/dgmann/document-manager/migrator/importer"
	"github.com/dgmann/document-manager/migrator/records/databasereader"
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/validator"
	"golang.org/x/sync/semaphore"
)

type Server struct {
	DatabaseManager   *databasereader.Manager
	FilesystemManager *filesystem.Manager
	ImportManager     *importer.Manager
	State             State
	Config            Config
}

type State struct {
	ImportRunning *semaphore.Weighted
	Resolvables   []validator.ResolvableValidationError
}

func NewServer(conf Config) (*Server, error) {
	recordManager := databasereader.NewManager(conf.DbName, conf.Username, conf.Password, conf.Hostname, conf.Instance)
	err := recordManager.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}
	fileystemManager := filesystem.NewManager(conf.RecordDirectory, conf.DataDirectory)
	importManager := importer.NewManager(fileystemManager, conf.DataDirectory, recordManager.Db, conf.ApiURL, conf.RetryCount)
	return &Server{DatabaseManager: recordManager, FilesystemManager: fileystemManager, ImportManager: importManager, Config: conf, State: State{ImportRunning: semaphore.NewWeighted(1)}}, nil
}

func (s *Server) Run(port string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("index.gohtml").ParseFiles("web/template/index.gohtml"))
		if err := t.Execute(w, s); err != nil {
			w.Write([]byte(err.Error()))
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	http.HandleFunc("/database/counts", func(w http.ResponseWriter, r *http.Request) {
		index, err := s.DatabaseManager.Index()
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		returnCounts(w, index)
	})

	http.HandleFunc("/filesystem/counts", func(w http.ResponseWriter, r *http.Request) {
		index, err := s.FilesystemManager.Index(r.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		returnCounts(w, index)
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		filesystemIndex, err := s.FilesystemManager.Index(r.Context())
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		databaseIndex, err := s.DatabaseManager.Index()
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		resolvable, validationErrors := validator.Validate(r.Context(), filesystemIndex, databaseIndex, s.DatabaseManager.Manager)
		s.State.Resolvables = resolvable
		records := make([]*models.Record, len(resolvable), len(resolvable))
		for i, res := range resolvable {
			records[i] = res.Record()
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"resolvables": records,
			"errors":      validationErrors.Messages,
		})
	})
	http.HandleFunc("/validate/resolve", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"resolvableCount": len(s.State.Resolvables),
			})
		} else if r.Method == http.MethodPost {
			errors := make([]string, 0)
			var remaining []validator.ResolvableValidationError
			for _, resolvable := range s.State.Resolvables {
				if err := resolvable.Resolve(); err != nil {
					errors = append(errors, err.Error())
					remaining = append(remaining, resolvable)
				}
			}
			s.State.Resolvables = remaining
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			filename := strings.TrimPrefix(r.URL.Path, "/files/")
			if filename == "" {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(s.ImportManager.Files())
			} else {
				importManager := importer.NewManager(s.FilesystemManager, s.Config.DataDirectory, s.DatabaseManager.Db, s.Config.ApiURL, s.Config.RetryCount)
				if err := importManager.Load(filename); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, err.Error())
					return
				}
				importable, err := importManager.DataToImport(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, err.Error())
					return
				}
				w.Header().Set("Content-Type", "application/json")
				records := importable.Records
				if len(records) > 500 {
					records = records[:500]
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"total":     len(importable.Records),
					"records":   records,
					"truncated": len(importable.Records) > 500,
				})
			}

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/import", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method == http.MethodGet {
			if !s.ImportManager.IsLoaded() {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"total":      0,
					"imported":   0,
					"errors":     nil,
					"categories": 0,
					"hasData":    false,
				})
				return
			}

			importable, err := s.ImportManager.DataToImport(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"total":      len(importable.Records),
				"imported":   len(s.ImportManager.ImportedRecords()),
				"errors":     s.ImportManager.ImportErrors,
				"categories": len(importable.Categories),
				"hasData":    true,
			})
		} else if r.Method == http.MethodPut {
			if ok := s.State.ImportRunning.TryAcquire(1); !ok {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "import already running")
				return
			}
			defer s.State.ImportRunning.Release(1)

			if file := r.URL.Query().Get("file"); file != "" {
				if err := s.ImportManager.Load(file); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, err.Error())
					return
				}
			}

			if err := s.ImportManager.Import(r.Context()); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err.Error())
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	})
	return http.ListenAndServe(":"+port, nil)
}

func returnCounts(w http.ResponseWriter, countable models.Countable) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{
		"records":  countable.GetTotalRecordCount(),
		"patients": countable.GetTotalPatientCount(),
	})
}
