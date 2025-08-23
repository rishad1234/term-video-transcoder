package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Version information
	version = "0.1.0"
	
	// Global flags
	verbose bool
	quiet   bool
	output  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "transcoder",
	Short: "A terminal video transcoder powered by ffmpeg",
	Long: color.CyanString(`
┌─────────────────────────────────────────────┐
│          Terminal Video Transcoder          │
│         Powered by FFmpeg & FFprobe         │
└─────────────────────────────────────────────┘

A fast, efficient command-line tool for video transcoding,
media analysis, and batch processing.

Examples:
  transcoder convert input.avi output.mp4
  transcoder info video.mkv
  transcoder batch *.mov --format mp4
  transcoder convert input.mp4 output.webm --preset high
	`),
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags - verbose is now default
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "verbose output (enabled by default)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode (minimal output)")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output file or directory")
	
	// Add version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("Terminal Video Transcoder %s\n", version))
}
