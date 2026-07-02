package task

type Task interface {
	Run(args ...interface{}) Completed
}

type Completed struct {
	Err error
}
