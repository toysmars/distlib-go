package task

import "time"

// Task represent a task to be done.
type Task struct {
	Message string
	Group   Group
	Option  Option
}

// Option specifies the option for a task.
type Option struct {
	// Priority is the priority of the task. The lower the number, the higher the priority.
	Priority uint32
	// ScheduledFor tells when this taks should be executed.
	ScheduledFor time.Time
}

type Group struct {
	Namespace string
	Name      string
}
