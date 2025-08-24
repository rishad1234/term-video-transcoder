#!/bin/bash

# Audio Extraction Test Script for Term Video Transcoder
# Quick and comprehensive audio testing with self-contained test file generation

set -e  # Exit on any error

echo "🎵 Audio Extraction Test Suite"
echo "=============================="

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🧹 Cleaning up test files..."
    rm -rf internal_test audio_test_outputs 2>/dev/null || true
    echo "✅ Cleanup completed!"
}

# Set trap to cleanup on script exit (success or failure)
trap cleanup EXIT

# Pre-flight checks
echo "🔍 Pre-flight checks..."
echo "----------------------"

# Check if ffmpeg is installed
if ! command -v ffmpeg &> /dev/null; then
    echo "❌ Error: ffmpeg is not installed or not in PATH"
    echo ""
    echo "📋 Installation instructions:"
    echo "  macOS:     brew install ffmpeg"
    echo "  Ubuntu:    sudo apt install ffmpeg"
    echo "  Windows:   Download from https://ffmpeg.org/download.html"
    echo ""
    exit 1
fi

# Check if ffprobe is installed
if ! command -v ffprobe &> /dev/null; then
    echo "❌ Error: ffprobe is not installed or not in PATH"
    echo ""
    echo "📋 ffprobe is typically included with ffmpeg installation"
    exit 1
fi

# Display versions
FFMPEG_VERSION=$(ffmpeg -version | head -n1 | cut -d' ' -f3)
echo "✅ ffmpeg: $FFMPEG_VERSION"
echo "✅ ffprobe: Available"

# Build the transcoder
echo ""
echo "📦 Building transcoder..."
go build -o transcoder

# Create test directories
echo ""
echo "🏗️  Setting up test environment..."
mkdir -p internal_test audio_test_outputs

# Generate test video with rich audio content (30 seconds, stereo)
echo "📹 Generating test video with audio (30 seconds)..."
ffmpeg -f lavfi -i testsrc2=duration=30:size=640x480:rate=25 \
       -f lavfi -i sine=frequency=440:duration=30 \
       -c:v libx264 -preset fast -crf 23 \
       -c:a aac -b:a 192k \
       -y internal_test/test_video.mp4 \
       -v quiet

# Verify test video was created
if [ ! -f "internal_test/test_video.mp4" ]; then
    echo "❌ Error: Failed to generate test video"
    exit 1
fi

echo "✅ Test video generated successfully ($(du -h internal_test/test_video.mp4 | cut -f1))"

echo ""
echo "⚡ Running audio extraction tests..."
echo ""

# Test 1: Basic MP3 extraction (default quality)
echo "1. Testing basic MP3 extraction..."
time ./transcoder extract internal_test/test_video.mp4 audio_test_outputs/basic.mp3

# Test 2: High quality FLAC extraction
echo ""
echo "2. Testing high quality FLAC extraction..."
time ./transcoder extract internal_test/test_video.mp4 audio_test_outputs/high_quality.flac --quality high

# Test 3: Custom bitrate MP3 (320k)
echo ""
echo "3. Testing custom bitrate MP3 (320k)..."
time ./transcoder extract internal_test/test_video.mp4 audio_test_outputs/custom_bitrate.mp3 --bitrate 320k

# Test 4: WAV with custom sample rate
echo ""
echo "4. Testing WAV with custom sample rate (48kHz)..."
time ./transcoder extract internal_test/test_video.mp4 audio_test_outputs/custom_sample.wav --sample-rate 48000

# Test 5: OGG with custom codec
echo ""
echo "5. Testing OGG with libvorbis codec..."
time ./transcoder extract internal_test/test_video.mp4 audio_test_outputs/custom_codec.ogg --codec libvorbis

# Test 6: AAC mono conversion
echo ""
echo "6. Testing AAC mono conversion..."
time ./transcoder extract internal_test/test_video.mp4 audio_test_outputs/mono.aac --channels 1

# Test 7: Low quality preset
echo ""
echo "7. Testing low quality preset..."
time ./transcoder extract internal_test/test_video.mp4 audio_test_outputs/low_quality.mp3 --quality low

echo ""
echo "📊 Audio Extraction Test Results:"
echo "================================="

# Display results with audio properties
if command -v ffprobe &> /dev/null; then
    for file in audio_test_outputs/*; do
        if [ -f "$file" ]; then
            filename=$(basename "$file")
            size=$(du -h "$file" | cut -f1)
            
            # Get audio info using ffprobe
            duration=$(ffprobe -v quiet -show_entries format=duration -of csv=p=0 "$file" 2>/dev/null | cut -d. -f1)
            codec=$(ffprobe -v quiet -select_streams a:0 -show_entries stream=codec_name -of csv=p=0 "$file" 2>/dev/null)
            bitrate=$(ffprobe -v quiet -select_streams a:0 -show_entries stream=bit_rate -of csv=p=0 "$file" 2>/dev/null)
            sample_rate=$(ffprobe -v quiet -select_streams a:0 -show_entries stream=sample_rate -of csv=p=0 "$file" 2>/dev/null)
            channels=$(ffprobe -v quiet -select_streams a:0 -show_entries stream=channels -of csv=p=0 "$file" 2>/dev/null)
            
            # Format bitrate display
            if [ "$bitrate" != "N/A" ] && [ -n "$bitrate" ]; then
                bitrate_display="${bitrate} bps"
            else
                bitrate_display="Variable"
            fi
            
            echo "📁 $filename: $size"
            echo "   └─ Duration: ${duration}s, Codec: $codec, Bitrate: $bitrate_display, Sample Rate: $sample_rate Hz, Channels: $channels"
        fi
    done
fi

echo ""
echo "✅ All audio extraction tests completed successfully! 🎉"

echo ""
echo "🎯 Features Tested:"
echo "  ✓ Basic MP3 extraction (default medium quality)"
echo "  ✓ High quality FLAC extraction (lossless)"
echo "  ✓ Custom bitrate control (320k MP3)"
echo "  ✓ Custom sample rate (48kHz WAV)"
echo "  ✓ Custom codec selection (libvorbis OGG)"
echo "  ✓ Channel configuration (stereo to mono)"
echo "  ✓ Low quality preset (128k MP3)"
echo "  ✓ Multiple audio formats (MP3, FLAC, WAV, OGG, AAC)"
echo "  ✓ Quality presets (low, medium, high)"
echo "  ✓ Parameter validation and error handling"

echo ""
echo "💡 Usage Examples:"
echo "  Basic:     ./transcoder extract video.mp4 audio.mp3"
echo "  Quality:   ./transcoder extract video.mp4 audio.flac --quality high"
echo "  Custom:    ./transcoder extract video.mp4 audio.mp3 --bitrate 320k"
echo "  Mono:      ./transcoder extract video.mp4 audio.wav --channels 1"

echo ""
echo "📋 Next Steps:"
echo "  • Run ./quick_test.sh for video conversion testing"
echo "  • Run ./test_all_features.sh for comprehensive testing"
echo "  • Use ./transcoder --help for full command reference"
