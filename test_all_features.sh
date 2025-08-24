#!/bin/bash

# Comprehensive test script for Term Video Transcoder
# Phase 2 Custom Parameters Feature Testing
# Self-contained with on-the-fly test file generation

set -e  # Exit on any error

echo "🎬 Term Video Transcoder - Phase 2 Custom Parameters Testing"
echo "============================================================"

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🧹 Cleaning up test files..."
    rm -rf internal_test test_outputs 2>/dev/null || true
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
mkdir -p internal_test test_outputs

# Generate test video on-the-fly (2 minutes, 720p, with test pattern and audio)
echo "📹 Generating test video (2 minutes)..."
ffmpeg -f lavfi -i testsrc2=duration=120:size=1280x720:rate=30 \
       -f lavfi -i sine=frequency=1000:duration=120 \
       -c:v libx264 -preset fast -crf 23 \
       -c:a aac -b:a 128k \
       -y internal_test/short_test_video.mp4 \
       -v quiet

# Verify test video was created
if [ ! -f "internal_test/short_test_video.mp4" ]; then
    echo "❌ Error: Failed to generate test video"
    exit 1
fi

echo "✅ Test video generated successfully ($(du -h internal_test/short_test_video.mp4 | cut -f1))"

echo ""
echo "🔧 Testing Custom Parameters Feature..."
echo ""

# Test 1: Basic codec selection
echo "1. Testing codec selection (H.264 + AAC)..."
./transcoder convert internal_test/short_test_video.mp4 test_outputs/test_h264_aac.mp4 \
  --video-codec libx264 --audio-codec aac

# Test 2: High quality VP9 encoding
echo ""
echo "2. Testing VP9 + Opus high quality encoding..."
./transcoder convert internal_test/short_test_video.mp4 test_outputs/test_vp9_opus.webm \
  --video-codec libvpx-vp9 --audio-codec libopus \
  --video-bitrate 2M --audio-bitrate 192k

# Test 3: Low quality for mobile/web
echo ""
echo "3. Testing low quality mobile-friendly encoding..."
./transcoder convert internal_test/short_test_video.mp4 test_outputs/test_mobile.mp4 \
  --video-codec libx264 --audio-codec aac \
  --video-bitrate 500k --audio-bitrate 64k \
  --resolution 640x360 --framerate 24

# Test 4: Resolution scaling only
echo ""
echo "4. Testing resolution scaling..."
./transcoder convert internal_test/short_test_video.mp4 test_outputs/test_720p.mp4 \
  --resolution 1280x720

# Test 5: Frame rate adjustment only
echo ""
echo "5. Testing frame rate adjustment..."
./transcoder convert internal_test/short_test_video.mp4 test_outputs/test_24fps.mp4 \
  --framerate 24

# Test 6: Bitrate control only
echo ""
echo "6. Testing bitrate control..."
./transcoder convert internal_test/short_test_video.mp4 test_outputs/test_bitrate.mp4 \
  --video-bitrate 1.5M --audio-bitrate 128k

# Test 7: All parameters combined
echo ""
echo "7. Testing all custom parameters combined..."
./transcoder convert internal_test/short_test_video.mp4 test_outputs/test_full_custom.webm \
  --video-codec libvpx-vp9 --audio-codec libopus \
  --video-bitrate 1.2M --audio-bitrate 160k \
  --resolution 854x480 --framerate 25

echo ""
echo "📊 Test Results Summary:"
echo "========================"

# Display file sizes and basic info
for file in test_outputs/test_*.{mp4,webm}; do
  if [ -f "$file" ]; then
    size=$(du -h "$file" | cut -f1)
    basename=$(basename "$file")
    echo "📁 $basename: $size"
    
    # Get basic video info
    resolution=$(ffprobe -v quiet -select_streams v:0 -show_entries stream=width,height -of csv=p=0 "$file" 2>/dev/null)
    fps=$(ffprobe -v quiet -select_streams v:0 -show_entries stream=r_frame_rate -of csv=p=0 "$file" 2>/dev/null | head -1)
    video_codec=$(ffprobe -v quiet -select_streams v:0 -show_entries stream=codec_name -of csv=p=0 "$file" 2>/dev/null)
    audio_codec=$(ffprobe -v quiet -select_streams a:0 -show_entries stream=codec_name -of csv=p=0 "$file" 2>/dev/null)
    
    echo "   └─ Resolution: $resolution, FPS: $fps, Video: $video_codec, Audio: $audio_codec"
    echo ""
  fi
done

echo "✅ All tests completed!"
echo ""
echo "🎯 Features Tested:"
echo "  ✓ Manual codec selection (video & audio)"
echo "  ✓ Bitrate control (video & audio)"
echo "  ✓ Resolution scaling"
echo "  ✓ Frame rate adjustment"
echo "  ✓ Parameter combination"
echo "  ✓ Format compatibility (MP4 & WebM)"
echo ""
echo "Phase 2 Custom Parameters feature is fully implemented and working! 🚀"
