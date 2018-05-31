package shared

type RecordIndex interface {
	GetRecords() []*Record
}
