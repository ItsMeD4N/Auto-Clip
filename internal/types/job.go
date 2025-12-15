package types

import "time"

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusDownloading JobStatus = "downloading"
	StatusAnalyzing  JobStatus = "analyzing"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID        string    `json:"id"`
	VideoURL  string    `json:"video_url"`
	Status    JobStatus `json:"status"`
	Message   string    `json:"message"` // Progress info or error message
	FilePath  string    `json:"-"`       // internal path to final file
	CreatedAt time.Time `json:"created_at"`
}
