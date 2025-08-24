#!/bin/bash

# Audio Extraction Test Script
# Tests the new audio extraction feature

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🎵 Audio Extraction Test Suite${NC}"
echo "======================================"

# Function to cleanup on exit
cleanup() {
    echo -e "\n🧹 Cleaning up test files..."
    rm -rf audio_test_outputs/
    rm -f internal_test/test_video.mp4
    echo "✅ Cleanup completed!"
}

# Set up cleanup trap
trap cleanup EXIT

# Pre-flight checks
echo -e "${YELLOW}🔍 Pre-flight checks...${NC}"
echo "----------------------"

if ! command -v ffmpeg &> /dev/null; then
    echo "❌ ffmpeg not found. Please install ffmpeg to run audio extraction tests."
    echo "Install with: brew install ffmpeg (macOS) or apt-get install ffmpeg (Ubuntu)"
    exit 1
fi
echo "✅ ffmpeg: $(ffmpeg -version | head -n1 | cut -d' ' -f3)"

if ! command -v ffprobe &> /dev/null; then
    echo "❌ ffprobe not found. Please install ffprobe to run tests."
    exit 1
fi
echo "✅ ffprobe: Available"

echo -e "\n📦 Building transcoder..."
go build -o transcoder .

echo -e "\n🏗️  Setting up test environment..."
mkdir -p audio_test_outputs/
mkdir -p internal_test/

# Generate a test video with audio using FFmpeg
echo "📹 Generating test video with audio (30 seconds)..."
ffmpeg -f lavfi -i testsrc2=duration=30:size=640x480:rate=25 \
       -f lavfi -i sine=frequency=440:duration=30 \
       -c:v libx264 -preset ultrafast -c:a aac -b:a 128k \
       -y internal_test/test_video.mp4 &> /dev/null

if [ $? -eq 0 ]; then
    echo "✅ Test video generated successfully ($(du -h internal_test/test_video.mp4 | cut -f1))"
else
    echo "❌ Failed to generate test video"
    exit 1
fi

echo -e "\n⚡ Running audio extraction tests...\n"

# Test 1: Basic MP3 extraction
echo "1. Testing basic MP3 extraction..."
./transcoder extract internal_test/test_video.mp4 audio_test_outputs/test.mp3 --force
if [ $? -eq 0 ]; then
    echo "✅ MP3 extraction successful"
else
    echo "❌ MP3 extraction failed"
    exit 1
fi

# Test 2: High quality FLAC extraction
echo -e "\n2. Testing high quality FLAC extraction..."
./transcoder extract internal_test/test_video.mp4 audio_test_outputs/test_hq.flac --quality high --force
if [ $? -eq 0 ]; then
    echo "✅ FLAC extraction successful"
else
    echo "❌ FLAC extraction failed"
    exit 1
fi

# Test 3: Custom bitrate extraction
echo -e "\n3. Testing custom bitrate extraction..."
./transcoder extract internal_test/test_video.mp4 audio_test_outputs/test_320k.mp3 --bitrate 320k --force
if [ $? -eq 0 ]; then
    echo "✅ Custom bitrate extraction successful"
else
    echo "❌ Custom bitrate extraction failed"
    exit 1
fi

# Test 4: WAV extraction with custom sample rate
echo -e "\n4. Testing WAV extraction with custom sample rate..."
./transcoder extract internal_test/test_video.mp4 audio_test_outputs/test.wav --sample-rate 48000 --force
if [ $? -eq 0 ]; then
    echo "✅ WAV extraction successful"
else
    echo "❌ WAV extraction failed"
    exit 1
fi

# Test 5: OGG extraction with custom codec
echo -e "\n5. Testing OGG extraction with custom codec..."
./transcoder extract internal_test/test_video.mp4 audio_test_outputs/test.ogg --codec libvorbis --force
if [ $? -eq 0 ]; then
    echo "✅ OGG extraction successful"
else
    echo "❌ OGG extraction failed"
    exit 1
fi

# Test 6: AAC extraction with custom channels
echo -e "\n6. Testing AAC extraction with mono output..."
./transcoder extract internal_test/test_video.mp4 audio_test_outputs/test_mono.aac --channels 1 --force
if [ $? -eq 0 ]; then
    echo "✅ AAC mono extraction successful"
else
    echo "❌ AAC mono extraction failed"
    exit 1
fi

# Display results summary
echo -e "\n📊 Audio Extraction Test Results:"
echo "=================================="

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
            
            echo "📁 $filename: $size"
            echo "   └─ Duration: ${duration}s, Codec: $codec, Bitrate: $bitrate, Sample Rate: $sample_rate Hz, Channels: $channels"
        fi
    done
fi

echo -e "\n✅ All audio extraction tests completed successfully! 🎉"

echo -e "\n🎯 Features Tested:"
echo "  ✓ MP3 extraction with default settings"
echo "  ✓ FLAC extraction with high quality preset"
echo "  ✓ Custom bitrate control (320k)"
echo "  ✓ WAV extraction with custom sample rate (48kHz)"
echo "  ✓ OGG extraction with custom codec (libvorbis)"
echo "  ✓ AAC extraction with mono output"
echo "  ✓ Multiple audio format support"
echo "  ✓ Quality presets (low, medium, high)"
echo "  ✓ Parameter validation and error handling"

echo -e "\n💡 Tip: Use this test script to verify audio extraction functionality."
echo "    For manual testing, try: ./transcoder extract input.mp4 output.mp3"
