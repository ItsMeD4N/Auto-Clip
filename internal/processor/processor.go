package processor

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProcessClip crops video and burns subtitles
func ProcessClip(videoPath, subPath, outputPath string, start, end string) error {
	// FFMPEG command:
	// ffmpeg -i video.mp4 -ss start -to end -vf "crop=ih*(9/16):ih,subtitles=subs.vtt" -c:a copy output.mp4
	// Note: burning subtitles requires re-encoding video.
	// crop=w:h:x:y. We want 9:16 aspect ratio.
	// Width = Height * (9/16).
	// We center it: x = (in_w - out_w)/2
	
	// Filter complex construction
	vf := "crop=ih*(9/16):ih:(iw-ow)/2:0"
	
	if subPath != "" {
		// FFMpeg requires escaping paths for the subtitles filter
		// Windows path escaping is tricky in FFmpeg.
		// safest is to use relative path if possible or forward slashes.
		// simple escape: replace \ with / and escape :
		
		escapedSub := filepath.ToSlash(subPath)
		// escapedSub = strings.ReplaceAll(escapedSub, ":", "\\:") // FFmpeg filter escaping
		// Actually, simpler to copy vtt to a temp file 'subs.vtt' in current dir?
		// Or using "subtitles='path'"
		
		// For windows, the path logic in filter string is messy. 
		// "subtitles='C\:/path/to/file.vtt'"
		
		escapedSub = strings.ReplaceAll(escapedSub, ":", "\\\\:")
		vf += fmt.Sprintf(",subtitles='%s'", escapedSub)
	}

	// -y to overwrite
	args := []string{
		"-y",
		"-i", videoPath,
		"-ss", start,
		"-to", end,
		"-vf", vf,
		"-c:v", "libx264", // Ensure h264
		"-c:a", "aac",     // Ensure aac
		"-strict", "experimental",
		outputPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %s, output: %s", err, string(output))
	}

	return nil
}
