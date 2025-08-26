package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rishad1234/term-video-transcoder/internal/security"
	"github.com/rishad1234/term-video-transcoder/internal/transcoder"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract [input] [output]",
	Short: "Extract audio from video files",
	Long: `Extract audio tracks from video files and convert to various audio formats.

Supported output formats: MP3, WAV, AAC, FLAC, OGG, M4A

The tool automatically detects the desired output format from the file extension
and applies appropriate codec selection and quality settings.

Examples:
  # Extract audio to MP3
  transcoder extract video.mp4 audio.mp3
  
  # Extract with high quality
  transcoder extract movie.mkv soundtrack.flac --quality high
  
  # Custom bitrate
  transcoder extract input.avi output.mp3 --bitrate 320k
  
  # Specific audio codec
  transcoder extract video.webm audio.ogg --codec libvorbis`,
	Args: cobra.ExactArgs(2),
	RunE: runExtract,
}

var (
	extractQuality    string
	extractBitrate    string
	extractCodec      string
	extractSampleRate string
	extractChannels   string
	extractForce      bool
)

func init() {
	rootCmd.AddCommand(extractCmd)

	// Audio quality preset (no shorthand to avoid conflict with global -q)
	extractCmd.Flags().StringVar(&extractQuality, "quality", "medium",
		"audio quality preset (low, medium, high)")

	// Custom audio parameters
	extractCmd.Flags().StringVarP(&extractBitrate, "bitrate", "b", "",
		"audio bitrate (e.g., 320k, 192k, 128k)")

	extractCmd.Flags().StringVarP(&extractCodec, "codec", "c", "",
		"audio codec (libmp3lame, aac, flac, libvorbis, etc.)")

	extractCmd.Flags().StringVarP(&extractSampleRate, "sample-rate", "s", "",
		"sample rate (e.g., 44100, 48000)")

	extractCmd.Flags().StringVar(&extractChannels, "channels", "",
		"number of channels (1=mono, 2=stereo, 6=5.1)")

	// Force overwrite flag
	extractCmd.Flags().BoolVarP(&extractForce, "force", "f", false,
		"overwrite output file if it exists")
}

func runExtract(cmd *cobra.Command, args []string) error {
	inputFile := args[0]
	outputFile := args[1]

	// Initialize security policy
	securityPolicy := security.NewDefaultSecurityPolicy()

	// Security validation for file paths
	if err := securityPolicy.ValidateFilePath(inputFile); err != nil {
		return fmt.Errorf("security validation failed for input path: %w", err)
	}

	if err := securityPolicy.ValidateFilePath(outputFile); err != nil {
		return fmt.Errorf("security validation failed for output path: %w", err)
	}

	if err := securityPolicy.ValidateFileFormat(outputFile); err != nil {
		return fmt.Errorf("security validation failed for output format: %w", err)
	}

	// Validate input file exists
	if !fileExists(inputFile) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Check if output file exists and handle overwrite
	if fileExists(outputFile) && !extractForce {
		return fmt.Errorf("output file already exists: %s (use --force to overwrite)", outputFile)
	}

	// Create audio extraction parameters
	params := transcoder.AudioExtractionParams{
		InputFile:  inputFile,
		OutputFile: outputFile,
		Quality:    extractQuality,
		Bitrate:    extractBitrate,
		Codec:      extractCodec,
		SampleRate: extractSampleRate,
		Channels:   extractChannels,
		Verbose:    verbose,
	}

	// Validate parameters
	if err := validateAudioParams(params); err != nil {
		return fmt.Errorf("invalid parameters: %v", err)
	}

	// Display extraction info
	if verbose {
		displayExtractionInfo(params)
	}

	// Perform audio extraction
	return transcoder.ExtractAudio(params)
}

func validateAudioParams(params transcoder.AudioExtractionParams) error {
	// Validate quality preset
	validQualities := []string{"low", "medium", "high"}
	if !contains(validQualities, params.Quality) {
		return fmt.Errorf("invalid quality preset: %s (valid: %s)",
			params.Quality, strings.Join(validQualities, ", "))
	}

	// Validate bitrate format if provided
	if params.Bitrate != "" {
		if !isValidAudioBitrate(params.Bitrate) {
			return fmt.Errorf("invalid bitrate format: %s (examples: 320k, 192k, 128k)", params.Bitrate)
		}
	}

	// Validate sample rate if provided
	if params.SampleRate != "" {
		validSampleRates := []string{"8000", "11025", "16000", "22050", "44100", "48000", "88200", "96000"}
		if !contains(validSampleRates, params.SampleRate) {
			return fmt.Errorf("invalid sample rate: %s (valid: %s)",
				params.SampleRate, strings.Join(validSampleRates, ", "))
		}
	}

	// Validate channels if provided
	if params.Channels != "" {
		validChannels := []string{"1", "2", "6", "8"}
		if !contains(validChannels, params.Channels) {
			return fmt.Errorf("invalid channel count: %s (valid: 1=mono, 2=stereo, 6=5.1, 8=7.1)", params.Channels)
		}
	}

	// Validate output format based on extension
	ext := strings.ToLower(filepath.Ext(params.OutputFile))
	supportedFormats := []string{".mp3", ".wav", ".aac", ".flac", ".ogg", ".m4a"}
	if !contains(supportedFormats, ext) {
		return fmt.Errorf("unsupported output format: %s (supported: %s)",
			ext, strings.Join(supportedFormats, ", "))
	}

	return nil
}

func displayExtractionInfo(params transcoder.AudioExtractionParams) {
	fmt.Println("🎵 Audio Extraction")
	fmt.Println("==================")
	fmt.Printf("📹 Input:   %s\n", params.InputFile)
	fmt.Printf("🎧 Output:  %s\n", params.OutputFile)
	fmt.Printf("🎯 Quality: %s\n", strings.ToUpper(params.Quality))

	if params.Bitrate != "" {
		fmt.Printf("📊 Bitrate: %s\n", params.Bitrate)
	}
	if params.Codec != "" {
		fmt.Printf("🔧 Codec:   %s\n", params.Codec)
	}
	if params.SampleRate != "" {
		fmt.Printf("📡 Sample:  %s Hz\n", params.SampleRate)
	}
	if params.Channels != "" {
		fmt.Printf("🔊 Channels: %s\n", params.Channels)
	}

	fmt.Println()
}

func isValidAudioBitrate(bitrate string) bool {
	// Check if bitrate ends with 'k' or 'K' and has valid number
	if len(bitrate) < 2 {
		return false
	}

	suffix := strings.ToLower(bitrate[len(bitrate)-1:])
	if suffix != "k" {
		return false
	}

	// Check if the prefix is a valid number
	numberPart := bitrate[:len(bitrate)-1]
	validBitrates := []string{"32", "64", "96", "128", "160", "192", "224", "256", "320"}
	return contains(validBitrates, numberPart)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
