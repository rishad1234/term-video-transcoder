package transcoder

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rishad1234/term-video-transcoder/internal/analyzer"
	"github.com/rishad1234/term-video-transcoder/internal/security"
)

// SupportedFormats defines the formats we can convert between
var SupportedFormats = map[string]bool{
	"mp4":  true,
	"avi":  true,
	"mkv":  true,
	"webm": true,
	"mov":  true,
}

// Global security policy for input validation
var securityPolicy = security.NewDefaultSecurityPolicy()

// CustomParameters holds user-specified custom encoding parameters
type CustomParameters struct {
	VideoCodec   string // User-specified video codec
	AudioCodec   string // User-specified audio codec
	VideoBitrate string // User-specified video bitrate (e.g., "2M", "1500k")
	AudioBitrate string // User-specified audio bitrate (e.g., "192k", "128k")
	Resolution   string // User-specified resolution (e.g., "1920x1080")
	Framerate    string // User-specified framerate (e.g., "30", "24")
}

// AudioExtractionParams holds parameters for audio extraction
type AudioExtractionParams struct {
	InputFile  string // Input video file path
	OutputFile string // Output audio file path
	Quality    string // Quality preset (low, medium, high)
	Bitrate    string // Custom bitrate (e.g., "320k", "192k")
	Codec      string // Custom codec (e.g., "libmp3lame", "aac")
	SampleRate string // Custom sample rate (e.g., "44100", "48000")
	Channels   string // Number of channels (e.g., "1", "2", "6")
	Verbose    bool   // Verbose output
}

// ConvertVideo converts a video file from one format to another (legacy function)
func ConvertVideo(inputPath, outputPath, preset string, presetExplicit, verbose bool) error {
	// Call the new function with empty custom parameters
	emptyParams := CustomParameters{}
	return ConvertVideoWithCustomParams(inputPath, outputPath, preset, presetExplicit, false, emptyParams, verbose)
}

