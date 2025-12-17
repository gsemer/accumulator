package domain

type Job interface {
	Process() error
}

type AddJob struct {
	State *State
	Value int64
}

func (j *AddJob) Process() error {
	err := j.State.Add(j.Value)
	return err
}
