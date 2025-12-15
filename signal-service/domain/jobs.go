package domain

type BackgroundService interface {
	FetchData() (int64, []int64, error)
	Run() (bool, error)
}

type BackgroundRepository interface {
	FetchAccumulator() (int64, error)
	FetchValues() ([]int64, error)
}
