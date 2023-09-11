package transfer

import (
	"noda/api/data/types"
)

type TaskCreation struct {
	Title       string             `json:"title"`
	Headline    string             `json:"headline"`
	Description string             `json:"description"`
	Priority    types.TaskPriority `json:"priority"`
	Status      types.TaskStatus   `json:"status"`
}
