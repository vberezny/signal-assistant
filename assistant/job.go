package assistant

// TODO: some way to track when the job is run.
type Job struct {
	Name string
}

func NewJob() *Job {
	return &Job{}
}
