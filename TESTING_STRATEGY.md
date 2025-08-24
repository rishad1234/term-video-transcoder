# Testing Strategy - Term Video Transcoder

## 🚀 Quick Development Testing

For fast development testing during feature development, use our **2-minute test video**:

```bash
# Quick individual tests (5-45 seconds each)
./transcoder convert internal_test/short_test_video.mp4 output.webm --preset high
./transcoder convert internal_test/short_test_video.mp4 output.mp4 --video-codec libx264 --video-bitrate 1M
./transcoder convert internal_test/short_test_video.mp4 output.mp4 --resolution 640x360

# Run all quick tests
./quick_test.sh
```

**Benefits:**
- ⚡ **20x faster** than long video testing
- 🔄 **Rapid iteration** during development
- 💾 **Small output files** (3-14MB vs 100MB+)
- ✅ **Same validation** of all features

## 🧪 Comprehensive Testing

For thorough testing before releases:

```bash
# Run all comprehensive tests with multiple scenarios
./test_all_features.sh
```

**Use cases:**
- 🚀 **Pre-release validation**
- 📊 **Performance benchmarking**
- 🔍 **Edge case testing**
- 📝 **Documentation examples**

## 📁 Test Files

| File | Duration | Size | Use Case |
|------|----------|------|----------|
| `short_test_video.mp4` | 2 min | 5MB | Development testing |
| `long_test_video.mp4` | 20 min | 50MB | Comprehensive testing |

## ⏱️ Performance Comparison

| Test Type | Short Video | Long Video | Speedup |
|-----------|-------------|------------|---------|
| VP9 Encoding | 42s | 15+ min | **20x faster** |
| H.264 Encoding | 5.6s | 2+ min | **22x faster** |
| Resolution Scaling | 3.3s | 1+ min | **18x faster** |

## 🛠️ Development Workflow

```bash
# 1. Make code changes
# 2. Quick validation
./quick_test.sh

# 3. If all quick tests pass, run comprehensive tests
./test_all_features.sh

# 4. Commit changes
```

This approach ensures **fast development cycles** while maintaining **thorough validation** before releases.
