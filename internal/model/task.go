package model

type Task struct {
	Status  string `json:"status"`
	Seconds int64  `json:"seconds"`
}
