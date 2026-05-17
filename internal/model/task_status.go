package model

type TaskStatus int

const (
	StatusBacklog    TaskStatus = 1
	StatusTodo       TaskStatus = 2
	StatusInProgress TaskStatus = 3
	StatusDone       TaskStatus = 4
	StatusClosed     TaskStatus = 5
)

func (s TaskStatus) String() string {
	switch s {
	case StatusBacklog:
		return "Backlog"
	case StatusTodo:
		return "Todo"
	case StatusInProgress:
		return "InProgress"
	case StatusDone:
		return "Done"
	case StatusClosed:
		return "Closed"
	default:
		return "Unknown"
	}
}

func ParseTaskStatus(v int) TaskStatus {
	switch TaskStatus(v) {
	case StatusBacklog, StatusTodo, StatusInProgress, StatusDone, StatusClosed:
		return TaskStatus(v)
	default:
		return StatusBacklog
	}
}
