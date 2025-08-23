package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rishad1234/term-video-transcoder/internal/analyzer"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [file]",
	Short: "Display detailed information about a media file",
	Long: `Analyze and display comprehensive information about a media file including:
- Format and container information
- Video streams (codec, resolution, frame rate, bitrate)
- Audio streams (codec, sample rate, channels, bitrate)
- Duration and file size
- Metadata

Example:
  transcoder info video.mp4
  transcoder info movie.mkv`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInfo(args[0])
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(filepath string) error {
	// Check if ffprobe is available
	if err := analyzer.CheckFFProbe(); err != nil {
		return fmt.Errorf("ffprobe check failed: %w", err)
	}

	// Analyze the media file
	info, err := analyzer.AnalyzeMedia(filepath)
	if err != nil {
		return fmt.Errorf("failed to analyze media: %w", err)
	}

	// Display the information
	displayMediaInfo(info)
	return nil
}

func displayMediaInfo(info *analyzer.MediaInfo) {
	// Header
	color.Cyan("┌─────────────────────────────────────────────┐")
	color.Cyan("│              Media Information              │")
	color.Cyan("└─────────────────────────────────────────────┘")
	fmt.Println()

	// File information
	color.Yellow("📁 File Information:")
	fmt.Printf("   Name: %s\n", filepath.Base(info.Filename))
	fmt.Printf("   Path: %s\n", info.Filename)
	fmt.Printf("   Format: %s\n", strings.ToUpper(info.Format))
	fmt.Printf("   Duration: %v\n", formatDuration(info.Duration))
	fmt.Printf("   Size: %s\n", formatBytes(info.Size))
	if info.Bitrate > 0 {
		fmt.Printf("   Overall Bitrate: %s\n", formatBitrate(info.Bitrate))
	}
	fmt.Println()

	// Video streams
	if len(info.VideoStreams) > 0 {
		color.Green("🎥 Video Streams:")
		for i, stream := range info.VideoStreams {
			fmt.Printf("   Stream %d:\n", i+1)
			fmt.Printf("     Codec: %s\n", stream.Codec)
			fmt.Printf("     Resolution: %dx%d\n", stream.Width, stream.Height)
			fmt.Printf("     Frame Rate: %s\n", stream.FrameRate)
			fmt.Printf("     Pixel Format: %s\n", stream.PixelFormat)
			if stream.Bitrate > 0 {
				fmt.Printf("     Bitrate: %s\n", formatBitrate(stream.Bitrate))
			}
			fmt.Println()
		}
	}

	// Audio streams
	if len(info.AudioStreams) > 0 {
		color.Magenta("🔊 Audio Streams:")
		for i, stream := range info.AudioStreams {
			fmt.Printf("   Stream %d:\n", i+1)
			fmt.Printf("     Codec: %s\n", stream.Codec)
			fmt.Printf("     Sample Rate: %d Hz\n", stream.SampleRate)
			fmt.Printf("     Channels: %d\n", stream.Channels)
			if stream.Bitrate > 0 {
				fmt.Printf("     Bitrate: %s\n", formatBitrate(stream.Bitrate))
			}
			if stream.Language != "" {
				fmt.Printf("     Language: %s\n", stream.Language)
			}
			fmt.Println()
		}
	}
}

// Helper functions for formatting
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatBitrate(bitrate int64) string {
	const unit = 1000
	if bitrate < unit {
		return fmt.Sprintf("%d bps", bitrate)
	}
	div, exp := int64(unit), 0
	for n := bitrate / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cbps", float64(bitrate)/float64(div), "kMGTPE"[exp])
}