// ConvertVideoWithCustomParams converts a video file with custom parameters support
func ConvertVideoWithCustomParams(inputPath, outputPath, preset string, presetExplicit, customParamsSet bool, customParams CustomParameters, verbose bool) error {
	// Step 1: Validate input file
	if err := validateInputFile(inputPath); err != nil {
		return err
	}

	// Step 2: Validate output format
	outputFormat := getFormatFromPath(outputPath)
	if !SupportedFormats[outputFormat] {
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	// Step 3: Security validation for file paths
	if err := securityPolicy.ValidateFilePath(inputPath); err != nil {
		return fmt.Errorf("security validation failed for input path: %w", err)
	}

	if err := securityPolicy.ValidateFilePath(outputPath); err != nil {
		return fmt.Errorf("security validation failed for output path: %w", err)
	}

	if err := securityPolicy.ValidateFileFormat(outputPath); err != nil {
		return fmt.Errorf("security validation failed for output format: %w", err)
	}

	// Step 4: Security validation for custom parameters
	if customParamsSet {
		if customParams.VideoCodec != "" {
			if err := securityPolicy.ValidateCodec(customParams.VideoCodec, "video"); err != nil {
				return fmt.Errorf("security validation failed for video codec: %w", err)
			}
		}

		if customParams.AudioCodec != "" {
			if err := securityPolicy.ValidateCodec(customParams.AudioCodec, "audio"); err != nil {
				return fmt.Errorf("security validation failed for audio codec: %w", err)
			}
		}

		if err := securityPolicy.ValidateBitrate(customParams.VideoBitrate); err != nil {
			return fmt.Errorf("security validation failed for video bitrate: %w", err)
		}

		if err := securityPolicy.ValidateBitrate(customParams.AudioBitrate); err != nil {
			return fmt.Errorf("security validation failed for audio bitrate: %w", err)
		}

		if err := securityPolicy.ValidateResolution(customParams.Resolution); err != nil {
			return fmt.Errorf("security validation failed for resolution: %w", err)
		}

		if err := securityPolicy.ValidateFramerate(customParams.Framerate); err != nil {
			return fmt.Errorf("security validation failed for framerate: %w", err)
		}
	}

	// Step 5: Analyze input media
	if verbose {
		color.Blue("🔍 Analyzing input media...")
	}
	inputInfo, err := analyzer.AnalyzeMedia(inputPath)
	if err != nil {
		return fmt.Errorf("failed to analyze input: %w", err)
	}

	// Step 6: Select optimal codecs (considering custom parameters and security)
	videoCodec, audioCodec, canCopy := selectCodecsWithCustomParamsSecure(inputInfo, outputFormat, preset, presetExplicit, customParamsSet, customParams, verbose)

	// Step 7: Apply preset-based bitrates if no custom bitrates specified
	secureCustomParams := customParams
	if !customParamsSet || customParams.VideoBitrate == "" {
		secureCustomParams.VideoBitrate = getPresetVideoBitrate(preset)
	}
	if !customParamsSet || customParams.AudioBitrate == "" {
		secureCustomParams.AudioBitrate = getPresetAudioBitrate(preset)
	}

	// Step 8: Build FFmpeg command (with security validation)
	cmd := buildFFmpegCommandWithCustomParams(inputPath, outputPath, videoCodec, audioCodec, preset, secureCustomParams, verbose)
	if cmd == nil {
		return fmt.Errorf("failed to build secure FFmpeg command")
	}

	// Step 9: Execute conversion
	if verbose {
		if canCopy {
			color.Green("⚡ Using stream copy (no re-encoding needed)")
		} else {
			color.Yellow("🔄 Re-encoding with selected codecs")
		}

		// Show custom parameters if any are set
		if customParamsSet {
			displayCustomParameters(secureCustomParams)
		}

		fmt.Printf("Command: %s\n\n", strings.Join(cmd.Args, " "))
	}

	return executeFFmpeg(cmd, inputInfo, verbose)
}

// selectCodecsWithCustomParamsSecure implements secure codec selection logic
func selectCodecsWithCustomParamsSecure(inputInfo *analyzer.MediaInfo, outputFormat, preset string, presetExplicit, customParamsSet bool, customParams CustomParameters, verbose bool) (string, string, bool) {
	// If custom codecs are specified, validate and use them
	if customParams.VideoCodec != "" && customParams.AudioCodec != "" {
		// Validation is already done in the calling function
		if verbose {
			color.Green("🎯 Using custom codecs specified by user")
			fmt.Printf("Video codec: %s\n", customParams.VideoCodec)
			fmt.Printf("Audio codec: %s\n", customParams.AudioCodec)
		}
		return customParams.VideoCodec, customParams.AudioCodec, false
	}

	// If any custom parameter is set, disable stream copy optimization
	if customParamsSet {
		videoCodec := customParams.VideoCodec
		audioCodec := customParams.AudioCodec

		// Use default codecs if not specified
		if videoCodec == "" {
			defaultVideo, _ := getDefaultCodecs(outputFormat)
			videoCodec = defaultVideo
		}
		if audioCodec == "" {
			_, defaultAudio := getDefaultCodecs(outputFormat)
			audioCodec = defaultAudio
		}

		// Apply quality presets (now returns only codec names)
		videoCodec = applyVideoPreset(videoCodec, preset)
		audioCodec = applyAudioPreset(audioCodec, preset)

		if verbose {
			color.Yellow("⚙️  Using custom parameters (stream copy disabled)")
			fmt.Printf("Video codec: %s\n", videoCodec)
			fmt.Printf("Audio codec: %s\n", audioCodec)
		}

		return videoCodec, audioCodec, false
	}

	// Fall back to original logic for automatic selection
	return selectCodecs(inputInfo, outputFormat, preset, presetExplicit, verbose)
}

// getPresetVideoBitrate returns video bitrate for quality presets
func getPresetVideoBitrate(preset string) string {
	switch preset {
	case "low":
		return "1M"
	case "medium":
		return "2M"
	case "high":
		return "4M"
	default:
		return "2M"
	}
}

// getPresetAudioBitrate returns audio bitrate for quality presets
func getPresetAudioBitrate(preset string) string {
	switch preset {
	case "low":
		return "128k"
	case "medium":
		return "192k"
	case "high":
		return "256k"
	default:
		return "192k"
	}
}

// validateInputFile checks if the input file exists and is readable
func validateInputFile(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}
	return nil
}

// getFormatFromPath extracts the file format from the file path
func getFormatFromPath(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		return strings.ToLower(ext[1:]) // Remove dot and convert to lowercase
	}
	return ""
}

