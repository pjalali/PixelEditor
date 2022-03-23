package imageUtils

import (
	"image"
	"sync"

	"pjalali.github.io/pixeleditor/internal/pkg/colourConversions"
)

func ModifyRGBParallel(img *image.RGBA, rOffset, gOffset, bOffset, nThreads int) {
	x := img.Rect.Max.X
	y := img.Rect.Max.Y
	yChunk := y / nThreads

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

		go modifyRGB(&wg, slice, rOffset, gOffset, bOffset)
	}

	wg.Wait()
}

func modifyRGB(wg *sync.WaitGroup, img []uint8, rOffset, gOffset, bOffset int) {
	defer wg.Done()

	for i := 0; i < len(img)-3; i += 4 {
		var oldR, oldG, oldB int
		oldR = int(img[i])
		oldG = int(img[i+1])
		oldB = int(img[i+2])

		img[i] = clamp(oldR + rOffset)
		img[i+1] = clamp(oldG + gOffset)
		img[i+2] = clamp(oldB + bOffset)
	}
}

func ModifyContrastParallel(img *image.RGBA, contrast, nThreads int) {

	x := img.Rect.Max.X
	y := img.Rect.Max.Y
	yChunk := y / nThreads
	contrastFactor := (259 * (float32(contrast) + 255)) / (255 * (259 - float32(contrast)))

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

		go modifyContrast(&wg, slice, contrastFactor)
	}

	wg.Wait()
}

func modifyContrast(wg *sync.WaitGroup, img []uint8, factor float32) {

	defer wg.Done()

	for i := 0; i < len(img)-3; i += 4 {
		var oldR, oldG, oldB float32
		oldR = float32(img[i])
		oldG = float32(img[i+1])
		oldB = float32(img[i+2])

		img[i] = clamp(int(factor*(oldR-128) + 128))
		img[i+1] = clamp(int(factor*(oldG-128) + 128))
		img[i+2] = clamp(int(factor*(oldB-128) + 128))
	}
}

func ImageHSLModifications(img *image.RGBA, hOffset, sOffset, lOffset int) {

	x := img.Rect.Max.X
	y := img.Rect.Max.Y

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			r := img.Pix[j*img.Stride+i*4]
			g := img.Pix[j*img.Stride+i*4+1]
			b := img.Pix[j*img.Stride+i*4+2]

			rgbPoint := colourConversions.RGBPoint{r, g, b}

			hslPoint := colourConversions.RGBtoHSL(rgbPoint)

			if hOffset != 0 {
				modifyHue(&hslPoint, hOffset)
			}
			if sOffset != 0 {
				modifySaturation(&hslPoint, sOffset)
			}
			if lOffset != 0 {
				modifyLight(&hslPoint, lOffset)
			}

			updatedRGBPoint := colourConversions.HSLToRGB(hslPoint)

			img.Pix[j*img.Stride+i*4] = updatedRGBPoint.R
			img.Pix[j*img.Stride+i*4+1] = updatedRGBPoint.G
			img.Pix[j*img.Stride+i*4+2] = updatedRGBPoint.B

		}
	}
}

func modifyHue(p *colourConversions.HSLPoint, hOffset int) {
	p.H = float64(hOffset)

	if p.H > 360 {
		p.H = 360
	} else if p.H < 0 {
		p.H = 0
	}
}

func modifySaturation(p *colourConversions.HSLPoint, sOffset int) {
	p.S += float64(sOffset)

	if p.S > 100 {
		p.S = 100
	} else if p.S < 0 {
		p.S = 0
	}
}

func modifyLight(p *colourConversions.HSLPoint, lOffset int) {
	p.L += float64(lOffset)

	if p.L > 100 {
		p.L = 100
	} else if p.L < 0 {
		p.L = 0
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
