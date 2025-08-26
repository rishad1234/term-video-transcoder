package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// manualCmd represents the manual command
var manualCmd = &cobra.Command{
	Use:   "manual",
	Short: "Display the comprehensive manual for all commands and options",
	Long: `Display the comprehensive manual for Terminal Video Transcoder.

This manual contains detailed information about all commands, options, examples,
supported formats, quality presets, and troubleshooting tips.

For the complete manual, see: MANUAL.md`,
	Run: func(cmd *cobra.Command, args []string) {
		showManual()
	},
}

func init() {
	rootCmd.AddCommand(manualCmd)
}

func showManual() {
	manual := `
┌─────────────────────────────────────────────┐
│      Terminal Video Transcoder Manual       │
│         Comprehensive Reference             │
└─────────────────────────────────────────────┘

QUICK REFERENCE
===============

COMMANDS:
  info     Analyze media files (duration, codecs, metadata)
  convert  Convert between video formats with custom options
  extract  Extract audio from videos to various formats
  manual   Show this manual

GLOBAL OPTIONS:
  -h, --help      Show help
  -v, --verbose   Verbose output (default)
  -q, --quiet     Quiet mode
  -o, --output    Output file/directory
  --version       Show version

CONVERT COMMAND
===============

Usage: transcoder convert [input] [output] [options]

QUALITY PRESETS:
  --preset low     Fast encoding, smaller files (mobile)
  --preset medium  Balanced quality/size (default)
  --preset high    Best quality, larger files (archival)

CUSTOM VIDEO OPTIONS:
  --video-codec      Video codec (libx264, libx265, libvpx-vp9)
  --video-bitrate    Bitrate (2M, 1500k, 4M)
  --resolution       Resolution (1920x1080, 1280x720, 640x360)
  --framerate        Frame rate (30, 24, 60)

CUSTOM AUDIO OPTIONS:
  --audio-codec      Audio codec (aac, libopus, libmp3lame)
  --audio-bitrate    Bitrate (192k, 128k, 256k)

OTHER OPTIONS:
  -f, --force        Overwrite existing files

SUPPORTED VIDEO FORMATS:
  MP4     Most compatible, web-friendly
  AVI     Legacy format, widely supported
  MKV     High quality, multiple tracks
  WebM    Modern web format, efficient
  MOV     Apple format, high quality

EXTRACT COMMAND
===============

Usage: transcoder extract [input] [output] [options]

QUALITY PRESETS:
  --quality low      128k bitrate (streaming)
  --quality medium   192k bitrate (default)
  --quality high     320k bitrate (high quality)

CUSTOM OPTIONS:
  -b, --bitrate      Audio bitrate (320k, 192k, 128k)
  -c, --codec        Audio codec (libmp3lame, flac, aac)
  -s, --sample-rate  Sample rate (44100, 48000)
  --channels         Channels (1=mono, 2=stereo)
  -f, --force        Overwrite existing files

SUPPORTED AUDIO FORMATS:
  MP3     Universal compatibility
  WAV     Uncompressed, best quality
  AAC     Modern, efficient compression
  FLAC    Lossless compression
  OGG     Open format, good compression
  M4A     Apple format, AAC container

COMMON EXAMPLES
===============

# Basic operations
transcoder info video.mp4
transcoder convert input.avi output.mp4
transcoder extract movie.mkv audio.mp3

# Quality control
transcoder convert input.mov output.mp4 --preset high
transcoder extract video.mp4 audio.flac --quality high

# Custom parameters
transcoder convert input.avi output.webm \
  --video-codec libvpx-vp9 --video-bitrate 2M \
  --audio-codec libopus --audio-bitrate 128k

# Resolution scaling
transcoder convert input.mkv mobile.mp4 \
  --resolution 640x360 --preset low

# Audio extraction with custom settings
transcoder extract video.mp4 podcast.mp3 \
  --bitrate 128k --channels 1

BITRATE GUIDELINES
==================

VIDEO BITRATES:
  480p:  500k-1M
  720p:  1M-3M
  1080p: 3M-6M
  4K:    15M-25M

AUDIO BITRATES:
  Voice:        64k-128k
  Music (std):  128k-192k
  Music (HQ):   256k-320k
  Lossless:     Variable (FLAC)

TROUBLESHOOTING
===============

Common Issues:
• "ffmpeg not found" → Install: brew install ffmpeg (macOS)
• "Permission denied" → Run: chmod +x transcoder
• Poor quality → Increase bitrate or use --preset high
• Slow conversion → Use --preset low or libx264 codec

Performance Tips:
• Use appropriate presets for your use case
• Match input resolution when possible
• Choose libx264 for speed, libx265 for quality
• Process files sequentially, not simultaneously

TESTING & VALIDATION
====================

Test Scripts:
  ./quick_test.sh          Fast testing (video + audio)
  ./audio_test.sh          Comprehensive audio testing  
  ./test_all_features.sh   Complete video testing
  ./security_test.sh       Security vulnerability testing

Help Commands:
  transcoder --help
  transcoder convert --help
  transcoder extract --help
  transcoder info --help

For the complete manual with detailed examples and format tables,
see the MANUAL.md file in the project directory.

Project: https://github.com/rishad1234/term-video-transcoder
`
	fmt.Print(manual)
}
