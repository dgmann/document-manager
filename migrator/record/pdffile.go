package record

type PdfFile struct {
	path string
}

func NewPdfFile(path string) *PdfFile {
	return &PdfFile{path: path}
}
