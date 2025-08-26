package analyzer

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// MediaInfo holds comprehensive information about a media file
type MediaInfo struct {
	Filename     string        `json:"filename"`
	Format       string        `json:"format"`
	Duration     time.Duration `json:"duration"`
	Size         int64         `json:"size"`
	Bitrate      int64         `json:"bitrate"`
	VideoStreams []VideoStream `json:"video_streams"`
	AudioStreams []AudioStream `json:"audio_streams"`
}

// VideoStream represents a video stream in the media file
type VideoStream struct {
	Index       int    `json:"index"`
	Codec       string `json:"codec"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	FrameRate   string `json:"frame_rate"`
	PixelFormat string `json:"pixel_format"`
	Bitrate     int64  `json:"bitrate"`
}

// AudioStream represents an audio stream in the media file
type AudioStream struct {
	Index      int    `json:"index"`
	Codec      string `json:"codec"`
	SampleRate int    `json:"sample_rate"`
	Channels   int    `json:"channels"`
	Bitrate    int64  `json:"bitrate"`
	Language   string `json:"language"`
}

// AnalyzeMedia uses ffprobe to extract comprehensive media information
func AnalyzeMedia(filepath string) (*MediaInfo, error) {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filepath)
	}

	// Run ffprobe command
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filepath)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	return parseFFProbeOutput(string(output), filepath)
}

// parseFFProbeOutput parses the JSON output from ffprobe
func parseFFProbeOutput(jsonOutput, filepath string) (*MediaInfo, error) {
	info := &MediaInfo{
		Filename: filepath,
	}

	if err := parseFormatInformation(jsonOutput, info); err != nil {
		return nil, fmt.Errorf("parsing format information: %w", err)
	}

	if err := parseStreamInformation(jsonOutput, info); err != nil {
		return nil, fmt.Errorf("parsing stream information: %w", err)
	}

	return info, nil
}

// parseFormatInformation extracts format-level metadata
func parseFormatInformation(jsonOutput string, info *MediaInfo) error {
	format := gjson.Get(jsonOutput, "format")
	if !format.Exists() {
		return nil
	}

	info.Format = format.Get("format_name").String()
	parseDuration(format, info)
	parseSize(format, info)
	parseBitrate(format, info)

	return nil
}

// parseDuration extracts and converts duration from format metadata
func parseDuration(format gjson.Result, info *MediaInfo) {
	if durationStr := format.Get("duration").String(); durationStr != "" {
		if duration, err := strconv.ParseFloat(durationStr, 64); err == nil {
			info.Duration = time.Duration(duration * float64(time.Second))
		}
	}
}

// parseSize extracts file size from format metadata
func parseSize(format gjson.Result, info *MediaInfo) {
	if sizeStr := format.Get("size").String(); sizeStr != "" {
		if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
			info.Size = size
		}
	}
}

// parseBitrate extracts overall bitrate from format metadata
func parseBitrate(format gjson.Result, info *MediaInfo) {
	if bitrateStr := format.Get("bit_rate").String(); bitrateStr != "" {
		if bitrate, err := strconv.ParseInt(bitrateStr, 10, 64); err == nil {
			info.Bitrate = bitrate
		}
	}
}

// parseStreamInformation extracts all stream information
func parseStreamInformation(jsonOutput string, info *MediaInfo) error {
	streams := gjson.Get(jsonOutput, "streams").Array()
	for _, stream := range streams {
		codecType := stream.Get("codec_type").String()

		switch codecType {
		case "video":
			parseVideoStream(stream, info)
		case "audio":
			parseAudioStream(stream, info)
		}
	}
	return nil
}

// parseVideoStream extracts video stream metadata
func parseVideoStream(stream gjson.Result, info *MediaInfo) {
	videoStream := VideoStream{
		Index:       int(stream.Get("index").Int()),
		Codec:       stream.Get("codec_name").String(),
		Width:       int(stream.Get("width").Int()),
		Height:      int(stream.Get("height").Int()),
		FrameRate:   stream.Get("r_frame_rate").String(),
		PixelFormat: stream.Get("pix_fmt").String(),
	}

	parseStreamBitrate(stream, &videoStream.Bitrate)
	info.VideoStreams = append(info.VideoStreams, videoStream)
}

// parseAudioStream extracts audio stream metadata
func parseAudioStream(stream gjson.Result, info *MediaInfo) {
	audioStream := AudioStream{
		Index:      int(stream.Get("index").Int()),
		Codec:      stream.Get("codec_name").String(),
		SampleRate: int(stream.Get("sample_rate").Int()),
		Channels:   int(stream.Get("channels").Int()),
		Language:   stream.Get("tags.language").String(),
	}

	parseStreamBitrate(stream, &audioStream.Bitrate)
	info.AudioStreams = append(info.AudioStreams, audioStream)
}

// parseStreamBitrate extracts bitrate for individual streams
func parseStreamBitrate(stream gjson.Result, bitrate *int64) {
	if bitrateStr := stream.Get("bit_rate").String(); bitrateStr != "" {
		if parsedBitrate, err := strconv.ParseInt(bitrateStr, 10, 64); err == nil {
			*bitrate = parsedBitrate
		}
	}
}

// CheckFFProbe verifies that ffprobe is available in the system
func CheckFFProbe() error {
	cmd := exec.Command("ffprobe", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffprobe not found or not working: %w", err)
	}
	return nil
}

// CheckFFMpeg verifies that ffmpeg is available in the system
func CheckFFMpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg not found or not working: %w", err)
	}
	return nil
}
