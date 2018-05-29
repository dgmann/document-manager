package record

type FileSystemIndex struct {
	data              []*Record
	TotalRecordCount  int
	TotalPatientCount int
}

func NewFileSystemIndex(data []*Record) *FileSystemIndex {
	patients := make(map[int]struct{})
	for _, r := range data {
		patients[r.PatId] = struct{}{}
	}
	return &FileSystemIndex{
		data:              data,
		TotalRecordCount:  len(data),
		TotalPatientCount: len(patients),
	}
}
