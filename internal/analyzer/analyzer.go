package analyzer

import (
	"os"
	"regexp"
)

// AnalyzeTranscript reads VTT and picks a clip
// In a real app, this would call OpenAI/Gemini
func AnalyzeTranscript(vttPath string) (string, string, error) {
	if vttPath == "" {
		// Fallback if no subs: catch first 30s
		return "00:00:10", "00:00:40", nil
	}

	content, err := os.ReadFile(vttPath)
	if err != nil {
		return "", "", err
	}

	// Parse VTT (Very simplistic)
	// WEBVTT
	//
	// 00:00:00.000 --> 00:00:05.000
	// Some text
	
	// For MVP, we will just take a segment from the middle of the text/video?
	// Or search for specific keywords like "funny", "laugh", etc.
	// But without timestamps-to-line mapping parsing, it's hard.
	
	// Quick parse: find all timestamps
	// Format: \d{2}:\d{2}:\d{2}\.\d{3}
	regex := regexp.MustCompile(`(\d{2}:\d{2}:\d{2}\.\d{3}) --> (\d{2}:\d{2}:\d{2}\.\d{3})`)
	matches := regex.FindAllStringSubmatch(string(content), -1)
	
	if len(matches) < 2 {
		return "00:00:00", "00:00:30", nil
	}

	// Heuristic: Pick a 30s-60s chunk from the middle of the video.
	// The matches are sorted by time.
	middleIndex := len(matches) / 2
	
	start := matches[middleIndex][1]
	
	// find end time ~30s later?
	// Simpler: just use start from middle, and manually add 30s.
	// Or find the timestamp that is closest to start + 30s.
	
	// Let's just return the middle timestamp as start, and let the caller/ffmpeg handle duration?
	// No, we need start and end for -to.
	
	// Assuming the video is long enough.
	// Let's just parse the H:M:S from the string, add 30s, format back.
	// Too complex for 5 mins?
	// Let's "Search" for the match roughly 30 lines later.
	endIndex := middleIndex + 10 // roughly 30s if lines are 3s
	if endIndex >= len(matches) {
		endIndex = len(matches) - 1
	}
	end := matches[endIndex][2]
	
	return start, end, nil
}

// TODO: Implement OpenAI integration here
// func CallLLM(transcript string) ...
