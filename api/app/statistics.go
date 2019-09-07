package app

type StatisticProvider interface {
	NumberOfElements() (int, error)
}
