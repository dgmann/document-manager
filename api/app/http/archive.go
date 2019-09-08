package http

import (
	"io"
	"io/ioutil"
	"net/http"
)

type ArchiveController struct {
	pdfs getter
}

type getter interface {
	Get(id string) (io.Reader, error)
}

func NewArchiveController(pdf getter) *ArchiveController {
	return &ArchiveController{pdfs: pdf}
}

func (a *ArchiveController) One(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	file, err := a.pdfs.Get(id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusNotFound)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/pdf")
	if _, err := w.Write(data); err != nil {
		w.WriteHeader(500)
	}
}
