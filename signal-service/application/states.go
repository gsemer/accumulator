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

func (bs *BackgroundService) FetchData() (domain.State, error) {
	state, err := bs.br.FetchData()
	if err != nil {
		return domain.State{}, err
	}

	return state, nil
}

func (bs *BackgroundService) Run() (bool, error) {
	state, err := bs.FetchData()
	if err != nil {
		return false, err
	}

	var sum int64
	for _, value := range state.Values {
		sum += value
	}

	log.Printf("Accumulator: %d", state.Accumulator)
	log.Printf("Sum: %d", sum)

	if state.Accumulator != sum {
		log.Printf("expected %v, got %v", state.Accumulator, sum)
		return false, errors.New("the sum of the values is not equal to the accumulator")
	}
	log.Println("the sum of the values is equal to the accumulator as expected")
	return true, nil
}
