package security

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// SecurityPolicy defines validation rules for user inputs
type SecurityPolicy struct {
	AllowedVideoCodecs map[string]bool
	AllowedAudioCodecs map[string]bool
	AllowedFormats     map[string]bool
	MaxPathLength      int
	MaxParameterLength int
}

// NewDefaultSecurityPolicy creates a security policy with safe defaults
func NewDefaultSecurityPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		AllowedVideoCodecs: map[string]bool{
			"libx264":    true,
			"libx265":    true,
			"libvpx-vp9": true,
			"libvpx":     true,
			"copy":       true,
		},
		AllowedAudioCodecs: map[string]bool{
			"aac":        true,
			"libopus":    true,
			"libmp3lame": true,
			"libvorbis":  true,
			"flac":       true,
			"pcm_s16le":  true,
			"copy":       true,
		},
		AllowedFormats: map[string]bool{
			"mp4":  true,
			"avi":  true,
			"mkv":  true,
			"webm": true,
			"mov":  true,
			"mp3":  true,
			"wav":  true,
			"aac":  true,
			"flac": true,
			"ogg":  true,
			"m4a":  true,
		},
		MaxPathLength:      255,
		MaxParameterLength: 50,
	}
}

// ValidateCodec validates video and audio codec parameters
func (p *SecurityPolicy) ValidateCodec(codec, codecType string) error {
	if len(codec) > p.MaxParameterLength {
		return fmt.Errorf("codec parameter too long (max %d characters)", p.MaxParameterLength)
	}

	// Check for dangerous characters that could enable command injection
	if containsDangerousChars(codec) {
		return fmt.Errorf("codec contains invalid characters: %s", codec)
	}

	var allowedCodecs map[string]bool
	switch codecType {
	case "video":
		allowedCodecs = p.AllowedVideoCodecs
	case "audio":
		allowedCodecs = p.AllowedAudioCodecs
	default:
		return fmt.Errorf("unknown codec type: %s", codecType)
	}

	if !allowedCodecs[codec] {
		return fmt.Errorf("codec not allowed: %s", codec)
	}

	return nil
}

// ValidateBitrate validates bitrate parameters
func (p *SecurityPolicy) ValidateBitrate(bitrate string) error {
	if bitrate == "" {
		return nil // Empty bitrate is allowed
	}

	if len(bitrate) > p.MaxParameterLength {
		return fmt.Errorf("bitrate parameter too long (max %d characters)", p.MaxParameterLength)
	}

	// Check for dangerous characters
	if containsDangerousChars(bitrate) {
		return fmt.Errorf("bitrate contains invalid characters: %s", bitrate)
	}

	// Validate bitrate format (e.g., "2M", "1500k", "192k")
	bitrateRegex := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?[kKmM]?$`)
	if !bitrateRegex.MatchString(bitrate) {
		return fmt.Errorf("invalid bitrate format: %s (use format like 2M, 1500k, 192k)", bitrate)
	}

	return nil
}

// ValidateResolution validates resolution parameters
func (p *SecurityPolicy) ValidateResolution(resolution string) error {
	if resolution == "" {
		return nil // Empty resolution is allowed
	}

	if len(resolution) > p.MaxParameterLength {
		return fmt.Errorf("resolution parameter too long (max %d characters)", p.MaxParameterLength)
	}

	// Check for dangerous characters
	if containsDangerousChars(resolution) {
		return fmt.Errorf("resolution contains invalid characters: %s", resolution)
	}

	// Validate resolution format (e.g., "1920x1080", "1280x720")
	resolutionRegex := regexp.MustCompile(`^[0-9]+x[0-9]+$`)
	if !resolutionRegex.MatchString(resolution) {
		return fmt.Errorf("invalid resolution format: %s (use format like 1920x1080)", resolution)
	}

	// Parse and validate reasonable resolution limits
	parts := strings.Split(resolution, "x")
	width, _ := strconv.Atoi(parts[0])
	height, _ := strconv.Atoi(parts[1])

	if width > 7680 || height > 4320 { // 8K max
		return fmt.Errorf("resolution too large: %s (max 7680x4320)", resolution)
	}

	if width < 1 || height < 1 {
		return fmt.Errorf("invalid resolution: %s (minimum 1x1)", resolution)
	}

	return nil
}

// ValidateFramerate validates framerate parameters
func (p *SecurityPolicy) ValidateFramerate(framerate string) error {
	if framerate == "" {
		return nil // Empty framerate is allowed
	}

	if len(framerate) > p.MaxParameterLength {
		return fmt.Errorf("framerate parameter too long (max %d characters)", p.MaxParameterLength)
	}

	// Check for dangerous characters
	if containsDangerousChars(framerate) {
		return fmt.Errorf("framerate contains invalid characters: %s", framerate)
	}

	// Validate framerate format (e.g., "30", "24", "60", "23.976")
	framerateRegex := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	if !framerateRegex.MatchString(framerate) {
		return fmt.Errorf("invalid framerate format: %s (use format like 30, 24, 60)", framerate)
	}

	// Parse and validate reasonable framerate limits
	fps, err := strconv.ParseFloat(framerate, 64)
	if err != nil {
		return fmt.Errorf("invalid framerate: %s", framerate)
	}

	if fps > 120 || fps <= 0 {
		return fmt.Errorf("framerate out of range: %s (must be between 0 and 120)", framerate)
	}

	return nil
}

// ValidateFilePath validates file paths to prevent directory traversal
func (p *SecurityPolicy) ValidateFilePath(path string) error {
	if len(path) > p.MaxPathLength {
		return fmt.Errorf("file path too long (max %d characters)", p.MaxPathLength)
	}

	// Clean the path and check for directory traversal attempts
	cleanPath := filepath.Clean(path)

	// Check for directory traversal patterns
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("directory traversal detected in path: %s", path)
	}

	// Check for dangerous characters in path
	if containsPathDangerousChars(path) {
		return fmt.Errorf("path contains invalid characters: %s", path)
	}

	return nil
}

// ValidateFileFormat validates file format based on extension
func (p *SecurityPolicy) ValidateFileFormat(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	if len(ext) > 1 {
		ext = ext[1:] // Remove the dot
	}

	if !p.AllowedFormats[ext] {
		return fmt.Errorf("file format not allowed: %s", ext)
	}

	return nil
}

// containsDangerousChars checks for characters that could enable command injection
func containsDangerousChars(input string) bool {
	// Characters that could be used for command injection
	dangerousChars := []string{
		";", "&", "|", "`", "$", "(", ")", "{", "}", "[", "]",
		"<", ">", "\\", "\"", "'", "\n", "\r", "\t",
	}

	for _, char := range dangerousChars {
		if strings.Contains(input, char) {
			return true
		}
	}

	return false
}

// containsPathDangerousChars checks for dangerous characters in file paths
func containsPathDangerousChars(path string) bool {
	// Characters that could be dangerous in file paths
	dangerousChars := []string{
		";", "&", "|", "`", "$", "\"", "'", "\n", "\r", "\t",
	}

	for _, char := range dangerousChars {
		if strings.Contains(path, char) {
			return true
		}
	}

	return false
}

// SanitizeCodecParameters safely parses codec parameters
func (p *SecurityPolicy) SanitizeCodecParameters(codec string, codecType string) (string, []string, error) {
	// Validate the codec first
	if err := p.ValidateCodec(codec, codecType); err != nil {
		return "", nil, err
	}

	// For whitelisted codecs, we only allow the codec name itself
	// No additional parameters to prevent injection
	return codec, []string{}, nil
}
