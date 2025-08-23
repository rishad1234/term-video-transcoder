package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

	// Determine output destination
	var writer io.Writer = os.Stdout
	var outputFile *os.File
	
	if output != "" {
		outputFile, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outputFile.Close()
		writer = outputFile
	}

	// Display the information with verbose mode consideration
	displayMediaInfo(info, verbose, writer)
	
	if output != "" {
		fmt.Printf("Media information saved to: %s\n", output)
	}
	
	return nil
}

func displayMediaInfo(info *analyzer.MediaInfo, verbose bool, writer io.Writer) {
	// Check if we're writing to a file (disable colors)
	isFile := writer != os.Stdout
	
	// Header
	if isFile {
		fmt.Fprintln(writer, "===============================================")
		if verbose {
			fmt.Fprintln(writer, "          Detailed Media Information")
		} else {
			fmt.Fprintln(writer, "              Media Information")
		}
		fmt.Fprintln(writer, "===============================================")
	} else {
		color.Cyan("┌─────────────────────────────────────────────┐")
		if verbose {
			color.Cyan("│          Detailed Media Information         │")
		} else {
			color.Cyan("│              Media Information              │")
		}
		color.Cyan("└─────────────────────────────────────────────┘")
	}
	fmt.Fprintln(writer)

	// File information
	if isFile {
		fmt.Fprintln(writer, "File Information:")
	} else {
		color.Yellow("📁 File Information:")
	}
	fmt.Fprintf(writer, "   Name: %s\n", filepath.Base(info.Filename))
	if verbose {
		fmt.Fprintf(writer, "   Full Path: %s\n", info.Filename)
	} else {
		fmt.Fprintf(writer, "   Path: %s\n", info.Filename)
	}
	fmt.Fprintf(writer, "   Format: %s\n", strings.ToUpper(info.Format))
	fmt.Fprintf(writer, "   Duration: %v\n", formatDuration(info.Duration))
	fmt.Fprintf(writer, "   Size: %s\n", formatBytes(info.Size))
	if info.Bitrate > 0 {
		fmt.Fprintf(writer, "   Overall Bitrate: %s\n", formatBitrate(info.Bitrate))
	}
	
	if verbose {
		fmt.Fprintf(writer, "   Duration (seconds): %.3f\n", info.Duration.Seconds())
		fmt.Fprintf(writer, "   Size (bytes): %d\n", info.Size)
		if info.Bitrate > 0 {
			fmt.Fprintf(writer, "   Bitrate (bps): %d\n", info.Bitrate)
		}
	}
	fmt.Fprintln(writer)

	// Video streams
	if len(info.VideoStreams) > 0 {
		if isFile {
			fmt.Fprintln(writer, "Video Streams:")
		} else {
			color.Green("🎥 Video Streams:")
		}
		for i, stream := range info.VideoStreams {
			fmt.Fprintf(writer, "   Stream %d:\n", i+1)
			if verbose {
				fmt.Fprintf(writer, "     Stream Index: %d\n", stream.Index)
			}
			fmt.Fprintf(writer, "     Codec: %s\n", stream.Codec)
			fmt.Fprintf(writer, "     Resolution: %dx%d\n", stream.Width, stream.Height)
			fmt.Fprintf(writer, "     Frame Rate: %s\n", stream.FrameRate)
			fmt.Fprintf(writer, "     Pixel Format: %s\n", stream.PixelFormat)
			if stream.Bitrate > 0 {
				fmt.Fprintf(writer, "     Bitrate: %s\n", formatBitrate(stream.Bitrate))
				if verbose {
					fmt.Fprintf(writer, "     Bitrate (bps): %d\n", stream.Bitrate)
				}
			}
			if verbose {
				fmt.Fprintf(writer, "     Aspect Ratio: %.2f:1\n", float64(stream.Width)/float64(stream.Height))
				totalPixels := stream.Width * stream.Height
				fmt.Fprintf(writer, "     Total Pixels: %d\n", totalPixels)
			}
			fmt.Fprintln(writer)
		}
	}

	// Audio streams
	if len(info.AudioStreams) > 0 {
		if isFile {
			fmt.Fprintln(writer, "Audio Streams:")
		} else {
			color.Magenta("🔊 Audio Streams:")
		}
		for i, stream := range info.AudioStreams {
			fmt.Fprintf(writer, "   Stream %d:\n", i+1)
			if verbose {
				fmt.Fprintf(writer, "     Stream Index: %d\n", stream.Index)
			}
			fmt.Fprintf(writer, "     Codec: %s\n", stream.Codec)
			fmt.Fprintf(writer, "     Sample Rate: %d Hz\n", stream.SampleRate)
			fmt.Fprintf(writer, "     Channels: %d\n", stream.Channels)
			if stream.Bitrate > 0 {
				fmt.Fprintf(writer, "     Bitrate: %s\n", formatBitrate(stream.Bitrate))
				if verbose {
					fmt.Fprintf(writer, "     Bitrate (bps): %d\n", stream.Bitrate)
				}
			}
			if stream.Language != "" && stream.Language != "und" {
				fmt.Fprintf(writer, "     Language: %s\n", stream.Language)
			} else if verbose {
				fmt.Fprintf(writer, "     Language: %s (undefined)\n", stream.Language)
			}
			if verbose {
				channelLayout := getChannelLayout(stream.Channels)
				fmt.Fprintf(writer, "     Channel Layout: %s\n", channelLayout)
			}
			fmt.Fprintln(writer)
		}
	}
	
	if verbose {
		if isFile {
			fmt.Fprintln(writer, "Technical Summary:")
		} else {
			color.Blue("🔧 Technical Summary:")
		}
		fmt.Fprintf(writer, "   Total Streams: %d\n", len(info.VideoStreams)+len(info.AudioStreams))
		fmt.Fprintf(writer, "   Video Streams: %d\n", len(info.VideoStreams))
		fmt.Fprintf(writer, "   Audio Streams: %d\n", len(info.AudioStreams))
		if len(info.VideoStreams) > 0 && info.Duration > 0 {
			fps := parseFrameRate(info.VideoStreams[0].FrameRate)
			totalFrames := int(info.Duration.Seconds() * fps)
			fmt.Fprintf(writer, "   Estimated Total Frames: %d\n", totalFrames)
		}
		fmt.Fprintln(writer)
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

// getChannelLayout returns a descriptive channel layout based on channel count
func getChannelLayout(channels int) string {
	switch channels {
	case 1:
		return "Mono"
	case 2:
		return "Stereo"
	case 3:
		return "2.1"
	case 4:
		return "4.0 (Quad)"
	case 5:
		return "5.0"
	case 6:
		return "5.1 Surround"
	case 7:
		return "6.1 Surround"
	case 8:
		return "7.1 Surround"
	default:
		return fmt.Sprintf("%d channels", channels)
	}
}

// parseFrameRate extracts FPS from frame rate string (e.g., "30/1" -> 30.0)
func parseFrameRate(frameRate string) float64 {
	parts := strings.Split(frameRate, "/")
	if len(parts) != 2 {
		return 0
	}
	
	numerator, err1 := strconv.ParseFloat(parts[0], 64)
	denominator, err2 := strconv.ParseFloat(parts[1], 64)
	
	if err1 != nil || err2 != nil || denominator == 0 {
		return 0
	}
	
	return numerator / denominator
}
