# Terminal Video Transcoder

A fast, efficient command-line tool for video transcoding, media analysis, and audio extraction, powered by FFmpeg and FFprobe.

## Features

✅ **Completed Features:**

- **Media Analysis**: Comprehensive file information with detailed metadata
- **Video Conversion**: Format conversion with automatic codec selection
- **Audio Extraction**: Extract audio from videos to multiple formats
- **Quality Presets**: Low, medium, high quality settings
- **Custom Parameters**: Full control over codecs, bitrates, resolution, frame rate
- **Progress Tracking**: Real-time conversion progress with ETA
- **Cross-Platform**: Works on macOS, Linux, and Windows

🎯 **Supported Formats:**

- **Video**: MP4, AVI, MKV, WebM, MOV
- **Audio**: MP3, WAV, AAC, FLAC, OGG, M4A

## Quick Start

```bash
# Analyze a media file
./transcoder info video.mp4

# Convert video with default settings
./transcoder convert input.avi output.mp4

# Extract audio from video
./transcoder extract movie.mkv soundtrack.mp3

# High quality conversion
./transcoder convert input.mov output.webm --preset high

# Custom parameters
./transcoder convert input.mp4 output.webm \
  --video-codec libvpx-vp9 --video-bitrate 2M \
  --audio-codec libopus --audio-bitrate 128k
```

## Documentation

📖 **Complete Manual**: See [MANUAL.md](MANUAL.md) for comprehensive documentation

📋 **Quick Reference**: Run `./transcoder manual` for command-line reference

🆘 **Command Help**: Use `./transcoder [command] --help` for specific help

## Installation

### Prerequisites

- **Go 1.21+** - Programming language and runtime
- **FFmpeg** - Video/audio processing engine  
- **FFprobe** - Media analysis tool (included with FFmpeg)

### Install FFmpeg

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Windows
# Download from https://ffmpeg.org/download.html
```

### Build from Source

```bash
git clone https://github.com/rishad1234/term-video-transcoder.git
cd term-video-transcoder
go mod tidy
go build -o transcoder .
```

## Usage Examples

### Media Analysis

```bash
# Basic file information
./transcoder info video.mp4

# Analyze different formats
./transcoder info movie.mkv
./transcoder info audio.mp3
```

### Video Conversion

```bash
# Basic conversion
./transcoder convert input.avi output.mp4

# Quality presets
./transcoder convert input.mov output.mp4 --preset high

# Custom resolution and bitrate
./transcoder convert input.mkv output.mp4 \
  --resolution 1280x720 --video-bitrate 2M
```

### Audio Extraction

```bash
# Basic audio extraction
./transcoder extract video.mp4 audio.mp3

# High quality FLAC
./transcoder extract movie.mkv soundtrack.flac --quality high

# Custom bitrate and format
./transcoder extract input.avi output.mp3 --bitrate 320k
```

## Testing

The project includes comprehensive test scripts:

```bash
# Quick test (basic functionality)
./quick_test.sh

# Audio extraction testing
./audio_test.sh

# Comprehensive testing
./full_test.sh
```

## Contributing

Contributions welcome! Areas for improvement:

- Additional format support
- Batch processing features
- Performance optimizations
- Platform-specific builds

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