// selectCodecsWithCustomParams implements codec selection logic with custom parameter support
func selectCodecsWithCustomParams(inputInfo *analyzer.MediaInfo, outputFormat, preset string, presetExplicit, customParamsSet bool, customParams CustomParameters, verbose bool) (string, string, bool) {
	// If custom codecs are specified, use them directly
	if customParams.VideoCodec != "" && customParams.AudioCodec != "" {
		if verbose {
			color.Green("🎯 Using custom codecs specified by user")
			fmt.Printf("Video codec: %s\n", customParams.VideoCodec)
			fmt.Printf("Audio codec: %s\n", customParams.AudioCodec)
		}
		return customParams.VideoCodec, customParams.AudioCodec, false
	}

	// If any custom parameter is set, disable stream copy optimization
	if customParamsSet {
		videoCodec := customParams.VideoCodec
		audioCodec := customParams.AudioCodec

		// Use default codecs if not specified
		if videoCodec == "" {
			defaultVideo, _ := getDefaultCodecs(outputFormat)
			videoCodec = defaultVideo
		}
		if audioCodec == "" {
			_, defaultAudio := getDefaultCodecs(outputFormat)
			audioCodec = defaultAudio
		}

		// Apply quality presets
		videoCodec = applyVideoPreset(videoCodec, preset)
		audioCodec = applyAudioPreset(audioCodec, preset)

		if verbose {
			color.Yellow("⚙️  Using custom parameters (stream copy disabled)")
			fmt.Printf("Video codec: %s\n", videoCodec)
			fmt.Printf("Audio codec: %s\n", audioCodec)
		}

		return videoCodec, audioCodec, false
	}

	// Fall back to original logic for automatic selection
	return selectCodecs(inputInfo, outputFormat, preset, presetExplicit, verbose)
}

// selectCodecs implements automatic codec selection logic
func selectCodecs(inputInfo *analyzer.MediaInfo, outputFormat, preset string, presetExplicit, verbose bool) (string, string, bool) {
	// Get default codecs for the output format
	defaultVideoCodec, defaultAudioCodec := getDefaultCodecs(outputFormat)

	// Check if we can use stream copy (no re-encoding)
	// Use stream copy only if:
	// 1. Formats are compatible, AND
	// 2. User did NOT explicitly set a preset (they want speed optimization)
	if canUseStreamCopy(inputInfo, outputFormat) && !presetExplicit {
		if verbose {
			color.Green("✨ Input codecs are compatible with output format")
		}
		return "copy", "copy", true
	}

	// Apply quality preset to codecs
	videoCodec := applyVideoPreset(defaultVideoCodec, preset)
	audioCodec := applyAudioPreset(defaultAudioCodec, preset)

	if verbose {
		fmt.Printf("Selected video codec: %s\n", videoCodec)
		fmt.Printf("Selected audio codec: %s\n", audioCodec)
	}

	return videoCodec, audioCodec, false
}

// displayCustomParameters shows the custom parameters being used
func displayCustomParameters(params CustomParameters) {
	color.Cyan("🔧 Custom Parameters:")
	if params.VideoCodec != "" {
		fmt.Printf("   Video Codec: %s\n", params.VideoCodec)
	}
	if params.AudioCodec != "" {
		fmt.Printf("   Audio Codec: %s\n", params.AudioCodec)
	}
	if params.VideoBitrate != "" {
		fmt.Printf("   Video Bitrate: %s\n", params.VideoBitrate)
	}
	if params.AudioBitrate != "" {
		fmt.Printf("   Audio Bitrate: %s\n", params.AudioBitrate)
	}
	if params.Resolution != "" {
		fmt.Printf("   Resolution: %s\n", params.Resolution)
	}
	if params.Framerate != "" {
		fmt.Printf("   Frame Rate: %s fps\n", params.Framerate)
	}
	fmt.Println()
}

// getDefaultCodecs returns the best default codecs for each format
func getDefaultCodecs(format string) (string, string) {
	switch format {
	case "mp4", "mov":
		return "libx264", "aac"
	case "webm":
		return "libvpx-vp9", "libopus"
	case "mkv":
		return "libx264", "aac"
	case "avi":
		return "libx264", "libmp3lame"
	default:
		return "libx264", "aac" // Safe defaults
	}
}

