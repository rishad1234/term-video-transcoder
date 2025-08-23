# Terminal Video Transcoder

A fast, efficient command-line tool for video transcoding, media analysis, and batch processing, powered by FFmpeg and FFprobe.

## Features

✅ **Phase 1 (Current):**

- Media file analysis with detailed information display
- FFmpeg/FFprobe integration and validation

🚧 **Coming Soon:**

- Video format conversion
- Quality presets (low, medium, high)
- Batch processing
- Progress tracking

## Installation

### Prerequisites

- Go 1.21+
- FFmpeg and FFprobe installed

### Build from Source

```bash
git clone https://github.com/rishad1234/term-video-transcoder.git
cd term-video-transcoder
go build -o transcoder .
```

## Usage

### Media Analysis

```bash
# Display detailed information about a media file
./transcoder info video.mp4
./transcoder info movie.mkv
```

### Help

```bash
# General help
./transcoder --help

# Command-specific help
./transcoder info --help
```

## Dependencies

This project uses the following Go libraries:

- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal colors
- `github.com/schollz/progressbar/v3` - Progress bars
- `github.com/tidwall/gjson` - JSON parsing

## Development Status

Currently in **Phase 1** development. See `initial_planning.md` for the complete roadmap.

## License

[To be determined]
