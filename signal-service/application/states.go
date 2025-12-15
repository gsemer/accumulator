package application

import (
	"errors"
	"log"
	"signal/domain"
)

type BackgroundService struct {
	br domain.BackgroundRepository
}

func NewBackgroundService(br domain.BackgroundRepository) *BackgroundService {
	return &BackgroundService{br: br}
}

func (bs *BackgroundService) FetchData() (int64, []int64, error) {
	accumulator, err := bs.br.FetchAccumulator()
	if err != nil {
		return 0, nil, err
	}

	values, err := bs.br.FetchValues()
	if err != nil {
		return 0, nil, err
	}

	return accumulator, values, nil
}

func (bs *BackgroundService) Run() (bool, error) {
	accumulator, values, err := bs.FetchData()
	if err != nil {
		return false, err
	}

	var sum int64
	for _, value := range values {
		sum += value
	}

	if accumulator != sum {
		log.Printf("expected %v, got %v", accumulator, sum)
		return false, errors.New("the sum of the values is not equal to the accumulator")
	}
	log.Println("the sum of the values is equal to the accumulator as expected")
	return true, nil
}