// canUseStreamCopy checks if we can copy streams without re-encoding
func canUseStreamCopy(inputInfo *analyzer.MediaInfo, outputFormat string) bool {
	if len(inputInfo.VideoStreams) == 0 || len(inputInfo.AudioStreams) == 0 {
		return false
	}

	videoCodec := inputInfo.VideoStreams[0].Codec
	audioCodec := inputInfo.AudioStreams[0].Codec

	// Check codec compatibility with output format
	switch outputFormat {
	case "mp4", "mov":
		return isCompatibleCodec(videoCodec, []string{"h264", "hevc"}) &&
			isCompatibleCodec(audioCodec, []string{"aac", "mp3"})
	case "webm":
		return isCompatibleCodec(videoCodec, []string{"vp8", "vp9", "av1"}) &&
			isCompatibleCodec(audioCodec, []string{"vorbis", "opus"})
	case "mkv":
		// MKV is very flexible, most codecs work
		return true
	case "avi":
		return isCompatibleCodec(videoCodec, []string{"h264", "xvid", "divx"}) &&
			isCompatibleCodec(audioCodec, []string{"mp3", "ac3"})
	}

	return false
}

// isCompatibleCodec checks if a codec is in the list of compatible codecs
func isCompatibleCodec(codec string, compatibleCodecs []string) bool {
	for _, compatible := range compatibleCodecs {
		if strings.Contains(strings.ToLower(codec), compatible) {
			return true
		}
	}
	return false
}

// applyVideoPreset applies quality settings to video codec
// Returns only the codec name for security - presets are handled via separate parameters
func applyVideoPreset(baseCodec, preset string) string {
	// For security, we only return the base codec name
	// Quality presets are now handled through separate validated parameters
	switch baseCodec {
	case "libx264", "libx265", "libvpx-vp9", "copy":
		return baseCodec
	default:
		// Default to safe codec if unknown
		return "libx264"
	}
}

// applyAudioPreset applies quality settings to audio codec
// Returns only the codec name for security - presets are handled via separate parameters
func applyAudioPreset(baseCodec, preset string) string {
	// For security, we only return the base codec name
	// Quality presets are now handled through separate validated parameters
	switch baseCodec {
	case "aac", "libopus", "libmp3lame", "libvorbis", "flac", "copy":
		return baseCodec
	default:
		// Default to safe codec if unknown
		return "aac"
	}
}

// buildFFmpegCommandWithCustomParams constructs the FFmpeg command with custom parameters
// This function now includes security validation to prevent command injection
func buildFFmpegCommandWithCustomParams(input, output, videoCodec, audioCodec, preset string, customParams CustomParameters, verbose bool) *exec.Cmd {
	// Security validation for all parameters
	if err := securityPolicy.ValidateFilePath(input); err != nil {
		if verbose {
			color.Red("Security validation failed for input path: %v", err)
		}
		return nil
	}

	if err := securityPolicy.ValidateFilePath(output); err != nil {
		if verbose {
			color.Red("Security validation failed for output path: %v", err)
		}
		return nil
	}

	args := []string{
		"ffmpeg",
		"-i", input,
	}

	// Add video codec parameters with security validation
	if videoCodec == "copy" {
		args = append(args, "-c:v", "copy")
	} else {
		// Validate video codec - prevent command injection
		if err := securityPolicy.ValidateCodec(videoCodec, "video"); err != nil {
			if verbose {
				color.Red("Security validation failed for video codec: %v", err)
			}
			return nil
		}

		// Only use the validated codec name - no additional parameters
		args = append(args, "-c:v", videoCodec)

		// Add custom video bitrate if specified and validated
		if customParams.VideoBitrate != "" {
			if err := securityPolicy.ValidateBitrate(customParams.VideoBitrate); err != nil {
				if verbose {
					color.Red("Security validation failed for video bitrate: %v", err)
				}
				return nil
			}
			args = append(args, "-b:v", customParams.VideoBitrate)
		}
	}

	// Add audio codec parameters with security validation
	if audioCodec == "copy" {
		args = append(args, "-c:a", "copy")
	} else {
		// Validate audio codec - prevent command injection
		if err := securityPolicy.ValidateCodec(audioCodec, "audio"); err != nil {
			if verbose {
				color.Red("Security validation failed for audio codec: %v", err)
			}
			return nil
		}

		// Only use the validated codec name - no additional parameters
		args = append(args, "-c:a", audioCodec)

		// Add custom audio bitrate if specified and validated
		if customParams.AudioBitrate != "" {
			if err := securityPolicy.ValidateBitrate(customParams.AudioBitrate); err != nil {
				if verbose {
					color.Red("Security validation failed for audio bitrate: %v", err)
				}
				return nil
			}
			args = append(args, "-b:a", customParams.AudioBitrate)
		}
	}

	// Add resolution scaling if specified and validated
	if customParams.Resolution != "" {
		if err := securityPolicy.ValidateResolution(customParams.Resolution); err != nil {
			if verbose {
				color.Red("Security validation failed for resolution: %v", err)
			}
			return nil
		}
		args = append(args, "-s", customParams.Resolution)
	}

	// Add framerate if specified and validated
	if customParams.Framerate != "" {
		if err := securityPolicy.ValidateFramerate(customParams.Framerate); err != nil {
			if verbose {
				color.Red("Security validation failed for framerate: %v", err)
			}
			return nil
		}
		args = append(args, "-r", customParams.Framerate)
	}

	// Add output file
	args = append(args, "-y", output) // -y to overwrite without asking

	return exec.Command(args[0], args[1:]...)
}

