# Terminal Video Transcoder - Manual

## Table of Contents

- [Overview](#overview)
- [Installation & Requirements](#installation--requirements)
- [Quick Start](#quick-start)
- [Commands](#commands)
  - [info](#info---media-analysis)
  - [convert](#convert---video-conversion)
  - [extract](#extract---audio-extraction)
  - [completion](#completion---shell-autocompletion)
- [Global Options](#global-options)
- [Examples](#examples)
- [Supported Formats](#supported-formats)
- [Quality Presets](#quality-presets)
- [Custom Parameters](#custom-parameters)
- [Troubleshooting](#troubleshooting)

## Overview

Terminal Video Transcoder is a fast, efficient command-line tool for video transcoding, media analysis, and audio extraction. Built on top of FFmpeg and FFprobe, it provides an intuitive interface for common media processing tasks while supporting advanced customization options.

**Key Features:**

- Media information analysis and display
- Video format conversion with automatic codec selection
- Audio extraction from video files
- Quality presets for different use cases
- Custom parameter control for advanced users
- Batch processing capabilities
- Real-time progress feedback
- Cross-platform support

## Installation & Requirements

### System Requirements

- **FFmpeg** (latest stable version) - Video/audio processing engine
- **FFprobe** (included with FFmpeg) - Media analysis tool

### Installation

```bash
# Install FFmpeg
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Windows
# Download from https://ffmpeg.org/download.html

# Build transcoder
git clone https://github.com/rishad1234/term-video-transcoder.git
cd term-video-transcoder
go build -o transcoder .
```

## Quick Start

```bash
# Analyze a video file
./transcoder info video.mp4

# Convert video with default settings
./transcoder convert input.avi output.mp4

# Extract audio from video
./transcoder extract movie.mkv soundtrack.mp3

# Convert with high quality preset
./transcoder convert input.mov output.webm --preset high
```

## Commands

### `info` - Media Analysis

Analyze and display comprehensive information about media files.

#### Usage

```bash
transcoder info [file] [flags]
```

#### What it shows

- Format and container information
- Video streams (codec, resolution, frame rate, bitrate)
- Audio streams (codec, sample rate, channels, bitrate)
- Duration and file size
- Metadata

#### Examples

```bash
# Basic media info
transcoder info video.mp4

# Analyze different formats
transcoder info movie.mkv
transcoder info audio.mp3
transcoder info presentation.webm
```

#### Flags

- `-h, --help` - Help for info command

---

### `convert` - Video Conversion

Convert video files between different formats with automatic codec selection and optimization.

#### Usage

```bash
transcoder convert [input] [output] [flags]
```

#### Supported Input/Output Formats

- **MP4** - Most compatible, good compression
- **AVI** - Legacy format, widely supported
- **MKV** - Open format, excellent for high quality
- **WebM** - Web-optimized, modern codecs
- **MOV** - Apple format, high quality

#### Quality Presets

- `--preset low` - Fast encoding, smaller files (mobile-friendly)
- `--preset medium` - Balanced quality and size (default)
- `--preset high` - Best quality, larger files (archival)

#### Custom Parameters

- `--video-codec` - Video codec (libx264, libx265, libvpx-vp9, etc.)
- `--audio-codec` - Audio codec (aac, libopus, libmp3lame, etc.)
- `--video-bitrate` - Video bitrate (e.g., 2M, 1500k)
- `--audio-bitrate` - Audio bitrate (e.g., 192k, 128k)
- `--resolution` - Output resolution (e.g., 1920x1080, 1280x720)
- `--framerate` - Output frame rate (e.g., 30, 24, 60)

#### Other Options

- `-f, --force` - Overwrite output file if it exists
- `-p, --preset` - Quality preset (low, medium, high)

#### Examples

```bash
# Basic conversion
transcoder convert input.avi output.mp4

# High quality conversion
transcoder convert movie.mkv movie.webm --preset high

# Custom codec selection
transcoder convert input.mp4 output.webm \
  --video-codec libvpx-vp9 --audio-codec libopus

# Bitrate control
transcoder convert input.mov output.mp4 \
  --video-bitrate 2M --audio-bitrate 192k

# Resolution scaling
transcoder convert input.mkv output.mp4 --resolution 1920x1080

# Frame rate adjustment
transcoder convert input.avi output.mp4 --framerate 30

# Combined parameters
transcoder convert input.avi output.mp4 \
  --video-codec libx264 --video-bitrate 4M --resolution 1280x720
```

---

### `extract` - Audio Extraction

Extract audio tracks from video files and convert to various audio formats.

#### Usage

```bash
transcoder extract [input] [output] [flags]
```

#### Supported Audio Formats

- **MP3** - Universal compatibility, good compression
- **WAV** - Uncompressed, highest quality
- **AAC** - Modern codec, efficient compression
- **FLAC** - Lossless compression, archival quality
- **OGG** - Open format, good compression
- **M4A** - Apple format, AAC container

#### Quality Presets

- `--quality low` - 128k bitrate (streaming quality)
- `--quality medium` - 192k bitrate (standard quality, default)
- `--quality high` - 320k bitrate (high quality)

#### Custom Parameters

- `-b, --bitrate` - Audio bitrate (e.g., 320k, 192k, 128k)
- `-c, --codec` - Audio codec (libmp3lame, aac, flac, libvorbis, etc.)
- `-s, --sample-rate` - Sample rate (e.g., 44100, 48000)
- `--channels` - Number of channels (1=mono, 2=stereo, 6=5.1)

#### Other Options

- `-f, --force` - Overwrite output file if it exists
- `--quality` - Audio quality preset (low, medium, high)

#### Examples

```bash
# Basic audio extraction
transcoder extract video.mp4 audio.mp3

# High quality FLAC
transcoder extract movie.mkv soundtrack.flac --quality high

# Custom bitrate
transcoder extract input.avi output.mp3 --bitrate 320k

# Specific codec
transcoder extract video.webm audio.ogg --codec libvorbis

# Custom sample rate
transcoder extract input.mp4 output.wav --sample-rate 48000

# Mono conversion
transcoder extract video.mkv audio.mp3 --channels 1

# Low quality for streaming
transcoder extract input.avi output.mp3 --quality low
```

---

### `completion` - Shell Autocompletion

Generate autocompletion scripts for your shell.

#### Usage

```bash
transcoder completion [bash|zsh|fish|powershell]
```

#### Examples

```bash
# Bash
transcoder completion bash > /etc/bash_completion.d/transcoder

# Zsh
transcoder completion zsh > "${fpath[1]}/_transcoder"

# Fish
transcoder completion fish > ~/.config/fish/completions/transcoder.fish
```

## Global Options

These options work with all commands:

- `-h, --help` - Show help information
- `-o, --output string` - Output file or directory
- `-q, --quiet` - Quiet mode (minimal output)
- `-v, --verbose` - Verbose output (enabled by default)
- `--version` - Show version information

## Examples

### Common Workflows

#### Media Analysis

```bash
# Quick file info
transcoder info video.mp4

# Analyze multiple files
transcoder info *.mp4
```

#### Format Conversion

```bash
# Standard conversion
transcoder convert input.avi output.mp4

# Batch conversion (coming soon)
transcoder batch *.avi --format mp4 --output-dir converted/
```

#### Audio Processing

```bash
# Extract soundtrack
transcoder extract movie.mkv soundtrack.flac --quality high

# Create podcast version
transcoder extract interview.mp4 podcast.mp3 --bitrate 128k --channels 1
```

#### Quality Optimization

```bash
# Mobile-friendly version
transcoder convert input.mkv mobile.mp4 \
  --preset low --resolution 640x360

# Archive quality
transcoder convert input.avi archive.mkv \
  --preset high --video-codec libx265
```

#### Web Optimization

```bash
# Modern web format
transcoder convert input.mp4 web.webm \
  --video-codec libvpx-vp9 --audio-codec libopus \
  --video-bitrate 1M --audio-bitrate 128k
```

## Supported Formats

### Video Formats

| Format | Extension | Best Use Case |
|--------|-----------|---------------|
| MP4    | .mp4      | General purpose, web, mobile |
| AVI    | .avi      | Legacy compatibility |
| MKV    | .mkv      | High quality, multiple tracks |
| WebM   | .webm     | Web streaming, modern codecs |
| MOV    | .mov      | Apple ecosystem, high quality |

### Audio Formats

| Format | Extension | Best Use Case |
|--------|-----------|---------------|
| MP3    | .mp3      | Universal compatibility |
| WAV    | .wav      | Uncompressed, editing |
| AAC    | .aac      | Modern devices, streaming |
| FLAC   | .flac     | Lossless, archival |
| OGG    | .ogg      | Open source, good compression |
| M4A    | .m4a      | Apple devices, AAC container |

### Video Codecs

| Codec     | Quality | Speed | Compatibility |
|-----------|---------|-------|---------------|
| libx264   | Good    | Fast  | Excellent     |
| libx265   | Better  | Slow  | Good          |
| libvpx-vp9| Best    | Slow  | Modern        |

### Audio Codecs

| Codec      | Quality | Size  | Compatibility |
|------------|---------|-------|---------------|
| aac        | Good    | Small | Excellent     |
| libmp3lame | Good    | Small | Universal     |
| libopus    | Better  | Small | Modern        |
| flac       | Perfect | Large | Good          |

## Quality Presets

### Low Quality

- **Use case**: Mobile devices, streaming, quick sharing
- **Video**: Lower bitrate, faster encoding
- **Audio**: 128k bitrate
- **File size**: Smallest

### Medium Quality (Default)

- **Use case**: General purpose, balanced quality/size
- **Video**: Balanced bitrate and quality
- **Audio**: 192k bitrate
- **File size**: Moderate

### High Quality

- **Use case**: Archival, professional work, large screens
- **Video**: Higher bitrate, best quality settings
- **Audio**: 320k bitrate
- **File size**: Largest

## Custom Parameters

### Bitrate Guidelines

#### Video Bitrates

- **480p**: 500k-1M
- **720p**: 1M-3M
- **1080p**: 3M-6M
- **4K**: 15M-25M

#### Audio Bitrates

- **Voice/Podcast**: 64k-128k
- **Music (standard)**: 128k-192k
- **Music (high quality)**: 256k-320k
- **Lossless**: Variable (FLAC)

### Resolution Guidelines

- **Mobile**: 640x360, 854x480
- **Web**: 1280x720, 1920x1080
- **HD**: 1920x1080
- **4K**: 3840x2160

### Frame Rate Guidelines

- **Film**: 24 fps
- **TV/Video**: 30 fps
- **Gaming/Sports**: 60 fps

## Troubleshooting

### Common Issues

#### "Command not found"

- Ensure the transcoder binary is in your PATH or use `./transcoder`
- Check that the binary was built successfully with `go build`

#### "ffmpeg not found"

- Install FFmpeg: `brew install ffmpeg` (macOS) or `apt install ffmpeg` (Ubuntu)
- Ensure FFmpeg is in your PATH

#### "Permission denied"

- Make the binary executable: `chmod +x transcoder`
- Check write permissions for output directory

#### Conversion Fails

- Check input file is valid: `transcoder info input.mp4`
- Verify enough disk space for output
- Try different codec combinations

#### Poor Quality Output

- Increase bitrate: `--video-bitrate 4M --audio-bitrate 256k`
- Use higher quality preset: `--preset high`
- Choose better codecs: `--video-codec libx265`

#### Slow Conversion

- Use faster presets: `--preset low`
- Choose faster codecs: `--video-codec libx264`
- Reduce resolution: `--resolution 1280x720`

### Getting Help

1. **Built-in help**: `transcoder --help` or `transcoder [command] --help`
2. **Test your setup**: Run `./quick_test.sh` or `./audio_test.sh`
3. **Check FFmpeg**: `ffmpeg -version` and `ffprobe -version`
4. **File an issue**: Report bugs on the project repository

### Performance Tips

1. **Use appropriate presets** for your use case
2. **Match input resolution** when possible (avoid unnecessary scaling)
3. **Choose efficient codecs** (libx264 over libx265 for speed)
4. **Monitor system resources** during conversion
5. **Process multiple files sequentially** rather than simultaneously

---

**Version**: 1.0.0  
**Last Updated**: August 2025  
**Project**: [Terminal Video Transcoder](https://github.com/rishad1234/term-video-transcoder)
