# gozxing A Barcode Scanning/Encoding Library for Go

[![Build Status](https://github.com/mechiko/gozxing/actions/workflows/main.yml/badge.svg)](https://github.com/mechiko/gozxing/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/mechiko/gozxing/branch/master/graph/badge.svg)](https://codecov.io/gh/mechiko/gozxing)

[ZXing](https://github.com/zxing/zxing) is an open-source, multi-format 1D/2D barcode image processing library for Java.
This project is a port of ZXing core library to pure Go.

## Porting Status (supported formats)

### 2D barcodes

| Format      | Scanning           | Encoding           |
|-------------|--------------------|--------------------|
| QR Code     | :heavy_check_mark: | :heavy_check_mark: |
| Data Matrix | :heavy_check_mark: | :heavy_check_mark: |
| Aztec       | :heavy_check_mark: |                    |
| PDF 417     |                    |                    |
| MaxiCode    |                    |                    |


### 1D product barcodes

| Format      | Scanning           | Encoding           |
|-------------|--------------------|--------------------|
| UPC-A       | :heavy_check_mark: | :heavy_check_mark: |
| UPC-E       | :heavy_check_mark: | :heavy_check_mark: |
| EAN-8       | :heavy_check_mark: | :heavy_check_mark: |
| EAN-13      | :heavy_check_mark: | :heavy_check_mark: |

### 1D industrial barcode

| Format       | Scanning           | Encoding           |
|--------------|--------------------|--------------------|
| Code 39      | :heavy_check_mark: | :heavy_check_mark: |
| Code 93      | :heavy_check_mark: | :heavy_check_mark: |
| Code 128     | :heavy_check_mark: | :heavy_check_mark: |
| Codabar      | :heavy_check_mark: | :heavy_check_mark: |
| ITF          | :heavy_check_mark: | :heavy_check_mark: |
| RSS-14       | :heavy_check_mark: | -                  |
| RSS-Expanded |                    |                    |

### Special reader/writer

| Reader/Writer                | Porting status     |
|------------------------------|--------------------|
| MultiFormatReader            |                    |
| MultiFormatWriter            |                    |
| ByQuadrantReader             |                    |
| GenericMultipleBarcodeReader |                    |
| QRCodeMultiReader            | :heavy_check_mark: |
| MultiFormatUPCEANReader      | :heavy_check_mark: |
| MultiFormatOneDReader        |                    |

## Usage Examples

### Scanning QR code

```Go
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"

	"github.com/mechiko/gozxing"
	"github.com/mechiko/gozxing/qrcode"
)

func main() {
	// open and decode image file
	file, _ := os.Open("qrcode.jpg")
	img, _, _ := image.Decode(file)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)

	fmt.Println(result)
}
```

### Generating CODE128 barcode

```Go
package main

import (
	"image/png"
	"os"

	"github.com/mechiko/gozxing"
	"github.com/mechiko/gozxing/oned"
)

func main() {
	// Generate a barcode image (*BitMatrix)
	enc := oned.NewCode128Writer()
	img, _ := enc.Encode("Hello, Gophers!", gozxing.BarcodeFormat_CODE_128, 250, 50, nil)

	file, _ := os.Create("barcode.png")
	defer file.Close()

	// *BitMatrix implements the image.Image interface,
	// so it is able to be passed to png.Encode directly.
	_ = png.Encode(file, img)
}
```

## Thread Safety

Starting from version v0.1.2, BinaryBitmap and HybridBinarizer are thread-safe for concurrent access. Multiple goroutines can safely call GetBlackMatrix() on the same instance without external synchronization.

### Why This Matters

Creating a BinaryBitmap from an image can be computationally significant. For high-performance services, caching these preprocessed bitmaps can offer tangible benefits.

### Performance Benchmarks

The following benchmark demonstrates why caching BinaryBitmap instances is valuable:

```go
func BenchmarkCachingImpact(b *testing.B) {
    // Generate a 400x400 QR code
    key, _ := totp.Generate(totp.GenerateOpts{
        Issuer:      "BenchmarkApp",
        AccountName: "bench@example.com",
    })
    img, _ := key.Image(400, 400)
    
    b.Run("WithoutCaching", func(b *testing.B) {
        reader := qrcode.NewQRCodeReader()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            // Create new BinaryBitmap each time (expensive!)
            bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
            _, _ = reader.Decode(bmp, nil)
        }
    })
    
    b.Run("WithCaching", func(b *testing.B) {
        // Create BinaryBitmap once and reuse
        bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
        reader := qrcode.NewQRCodeReader()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            // Reuse the same BinaryBitmap (fast!)
            _, _ = reader.Decode(bmp, nil)
        }
    })
}
```

Results on Apple M1 Ultra:
```
BenchmarkCachingImpact/WithoutCaching-20     498    2392097 ns/op
BenchmarkCachingImpact/WithCaching-20       5311     188575 ns/op
```

This shows a **12.7x performance improvement** when caching BinaryBitmap instances.

For a pseudocode example of how you might leverage the added thread-safety of `gozxing`'s BinaryBitmap now:

```go
import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "image/jpeg"
    "sync"
    
    "github.com/mechiko/gozxing"
    "github.com/mechiko/gozxing/qrcode"
)

type QRResult struct {
    Text string
}

// High-performance QR service that caches preprocessed images
type QRService struct {
    cache sync.Map // image_hash -> *gozxing.BinaryBitmap
}

// ProcessQR handles QR detection for uploaded images.
// Without caching: Each request creates a new BinaryBitmap (1.7ms overhead)
// With caching: Reuse BinaryBitmap for duplicate images (12x faster)
func (s *QRService) ProcessQR(imageData []byte) (*QRResult, error) {
    hash := sha256.Sum256(imageData)
    hashStr := hex.EncodeToString(hash[:])
    
    // Check if we've already preprocessed this image
    if cached, ok := s.cache.Load(hashStr); ok {
        // Multiple goroutines may decode the same cached bitmap
        // This is now safe with v0.1.2+
        return s.decodeQR(cached.(*gozxing.BinaryBitmap))
    }
    
    // Preprocess new image (expensive: ~1.8ms for 400x400)
    img, _ := jpeg.Decode(bytes.NewReader(imageData))
    bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
    
    // Cache for future requests
    s.cache.Store(hashStr, bmp)
    
    return s.decodeQR(bmp)
}

func (s *QRService) decodeQR(bmp *gozxing.BinaryBitmap) (*QRResult, error) {
    reader := qrcode.NewQRCodeReader()
    result, err := reader.Decode(bmp, nil) // Safe for concurrent use
    if err != nil {
        return nil, err
    }
    return &QRResult{Text: result.GetText()}, nil
}
```