// buildFFmpegCommand constructs the FFmpeg command with all parameters (legacy function)
func buildFFmpegCommand(input, output, videoCodec, audioCodec, preset string, verbose bool) *exec.Cmd {
	emptyParams := CustomParameters{}
	return buildFFmpegCommandWithCustomParams(input, output, videoCodec, audioCodec, preset, emptyParams, verbose)
}

// executeFFmpeg runs the FFmpeg command and handles output
func executeFFmpeg(cmd *exec.Cmd, inputInfo *analyzer.MediaInfo, verbose bool) error {
	if verbose {
		color.Blue("🚀 Starting FFmpeg conversion...")
		// In verbose mode, show FFmpeg output directly
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Non-verbose mode: show progress bar
	return executeFFmpegWithProgress(cmd, inputInfo)
}

// executeFFmpegWithProgress runs FFmpeg and displays a progress indicator
func executeFFmpegWithProgress(cmd *exec.Cmd, inputInfo *analyzer.MediaInfo) error {
	color.Blue("🚀 Starting FFmpeg conversion...")

	// Show initial progress
	totalSeconds := inputInfo.Duration.Seconds()
	fmt.Printf("⏳ Processing %.1fs video...\n", totalSeconds)

	// Add progress reporting to stderr using -stats_period
	newArgs := make([]string, 0, len(cmd.Args)+2)
	newArgs = append(newArgs, cmd.Args[0])            // ffmpeg
	newArgs = append(newArgs, "-stats_period", "0.2") // Update stats every 0.2 seconds
	newArgs = append(newArgs, cmd.Args[1:]...)        // Rest of arguments
	cmd.Args = newArgs

	// Create pipes for stderr (stats)
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Suppress stdout in non-verbose mode
	cmd.Stdout = nil

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Track if we've shown any progress
	progressShown := false

	// Parse FFmpeg stats output for progress
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		timeRegex := regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2})\.(\d{2})`)
		speedRegex := regexp.MustCompile(`speed=\s*([0-9.]+)x`)

		for scanner.Scan() {
			line := scanner.Text()

			// Parse time progress
			if matches := timeRegex.FindStringSubmatch(line); len(matches) > 4 {
				hours, _ := strconv.Atoi(matches[1])
				minutes, _ := strconv.Atoi(matches[2])
				seconds, _ := strconv.Atoi(matches[3])
				centiseconds, _ := strconv.Atoi(matches[4])

				currentSeconds := float64(hours*3600+minutes*60+seconds) + float64(centiseconds)/100.0
				progressPercent := (currentSeconds / totalSeconds) * 100
				if progressPercent > 100 {
					progressPercent = 100
				}

				// Parse speed
				speed := 0.0
				if speedMatches := speedRegex.FindStringSubmatch(line); len(speedMatches) > 1 {
					speed, _ = strconv.ParseFloat(speedMatches[1], 64)
				}

				// Calculate ETA
				eta := ""
				if speed > 0 && currentSeconds < totalSeconds {
					remainingSeconds := (totalSeconds - currentSeconds) / speed
					eta = fmt.Sprintf(" (ETA: %s)", formatDuration(time.Duration(remainingSeconds)*time.Second))
				}

				// Display progress bar
				barWidth := 30
				filled := int((progressPercent / 100) * float64(barWidth))
				bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

				fmt.Printf("\r📊 [%s] %.1f%% - %.1fx speed%s", bar, progressPercent, speed, eta)
				progressShown = true
			}
		}
	}()

	// Wait for command to complete
	err = cmd.Wait()

	// Clear the progress line if we showed any
	if progressShown {
		fmt.Printf("\r%s\r", strings.Repeat(" ", 100))
	}

	if err != nil {
		return fmt.Errorf("ffmpeg execution failed: %w", err)
	}

	return nil
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	} else {
		return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
	}
}

// ExtractAudio extracts audio from a video file with specified parameters
func ExtractAudio(params AudioExtractionParams) error {
	// Step 1: Validate all parameters
	if err := validateAudioExtractionParams(params); err != nil {
		return err
	}

	// Step 2: Analyze input media
	mediaInfo, err := analyzeInputForAudioExtraction(params)
	if err != nil {
		return err
	}

	// Step 3: Select codec and build command
	codec, command, err := prepareAudioExtractionCommand(params, mediaInfo)
	if err != nil {
		return err
	}

	// Step 4: Display extraction info if verbose
	if params.Verbose {
		displayAudioExtractionInfo(params, codec, command)
	}

	// Step 5: Execute extraction
	return executeAudioExtraction(params, command, mediaInfo)
}

// validateAudioExtractionParams performs comprehensive validation of audio extraction parameters
func validateAudioExtractionParams(params AudioExtractionParams) error {
	// Validate input file
	if err := validateInputFile(params.InputFile); err != nil {
		return err
	}

	// Security validation for file paths
	if err := validateAudioExtractionPaths(params); err != nil {
		return err
	}

	// Security validation for audio parameters
	return validateAudioExtractionAudioParams(params)
}

// validateAudioExtractionPaths validates file paths for audio extraction
func validateAudioExtractionPaths(params AudioExtractionParams) error {
	if err := securityPolicy.ValidateFilePath(params.InputFile); err != nil {
		return fmt.Errorf("security validation failed for input path: %w", err)
	}

	if err := securityPolicy.ValidateFilePath(params.OutputFile); err != nil {
		return fmt.Errorf("security validation failed for output path: %w", err)
	}

	if err := securityPolicy.ValidateFileFormat(params.OutputFile); err != nil {
		return fmt.Errorf("security validation failed for output format: %w", err)
	}

	return nil
}

// validateAudioExtractionAudioParams validates audio-specific parameters
func validateAudioExtractionAudioParams(params AudioExtractionParams) error {
	if params.Codec != "" {
		if err := securityPolicy.ValidateCodec(params.Codec, "audio"); err != nil {
			return fmt.Errorf("security validation failed for audio codec: %w", err)
		}
	}

	if err := securityPolicy.ValidateBitrate(params.Bitrate); err != nil {
		return fmt.Errorf("security validation failed for bitrate: %w", err)
	}

	if params.SampleRate != "" {
		if err := validateSampleRate(params.SampleRate); err != nil {
			return fmt.Errorf("invalid sample rate: %w", err)
		}
	}

	if params.Channels != "" {
		if err := validateChannels(params.Channels); err != nil {
			return fmt.Errorf("invalid channels: %w", err)
		}
	}

	return nil
}

// analyzeInputForAudioExtraction analyzes the input media and validates audio streams
func analyzeInputForAudioExtraction(params AudioExtractionParams) (*analyzer.MediaInfo, error) {
	if params.Verbose {
		color.Cyan("🔍 Analyzing input media...")
	}

	mediaInfo, err := analyzer.AnalyzeMedia(params.InputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze input media: %w", err)
	}

	// Check if input has audio streams
	if len(mediaInfo.AudioStreams) == 0 {
		return nil, fmt.Errorf("no audio streams found in input file: %s", params.InputFile)
	}

	return mediaInfo, nil
}

// prepareAudioExtractionCommand selects codec and builds the FFmpeg command
func prepareAudioExtractionCommand(params AudioExtractionParams, mediaInfo *analyzer.MediaInfo) (string, []string, error) {
	// Determine output format and codec
	outputExt := strings.ToLower(filepath.Ext(params.OutputFile))
	codec, err := selectAudioCodec(outputExt, params.Codec)
	if err != nil {
		return "", nil, err
	}

	// Build FFmpeg command with security validation
	command := buildAudioExtractionCommandSecure(params, codec, mediaInfo)
	if command == nil {
		return "", nil, fmt.Errorf("failed to build secure audio extraction command")
	}

	return codec, command, nil
}

// displayAudioExtractionInfo shows extraction information in verbose mode
func displayAudioExtractionInfo(params AudioExtractionParams, codec string, command []string) {
	outputExt := strings.ToLower(filepath.Ext(params.OutputFile))
	
	fmt.Printf("🎵 Extracting audio to %s format\n", strings.TrimPrefix(outputExt, "."))
	fmt.Printf("🔧 Using codec: %s\n", codec)
	
	if params.Bitrate != "" || hasQualityBitrate(params.Quality) {
		bitrate := params.Bitrate
		if bitrate == "" {
			bitrate = getQualityBitrate(params.Quality)
		}
		fmt.Printf("📊 Bitrate: %s\n", bitrate)
	}
	
	fmt.Printf("Command: %s\n", strings.Join(command, " "))
	fmt.Println()
}

// executeAudioExtraction executes the audio extraction command
func executeAudioExtraction(params AudioExtractionParams, command []string, mediaInfo *analyzer.MediaInfo) error {
	if params.Verbose {
		color.Green("🚀 Starting audio extraction...")
	}

	cmd := exec.Command(command[0], command[1:]...)

	var err error
	if params.Verbose {
		// For verbose mode, show real-time progress
		err = executeFFmpegWithProgress(cmd, mediaInfo)
	} else {
		// For quiet mode, just run and wait
		output, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			err = fmt.Errorf("audio extraction failed: %w\nOutput: %s", cmdErr, string(output))
		}
	}

	if err != nil {
		return err
	}

	if params.Verbose {
		color.Green("✅ Audio extraction completed successfully!")
		fmt.Printf("Output saved to: %s\n", params.OutputFile)
	}

	return nil
}

// selectAudioCodec determines the appropriate audio codec for the output format
func selectAudioCodec(outputExt, customCodec string) (string, error) {
	// If user specified a codec, use it
	if customCodec != "" {
		return customCodec, nil
	}

	// Auto-select codec based on output format
	switch outputExt {
	case ".mp3":
		return "libmp3lame", nil
	case ".aac", ".m4a":
		return "aac", nil
	case ".wav":
		return "pcm_s16le", nil
	case ".flac":
		return "flac", nil
	case ".ogg":
		return "libvorbis", nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", outputExt)
	}
}

// buildAudioExtractionCommand builds the FFmpeg command for audio extraction
func buildAudioExtractionCommand(params AudioExtractionParams, codec string, mediaInfo *analyzer.MediaInfo) []string {
	command := []string{"ffmpeg", "-i", params.InputFile}

	// Disable video stream
	command = append(command, "-vn")

	// Set audio codec
	command = append(command, "-c:a", codec)

	// Set bitrate (custom or from quality preset)
	if params.Bitrate != "" {
		command = append(command, "-b:a", params.Bitrate)
	} else {
		// Apply quality preset bitrates
		bitrate := getQualityBitrate(params.Quality)
		if bitrate != "" {
			command = append(command, "-b:a", bitrate)
		}
	}

	// Set sample rate if specified
	if params.SampleRate != "" {
		command = append(command, "-ar", params.SampleRate)
	}

	// Set channels if specified
	if params.Channels != "" {
		command = append(command, "-ac", params.Channels)
	}

	// Set additional codec-specific options
	switch codec {
	case "libmp3lame":
		// For MP3, use VBR if no bitrate specified
		if params.Bitrate == "" && !hasQualityBitrate(params.Quality) {
			qualityValue := getMP3Quality(params.Quality)
			command = append(command, "-q:a", qualityValue)
		}
	case "flac":
		// For FLAC, set compression level
		compressionLevel := getFLACCompression(params.Quality)
		command = append(command, "-compression_level", compressionLevel)
	case "libvorbis":
		// For Vorbis, use quality-based encoding if no bitrate
		if params.Bitrate == "" && !hasQualityBitrate(params.Quality) {
			qualityValue := getVorbisQuality(params.Quality)
			command = append(command, "-q:a", qualityValue)
		}
	}

	// Output file (overwrite without asking)
	command = append(command, "-y", params.OutputFile)

	return command
}

// getQualityBitrate returns the bitrate for a quality preset
func getQualityBitrate(quality string) string {
	switch strings.ToLower(quality) {
	case "low":
		return "128k"
	case "medium":
		return "192k"
	case "high":
		return "320k"
	default:
		return "192k"
	}
}

// hasQualityBitrate checks if a quality preset should use bitrate-based encoding
func hasQualityBitrate(quality string) bool {
	// For some codecs, we prefer quality-based encoding over bitrate
	return true
}

// getMP3Quality returns the VBR quality setting for MP3
func getMP3Quality(quality string) string {
	switch strings.ToLower(quality) {
	case "low":
		return "5" // ~130 kbps
	case "medium":
		return "2" // ~190 kbps
	case "high":
		return "0" // ~245 kbps
	default:
		return "2"
	}
}

// getFLACCompression returns the compression level for FLAC
func getFLACCompression(quality string) string {
	switch strings.ToLower(quality) {
	case "low":
		return "0" // Fastest compression
	case "medium":
		return "5" // Balanced
	case "high":
		return "8" // Best compression
	default:
		return "5"
	}
}

// getVorbisQuality returns the VBR quality setting for Vorbis
func getVorbisQuality(quality string) string {
	switch strings.ToLower(quality) {
	case "low":
		return "3" // ~112 kbps
	case "medium":
		return "6" // ~192 kbps
	case "high":
		return "9" // ~320 kbps
	default:
		return "6"
	}
}

// validateSampleRate validates audio sample rate parameters
func validateSampleRate(sampleRate string) error {
	if sampleRate == "" {
		return nil
	}

	// Check for dangerous characters
	if strings.ContainsAny(sampleRate, ";|&`$<>\"'") {
		return fmt.Errorf("sample rate contains invalid characters: %s", sampleRate)
	}

	// Validate sample rate format
	validSampleRates := map[string]bool{
		"8000": true, "11025": true, "16000": true, "22050": true,
		"44100": true, "48000": true, "88200": true, "96000": true,
	}

	if !validSampleRates[sampleRate] {
		return fmt.Errorf("invalid sample rate: %s (valid: 8000, 11025, 16000, 22050, 44100, 48000, 88200, 96000)", sampleRate)
	}

	return nil
}

