package validator

type Record interface {
	GetSubRecords() SubRecord
}

type SubRecord interface {
}
