package types

type TaskStatus string

const (
	TaskStatusIncomplete TaskStatus = "incomplete"
	TaskStatusComplete   TaskStatus = "complete"
	TaskStatusDecayed    TaskStatus = "decayed"
)
