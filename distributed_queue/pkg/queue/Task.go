package queue

import "time"

type Task struct {
	id          int
	workerid    int
	status      string
	data        []byte
	createdTime time.Time
}
