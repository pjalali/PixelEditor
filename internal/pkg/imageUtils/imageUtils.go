package imageUtils

import (
	"image"
	"sync"

	"pjalali.github.io/pixeleditor/internal/pkg/colourConversions"
)

func ModifyImageParallel(img *image.RGBA, rOffset, gOffset, bOffset, contrast, hOffset, sOffset, lOffset, nThreads int) {
	x := img.Rect.Max.X
	y := img.Rect.Max.Y
	yChunk := y / nThreads

	var contrastFactor float64 = 1.0
	if contrast != 0 {
		contrastFactor = (259 * (float64(contrast) + 255)) / (255 * (259 - float64(contrast)))
	}

	// // Citation: https://goinbigdata.com/golang-wait-for-all-goroutines-to-finish/
	var wg sync.WaitGroup

	for i := 0; i < nThreads; i++ {
		startY := i * yChunk
		endY := (i+1)*yChunk - 1

		start := startY * img.Stride
		end := endY*img.Stride + x*4

		var slice []uint8
		if i == nThreads-1 {
			slice = img.Pix[start:]

		} else {
			slice = img.Pix[start:end]
		}
		wg.Add(1)

		go modifySlice(&wg, slice, rOffset, gOffset, bOffset, contrastFactor, hOffset, sOffset, lOffset)
	}

	wg.Wait()
}

func modifySlice(wg *sync.WaitGroup, img []uint8, rOffset, gOffset, bOffset int, contrastFactor float64, hOffset, sOffset, lOffset int) {
	defer wg.Done()

	for i := 0; i < len(img)-3; i += 4 {
		if rOffset != 0 {
			img[i] = clamp(int(img[i]) + rOffset)
		}

		if gOffset != 0 {
			img[i+1] = clamp(int(img[i+1]) + gOffset)
		}

		if bOffset != 0 {
			img[i+2] = clamp(int(img[i+2]) + bOffset)
		}

		if contrastFactor != 1.0 {
			oldR := float64(img[i])
			oldG := float64(img[i+1])
			oldB := float64(img[i+2])

			img[i] = clamp(int(contrastFactor*(oldR-128) + 128))
			img[i+1] = clamp(int(contrastFactor*(oldG-128) + 128))
			img[i+2] = clamp(int(contrastFactor*(oldB-128) + 128))
		}

		if hOffset != 0 || sOffset != 0 || lOffset != 0 {
			r := img[i]
			g := img[i+1]
			b := img[i+2]

			rgbPoint := colourConversions.RGBPoint{r, g, b}

			hslPoint := colourConversions.RGBtoHSL(rgbPoint)

			if hOffset != 0 {
				hslPoint.H = float64(hOffset)
			}
			if sOffset != 0 {
				modifyAndClipFloat(&hslPoint.S, float64(sOffset), 0, 100)
			}
			if lOffset != 0 {
				modifyAndClipFloat(&hslPoint.L, float64(lOffset), 0, 100)
			}

			updatedRGBPoint := colourConversions.HSLToRGB(hslPoint)

			img[i] = updatedRGBPoint.R
			img[i+1] = updatedRGBPoint.G
			img[i+2] = updatedRGBPoint.B
		}
	}
}

func modifyAndClipFloat(initialValue *float64, offset, min, max float64) {
	*initialValue += float64(offset)

	if *initialValue > max {
		*initialValue = max
	} else if *initialValue < min {
		*initialValue = min
	}
}

func clamp(value int) uint8 {
	if value > 255 {
		return 255
	} else if value < 0 {
		return 0
	}
	return uint8(value)
}
