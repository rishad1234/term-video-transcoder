package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/rishad1234/term-video-transcoder/internal/transcoder"
	"github.com/spf13/cobra"
)

var (
	// Convert command flags
	preset        string
	presetChanged bool // Track if preset flag was explicitly set
	force         bool
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [input] [output]",
	Short: "Convert video files between different formats",
	Long: `Convert video files between common formats with automatic codec selection.

Supported formats: MP4, AVI, MKV, WebM, MOV

The transcoder automatically selects the best codecs for the target format
and applies intelligent optimizations like stream copying when possible.

Examples:
  transcoder convert input.avi output.mp4
  transcoder convert movie.mkv movie.webm --preset high
  transcoder convert video.mov video.mp4 --force`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConvert(cmd, args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Add flags
	convertCmd.Flags().StringVarP(&preset, "preset", "p", "medium", "quality preset (low, medium, high)")
	convertCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite output file if it exists")
	
	// Track when preset flag is explicitly set
	convertCmd.Flags().Lookup("preset").Changed = false
}

func runConvert(cmd *cobra.Command, inputPath, outputPath string) error {
	// Validate preset
	if !isValidPreset(preset) {
		return fmt.Errorf("invalid preset '%s'. Valid options: low, medium, high", preset)
	}
	
	// Check if output file exists (unless force is specified)
	if !force {
		if err := checkOutputFile(outputPath); err != nil {
			return err
		}
	}
	
	// Display conversion info (unless quiet mode)
	if !quiet {
		displayConversionInfo(inputPath, outputPath, preset)
	}
	
	// Check if preset was explicitly set by user
	presetExplicit := cmd.Flags().Lookup("preset").Changed
	
	// Determine verbosity: quiet overrides verbose
	useVerbose := verbose && !quiet
	
	// Perform the conversion
	err := transcoder.ConvertVideo(inputPath, outputPath, preset, presetExplicit, useVerbose)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	
	// Success message (unless quiet mode)
	if !quiet {
		color.Green("✅ Conversion completed successfully!")
		fmt.Printf("Output saved to: %s\n", outputPath)
	}
	
	return nil
}

func isValidPreset(preset string) bool {
	validPresets := []string{"low", "medium", "high"}
	for _, valid := range validPresets {
		if preset == valid {
			return true
		}
	}
	return false
}

func checkOutputFile(outputPath string) error {
	if _, err := filepath.Abs(outputPath); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}
	return nil
}

func displayConversionInfo(inputPath, outputPath, preset string) {
	color.Cyan("🔄 Starting Video Conversion")
	fmt.Println()
	fmt.Printf("   Input:   %s\n", inputPath)
	fmt.Printf("   Output:  %s\n", outputPath)
	fmt.Printf("   Preset:  %s\n", strings.ToUpper(preset))
	fmt.Printf("   Format:  %s → %s\n", 
		strings.ToUpper(getFileExtension(inputPath)), 
		strings.ToUpper(getFileExtension(outputPath)))
	fmt.Println()
}

func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		return ext[1:] // Remove the dot
	}
	return ""
}
