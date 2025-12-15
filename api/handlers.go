package api

import (
	"auto-clip/internal/analyzer"
	"auto-clip/internal/downloader"
	"auto-clip/internal/processor"
	"auto-clip/internal/queue"
	"auto-clip/internal/types"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenerateRequest struct {
	URL string `json:"url" binding:"required"`
}

func HandleGenerate(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobID := uuid.New().String()
	q := queue.GetQueue()
	q.AddJob(jobID, req.URL)

	// Trigger processing in background
	go processVideo(jobID, req.URL)

	c.JSON(http.StatusOK, gin.H{"job_id": jobID})
}

func HandleStatus(c *gin.Context) {
	id := c.Param("id")
	q := queue.GetQueue()
	job, exists := q.GetJob(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func HandleDownload(c *gin.Context) {
	id := c.Param("id")
	q := queue.GetQueue()
	job, exists := q.GetJob(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if job.Status != types.StatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job not completed yet"})
		return
	}

	c.File(job.FilePath) // Serve the file
}

// Background Processor
func processVideo(jobID string, url string) {
	q := queue.GetQueue()
	workDir := "./temp"
	
	q.UpdateStatus(jobID, types.StatusDownloading, "Downloading video...")
	
	// 1. Download
	videoPath, subPath, err := downloader.DownloadVideo(url, workDir, jobID)
	if err != nil {
		q.UpdateStatus(jobID, types.StatusFailed, fmt.Sprintf("Download failed: %v", err))
		return
	}

	// 2. Analyze
	q.UpdateStatus(jobID, types.StatusAnalyzing, "Analyzing content...")
	start, end, err := analyzer.AnalyzeTranscript(subPath)
	if err != nil {
		// Log but maybe continue with default?
		fmt.Printf("Analysis failed, using default: %v\n", err)
		start = "00:00:10"
		end = "00:00:40"
	}

	// 3. Process
	q.UpdateStatus(jobID, types.StatusProcessing, "Cropping and burning subtitles...")
	outputPath := filepath.Join(workDir, jobID+"_final.mp4")
	
	err = processor.ProcessClip(videoPath, subPath, outputPath, start, end)
	if err != nil {
		q.UpdateStatus(jobID, types.StatusFailed, fmt.Sprintf("Processing failed: %v", err))
		return
	}

	// Done
	q.SetFilePath(jobID, outputPath)
	q.UpdateStatus(jobID, types.StatusCompleted, "Ready to download")
}
