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

	// Parse format information
	format := gjson.Get(jsonOutput, "format")
	if format.Exists() {
		info.Format = format.Get("format_name").String()
		
		// Parse duration
		if durationStr := format.Get("duration").String(); durationStr != "" {
			if duration, err := strconv.ParseFloat(durationStr, 64); err == nil {
				info.Duration = time.Duration(duration * float64(time.Second))
			}
		}
		
		// Parse size
		if sizeStr := format.Get("size").String(); sizeStr != "" {
			if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
				info.Size = size
			}
		}
		
		// Parse bitrate
		if bitrateStr := format.Get("bit_rate").String(); bitrateStr != "" {
			if bitrate, err := strconv.ParseInt(bitrateStr, 10, 64); err == nil {
				info.Bitrate = bitrate
			}
		}
	}

	// Parse streams
	streams := gjson.Get(jsonOutput, "streams").Array()
	for _, stream := range streams {
		codecType := stream.Get("codec_type").String()
		
		switch codecType {
		case "video":
			videoStream := VideoStream{
				Index:       int(stream.Get("index").Int()),
				Codec:       stream.Get("codec_name").String(),
				Width:       int(stream.Get("width").Int()),
				Height:      int(stream.Get("height").Int()),
				FrameRate:   stream.Get("r_frame_rate").String(),
				PixelFormat: stream.Get("pix_fmt").String(),
			}
			
			if bitrateStr := stream.Get("bit_rate").String(); bitrateStr != "" {
				if bitrate, err := strconv.ParseInt(bitrateStr, 10, 64); err == nil {
					videoStream.Bitrate = bitrate
				}
			}
			
			info.VideoStreams = append(info.VideoStreams, videoStream)
			
		case "audio":
			audioStream := AudioStream{
				Index:      int(stream.Get("index").Int()),
				Codec:      stream.Get("codec_name").String(),
				SampleRate: int(stream.Get("sample_rate").Int()),
				Channels:   int(stream.Get("channels").Int()),
				Language:   stream.Get("tags.language").String(),
			}
			
			if bitrateStr := stream.Get("bit_rate").String(); bitrateStr != "" {
				if bitrate, err := strconv.ParseInt(bitrateStr, 10, 64); err == nil {
					audioStream.Bitrate = bitrate
				}
			}
			
			info.AudioStreams = append(info.AudioStreams, audioStream)
		}
	}

	return info, nil
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
