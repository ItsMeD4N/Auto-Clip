package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DownloadVideo fetches video and subtitles
// Returns: videoPath, subtitlePath, info, error
func DownloadVideo(url string, outputDir string, jobID string) (string, string, error) {
	// Ensure dir exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", "", err
	}

	// Output template: outputDir/jobID.%(ext)s
	outputTemplate := filepath.Join(outputDir, jobID+".%(ext)s")

	// Command: yt-dlp -f "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best" --write-auto-sub --sub-lang en --skip-download --print-json URL
	// Actually we need to download it.
	// We want mp4.
	// We want subtitles (auto-generated or manual).
	// We'll prioritize English.
	
	// Flags:
	// -f mp4: ensure mp4 container (or merge).
	// --write-sub --write-auto-sub: get subs.
	// --sub-lang en: english.
	// --output ...
	
	cmd := exec.Command("yt-dlp",
		"-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best",
		"--write-auto-sub",
		"--sub-lang", "en",
		"--output", outputTemplate,
		url,
	)
	
	// Run command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("yt-dlp failed: %s, output: %s", err, string(output))
	}

	// Find the file. yt-dlp might merge files.
	// Expected files: jobID.mp4, jobID.en.vtt
	
	videoPath := filepath.Join(outputDir, jobID+".mp4")
	subPath := filepath.Join(outputDir, jobID+".en.vtt")
	
	// Check if file exists
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		// Try to find any mp4 starting with jobID
		files, _ := filepath.Glob(filepath.Join(outputDir, jobID+"*.mp4"))
		if len(files) > 0 {
			videoPath = files[0]
		} else {
			return "", "", fmt.Errorf("video file not found after download")
		}
	}

	// Check for subtitles
	if _, err := os.Stat(subPath); os.IsNotExist(err) {
		// Try glob
		files, _ := filepath.Glob(filepath.Join(outputDir, jobID+"*.vtt"))
		if len(files) > 0 {
			subPath = files[0]
		} else {
			// Ensure we have a sub path even if empty, but processor might fail.
			// Making it optional? The prompt says "Burns subtitles".
			// We'll allow returning empty string and handle it.
			subPath = ""
		}
	}

	return videoPath, subPath, nil
}
