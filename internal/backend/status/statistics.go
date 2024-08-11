package status

import (
	"context"
	"fmt"
)

type StatisticProvider interface {
	NumberOfElements(ctx context.Context) (int, error)
}

type Providers map[string]StatisticProvider

type StatisticsService struct {
	providers Providers
}

func NewStatisticsService(checkables Providers) *StatisticsService {
	return &StatisticsService{checkables}
}

func (s *StatisticsService) Collect(ctx context.Context) map[string]string {
	messages := make(map[string]string)
	for key, provider := range s.providers {
		numberOfElements, err := provider.NumberOfElements(ctx)
		if err != nil {
			messages[key] = err.Error()
		} else {
			messages[key] = fmt.Sprintf("number of elements: %d", numberOfElements)
		}
	}
	return messages
}
