package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
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
	
	// Phase 2: Custom Parameters
	videoCodec    string
	audioCodec    string
	videoBitrate  string
	audioBitrate  string
	resolution    string
	framerate     string
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
  # Basic conversion with presets
  transcoder convert input.avi output.mp4
  transcoder convert movie.mkv movie.webm --preset high
  
  # Custom codec selection
  transcoder convert input.mp4 output.webm --video-codec libvpx-vp9 --audio-codec libopus
  
  # Bitrate control
  transcoder convert input.mov output.mp4 --video-bitrate 2M --audio-bitrate 192k
  
  # Resolution and frame rate
  transcoder convert input.mkv output.mp4 --resolution 1920x1080 --framerate 30
  
  # Combined custom parameters
  transcoder convert input.avi output.mp4 --video-codec libx264 --video-bitrate 4M --resolution 1280x720`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConvert(cmd, args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Basic flags
	convertCmd.Flags().StringVarP(&preset, "preset", "p", "medium", "quality preset (low, medium, high)")
	convertCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite output file if it exists")
	
	// Phase 2: Custom Parameters
	convertCmd.Flags().StringVar(&videoCodec, "video-codec", "", "video codec (libx264, libx265, libvpx-vp9, etc.)")
	convertCmd.Flags().StringVar(&audioCodec, "audio-codec", "", "audio codec (aac, libopus, libmp3lame, etc.)")
	convertCmd.Flags().StringVar(&videoBitrate, "video-bitrate", "", "video bitrate (e.g., 2M, 1500k)")
	convertCmd.Flags().StringVar(&audioBitrate, "audio-bitrate", "", "audio bitrate (e.g., 192k, 128k)")
	convertCmd.Flags().StringVar(&resolution, "resolution", "", "output resolution (e.g., 1920x1080, 1280x720)")
	convertCmd.Flags().StringVar(&framerate, "framerate", "", "output frame rate (e.g., 30, 24, 60)")
	
	// Track when preset flag is explicitly set
	convertCmd.Flags().Lookup("preset").Changed = false
}

func runConvert(cmd *cobra.Command, inputPath, outputPath string) error {
	// Validate preset
	if !isValidPreset(preset) {
		return fmt.Errorf("invalid preset '%s'. Valid options: low, medium, high", preset)
	}
	
	// Validate custom parameters
	if err := validateCustomParameters(); err != nil {
		return err
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
	
	// Check if any custom parameters were set
	customParamsSet := hasCustomParameters()
	
	// Determine verbosity: quiet overrides verbose
	useVerbose := verbose && !quiet
	
	// Create custom parameters struct
	customParams := transcoder.CustomParameters{
		VideoCodec:   videoCodec,
		AudioCodec:   audioCodec,
		VideoBitrate: videoBitrate,
		AudioBitrate: audioBitrate,
		Resolution:   resolution,
		Framerate:    framerate,
	}
	
	// Perform the conversion
	err := transcoder.ConvertVideoWithCustomParams(inputPath, outputPath, preset, presetExplicit, customParamsSet, customParams, useVerbose)
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

// validateCustomParameters validates the custom parameter values
func validateCustomParameters() error {
	// Validate resolution format
	if resolution != "" {
		if !isValidResolution(resolution) {
			return fmt.Errorf("invalid resolution format '%s'. Use format like 1920x1080", resolution)
		}
	}
	
	// Validate framerate
	if framerate != "" {
		if !isValidFramerate(framerate) {
			return fmt.Errorf("invalid framerate '%s'. Use positive numbers like 30, 24, 60", framerate)
		}
	}
	
	// Validate bitrate formats
	if videoBitrate != "" {
		if !isValidBitrate(videoBitrate) {
			return fmt.Errorf("invalid video bitrate format '%s'. Use format like 2M, 1500k", videoBitrate)
		}
	}
	
	if audioBitrate != "" {
		if !isValidBitrate(audioBitrate) {
			return fmt.Errorf("invalid audio bitrate format '%s'. Use format like 192k, 128k", audioBitrate)
		}
	}
	
	return nil
}

// hasCustomParameters checks if any custom parameters were set
func hasCustomParameters() bool {
	return videoCodec != "" || audioCodec != "" || videoBitrate != "" || 
		   audioBitrate != "" || resolution != "" || framerate != ""
}

// isValidResolution checks if resolution is in format WIDTHxHEIGHT
func isValidResolution(res string) bool {
	parts := strings.Split(res, "x")
	if len(parts) != 2 {
		return false
	}
	
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return false
		}
	}
	return true
}

// isValidFramerate checks if framerate is a valid positive number
func isValidFramerate(fps string) bool {
	rate, err := strconv.ParseFloat(fps, 64)
	return err == nil && rate > 0
}

// isValidBitrate checks if bitrate is in valid format (e.g., 2M, 1500k)
func isValidBitrate(bitrate string) bool {
	if len(bitrate) < 2 {
		return false
	}
	
	// Check if it ends with k, K, m, or M
	suffix := strings.ToLower(bitrate[len(bitrate)-1:])
	if suffix != "k" && suffix != "m" {
		return false
	}
	
	// Check if the numeric part is valid
	numeric := bitrate[:len(bitrate)-1]
	_, err := strconv.ParseFloat(numeric, 64)
	return err == nil
}
