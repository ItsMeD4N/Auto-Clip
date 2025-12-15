package queue

import (
	"auto-clip/internal/types"
	"sync"
	"time"
)

type JobQueue struct {
	jobs map[string]*types.Job
	mu   sync.RWMutex
}

var instance *JobQueue
var once sync.Once

func GetQueue() *JobQueue {
	once.Do(func() {
		instance = &JobQueue{
			jobs: make(map[string]*types.Job),
		}
	})
	return instance
}

func (q *JobQueue) AddJob(id string, url string) *types.Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	job := &types.Job{
		ID:        id,
		VideoURL:  url,
		Status:    types.StatusPending,
		CreatedAt: time.Now(),
		Message:   "Queued",
	}
	q.jobs[id] = job
	return job
}

func (q *JobQueue) GetJob(id string) (*types.Job, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	job, exists := q.jobs[id]
	return job, exists
}

func (q *JobQueue) UpdateStatus(id string, status types.JobStatus, message string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if job, exists := q.jobs[id]; exists {
		job.Status = status
		job.Message = message
	}
}

func (q *JobQueue) SetFilePath(id string, path string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	if job, exists := q.jobs[id]; exists {
		job.FilePath = path
	}
}
