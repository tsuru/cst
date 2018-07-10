package queue

import (
	"github.com/tsuru/monsterqueue"
	"github.com/tsuru/monsterqueue/mongodb"
)

// ScanTaskName holds the task name to which jobs are assigned.
const ScanTaskName = `scan`

var queueInstance monsterqueue.Queue

// NewQueue creates a new instance of Queue (backed by monsterqueue.MongoDB).
func NewQueue(rawURL string) (monsterqueue.Queue, error) {

	return mongodb.NewQueue(mongodb.QueueConfig{
		Url: rawURL,
	})
}

// SetQueue sets the queue parameter to be used globally by any that calls
// GetQueue function.
func SetQueue(queue monsterqueue.Queue) {
	queueInstance = queue
}

// GetQueue returns the current queue instance. Make sure you set a queue
// instance (calling the SetQueue function) before.
func GetQueue() monsterqueue.Queue {
	return queueInstance
}
