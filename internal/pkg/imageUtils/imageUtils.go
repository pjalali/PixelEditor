package imageUtils

import (
	"image"
	"sync"

	"pjalali.github.io/pixeleditor/internal/pkg/colourUtils"
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
		if rOffset != 0 || gOffset != 0 || bOffset != 0 || contrastFactor != 1.0 {
			modifyRGBValues(img, i, rOffset, gOffset, bOffset, contrastFactor)
		}

		if hOffset != 0 || sOffset != 0 || lOffset != 0 {
			modifyHSLValues(img, i, hOffset, sOffset, lOffset)
		}
	}
}

func modifyRGBValues(img []uint8, index, rOffset, gOffset, bOffset int, contrastFactor float64) {
	if rOffset != 0 {
		img[index] = clampToUInt8(int(img[index]) + rOffset)
	}

	if gOffset != 0 {
		img[index+1] = clampToUInt8(int(img[index+1]) + gOffset)
	}

	if bOffset != 0 {
		img[index+2] = clampToUInt8(int(img[index+2]) + bOffset)
	}

	if contrastFactor != 1.0 {
		oldR := float64(img[index])
		oldG := float64(img[index+1])
		oldB := float64(img[index+2])

		img[index] = clampToUInt8(int(contrastFactor*(oldR-128) + 128))
		img[index+1] = clampToUInt8(int(contrastFactor*(oldG-128) + 128))
		img[index+2] = clampToUInt8(int(contrastFactor*(oldB-128) + 128))
	}
}

func modifyHSLValues(img []uint8, index, hOffset, sOffset, lOffset int) {
	r := img[index]
	g := img[index+1]
	b := img[index+2]

	rgbPoint := colourUtils.RGBPoint{r, g, b}

	hslPoint := colourUtils.RGBtoHSL(rgbPoint)

	if hOffset != 0 {
		hslPoint.H = float64(hOffset)
	}
	if sOffset != 0 {
		modifyAndClampFloat(&hslPoint.S, float64(sOffset), 0, 100)
	}
	if lOffset != 0 {
		modifyAndClampFloat(&hslPoint.L, float64(lOffset), 0, 100)
	}

	updatedRGBPoint := colourUtils.HSLToRGB(hslPoint)

	img[index] = updatedRGBPoint.R
	img[index+1] = updatedRGBPoint.G
	img[index+2] = updatedRGBPoint.B
}

func modifyAndClampFloat(initialValue *float64, offset, min, max float64) {
	*initialValue += float64(offset)

	if *initialValue > max {
		*initialValue = max
	} else if *initialValue < min {
		*initialValue = min
	}
}

func clampToUInt8(value int) uint8 {
	if value > 255 {
		return 255
	} else if value < 0 {
		return 0
	}
	return uint8(value)
}