// validateChannels validates audio channel parameters
func validateChannels(channels string) error {
	if channels == "" {
		return nil
	}

	// Check for dangerous characters
	if strings.ContainsAny(channels, ";|&`$<>\"'") {
		return fmt.Errorf("channels contains invalid characters: %s", channels)
	}

	// Validate channel count
	validChannels := map[string]bool{
		"1": true, "2": true, "6": true, "8": true,
	}

	if !validChannels[channels] {
		return fmt.Errorf("invalid channel count: %s (valid: 1=mono, 2=stereo, 6=5.1, 8=7.1)", channels)
	}

	return nil
}

// buildAudioExtractionCommandSecure builds the FFmpeg command for audio extraction with security validation
func buildAudioExtractionCommandSecure(params AudioExtractionParams, codec string, mediaInfo *analyzer.MediaInfo) []string {
	command := []string{"ffmpeg", "-i", params.InputFile}

	// Disable video stream
	command = append(command, "-vn")

	// Set audio codec (already validated)
	command = append(command, "-c:a", codec)

	// Set bitrate (custom or from quality preset) - already validated
	if params.Bitrate != "" {
		command = append(command, "-b:a", params.Bitrate)
	} else {
		// Apply quality preset bitrates
		bitrate := getQualityBitrate(params.Quality)
		if bitrate != "" {
			command = append(command, "-b:a", bitrate)
		}
	}

	// Set sample rate if specified (already validated)
	if params.SampleRate != "" {
		command = append(command, "-ar", params.SampleRate)
	}

	// Set channels if specified (already validated)
	if params.Channels != "" {
		command = append(command, "-ac", params.Channels)
	}

	// Set additional codec-specific options (safe, predefined values only)
	switch codec {
	case "libmp3lame":
		// For MP3, use VBR if no bitrate specified
		if params.Bitrate == "" && !hasQualityBitrate(params.Quality) {
			qualityValue := getMP3Quality(params.Quality)
			command = append(command, "-q:a", qualityValue)
		}
	case "flac":
		// For FLAC, set compression level
		compressionLevel := getFLACCompression(params.Quality)
		command = append(command, "-compression_level", compressionLevel)
	case "libvorbis":
		// For Vorbis, use quality-based encoding if no bitrate
		if params.Bitrate == "" && !hasQualityBitrate(params.Quality) {
			qualityValue := getVorbisQuality(params.Quality)
			command = append(command, "-q:a", qualityValue)
		}
	}

	// Output file (overwrite without asking) - already validated
	command = append(command, "-y", params.OutputFile)

	return command
}
