package types

import "time"

type Todo struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    string    `json:"priority"`
	Complete    bool      `json:"complete"`
}
