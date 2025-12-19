package domain

type State struct {
	Accumulator int64   `json:"accumulator"`
	Values      []int64 `json:"values"`
}

type BackgroundService interface {
	FetchData() (int64, []int64, error)
	Run() (bool, error)
}

type BackgroundRepository interface {
	FetchData() (State, error)
}
