# Auto-Clip 🎬

Turn YouTube Videos into TikTok-ready clips automatically.

## Tech Stack
- **Backend**: Go (Gin)
- **Frontend**: React + TypeScript (Vite)
- **Tools**: `yt-dlp` (Download), `ffmpeg` (Processing)

## Prerequisites
Ensure the following tools are installed and in your system PATH:
1.  **Go** (1.20+)
2.  **Node.js** (18+)
3.  **yt-dlp**: [Installation Guide](https://github.com/yt-dlp/yt-dlp#installation)
4.  **FFmpeg**: [Download](https://ffmpeg.org/download.html)

## Getting Started

### 1. Start Backend
```bash
go run cmd/server/main.go
```
Server runs on `http://localhost:8080`

### 2. Start Frontend
```bash
cd frontend
npm install
npm run dev
```
Access UI at `http://localhost:5173`

## How it works
1.  **Download**: Fetches video & subtitles using `yt-dlp`.
2.  **Analyze**: Simple heuristic (or AI) picks a 30s segment.
3.  **Process**: FFmpeg crops to 9:16 and burns subtitles.
4.  **Download**: File is ready for download.
