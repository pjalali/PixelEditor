package imageUtils

import (
	"image"
	"sync"

	"pjalali.github.io/pixeleditor/internal/pkg/colourConversions"
)

func ModifyRGBParallel(img *image.RGBA, rOffset, gOffset, bOffset, contrast, nThreads int) {
	x := img.Rect.Max.X
	y := img.Rect.Max.Y
	yChunk := y / nThreads

	var contrastFactor float32 = 1.0
	if contrast != 0 {
		contrastFactor = (259 * (float32(contrast) + 255)) / (255 * (259 - float32(contrast)))
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

		go modifySlice(&wg, slice, rOffset, gOffset, bOffset, contrastFactor)
	}

	wg.Wait()
}

func modifySlice(wg *sync.WaitGroup, img []uint8, rOffset, gOffset, bOffset int, contrastFactor float32) {
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

		if contrastFactor != 1 {
			oldR := float32(img[i])
			oldG := float32(img[i+1])
			oldB := float32(img[i+2])

			img[i] = clamp(int(contrastFactor*(oldR-128) + 128))
			img[i+1] = clamp(int(contrastFactor*(oldG-128) + 128))
			img[i+2] = clamp(int(contrastFactor*(oldB-128) + 128))

		}
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
