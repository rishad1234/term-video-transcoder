#!/bin/bash

# Quick Test Script for Term Video Transcoder
# Self-contained testing with on-the-fly test file generation

set -e  # Exit on any error

echo "🚀 Quick Test - Term Video Transcoder"
echo "====================================="

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🧹 Cleaning up test files..."
    rm -rf internal_test quick_test_outputs 2>/dev/null || true
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
mkdir -p internal_test quick_test_outputs

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
echo "⚡ Running quick tests with 2-minute video..."
echo ""

# Quick test 1: Basic preset conversion
echo "1. Testing basic preset conversion..."
time ./transcoder convert internal_test/short_test_video.mp4 quick_test_outputs/preset_test.webm --preset high

# Quick test 2: Custom parameters
echo ""
echo "2. Testing custom parameters..."
time ./transcoder convert internal_test/short_test_video.mp4 quick_test_outputs/custom_test.mp4 \
  --video-codec libx264 --audio-codec aac --video-bitrate 1M --audio-bitrate 128k

# Quick test 3: Resolution scaling
echo ""
echo "3. Testing resolution scaling..."
time ./transcoder convert internal_test/short_test_video.mp4 quick_test_outputs/resolution_test.mp4 \
  --resolution 640x360

echo ""
echo "📊 Quick Test Results:"
echo "====================="

for file in quick_test_outputs/*.{mp4,webm}; do
  if [ -f "$file" ]; then
    size=$(du -h "$file" | cut -f1)
    basename=$(basename "$file")
    echo "📁 $basename: $size"
  fi
done

echo ""
echo "✅ Quick tests completed! All features working correctly. 🎉"
echo ""
echo "💡 Tip: Use this script for fast development testing."
echo "    For comprehensive testing, run: ./test_all_features.sh"
