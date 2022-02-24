package imageUtils

import (
	"image"
	"math"
	"sync"
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

type hslpoint struct {
	h, s, l float64
}

type rgbpoint struct {
	r, g, b uint8
}

// Formula from https://www.niwa.nu/2013/05/math-behind-colorspace-conversions-rgb-hsl/
func RGBtoHSL(p rgbpoint) hslpoint {
	var h, s, l float64

	r := float64(p.r) / 255
	g := float64(p.g) / 255
	b := float64(p.b) / 255

	xMax := math.Max(math.Max(r, g), b)
	xMin := math.Min(math.Min(r, g), b)

	c := xMax - xMin

	l = (xMax + xMin) / 2

	if c == 0 {
		s = 0
	} else if l < 0.5 {
		s = (xMax - xMin) / (xMax + xMin)
	} else {
		s = (xMax - xMin) / (2 - xMax - xMin)
	}

	if c == 0 {
		h = 0
	} else if xMax == r {
		h = 60 * ((g - b) / c)
	} else if xMax == g {
		h = 60 * (2 + ((b - r) / c))
	} else if xMax == b {
		h = 60 * (4 + ((r - g) / c))
	}

	if h < 0 {
		h += 360
	}

	s *= 100
	l *= 100

	return hslpoint{h, s, l}
}

// Formula from: https://www.niwa.nu/2013/05/math-behind-colorspace-conversions-rgb-hsl/
func HSLToRGB(p hslpoint) rgbpoint {
	var h, s, l float64
	h = p.h / 360.0
	s = p.s / 100.0
	l = p.l / 100.0

	if s == 0 {
		v := l * 255
		return rgbpoint{uint8(v + 0.5), uint8(v + 0.5), uint8(v + 0.5)}
	}

	var t1 float64
	if l < 0.5 {
		t1 = l * (1 + s)
	} else {
		t1 = l + s - (l * s)
	}

	t2 := 2*l - t1

	tr := h + 1.0/3.0
	tg := h
	tb := h - 1.0/3.0

	var r, g, b float64

	r = hueToRGB(t1, t2, tr)
	g = hueToRGB(t1, t2, tg)
	b = hueToRGB(t1, t2, tb)

	r *= 255
	g *= 255
	b *= 255

	return rgbpoint{uint8(r + 0.5), uint8(g + 0.5), uint8(b + 0.5)}
}

// Citation: https://github.com/gerow/go-color/blob/master/color.go
func hueToRGB(t1, t2, tc float64) float64 {
	if tc < 0 {
		tc++
	}
	if tc > 1 {
		tc--
	}
	switch {
	case 6*tc < 1:
		return t2 + (t1-t2)*6*tc
	case 2*tc < 1:
		return t1
	case 3*tc < 2:
		return t2 + (t1-t2)*(2.0/3.0-tc)*6
	}
	return t2
}

func ImageHSLModifications(img *image.RGBA, hOffset, sOffset, lOffset int) {

	x := img.Rect.Max.X
	y := img.Rect.Max.Y

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			r := img.Pix[j*img.Stride+i*4]
			g := img.Pix[j*img.Stride+i*4+1]
			b := img.Pix[j*img.Stride+i*4+2]

			rgbPoint := rgbpoint{r, g, b}

			hslPoint := RGBtoHSL(rgbPoint)

			if hOffset != 0 {
				modifyHue(&hslPoint, hOffset)
			}
			if sOffset != 0 {
				modifySaturation(&hslPoint, sOffset)
			}
			if lOffset != 0 {
				modifyLight(&hslPoint, lOffset)
			}

			updatedRGBPoint := HSLToRGB(hslPoint)

			img.Pix[j*img.Stride+i*4] = updatedRGBPoint.r
			img.Pix[j*img.Stride+i*4+1] = updatedRGBPoint.g
			img.Pix[j*img.Stride+i*4+2] = updatedRGBPoint.b

		}
	}
}

func modifyHue(p *hslpoint, hOffset int) {
	p.h = float64(hOffset)

	if p.h > 360 {
		p.h = 360
	} else if p.h < 0 {
		p.h = 0
	}
}

func modifySaturation(p *hslpoint, sOffset int) {
	p.s += float64(sOffset)

	if p.s > 100 {
		p.s = 100
	} else if p.s < 0 {
		p.s = 0
	}
}

func modifyLight(p *hslpoint, lOffset int) {
	p.l += float64(lOffset)

	if p.l > 100 {
		p.l = 100
	} else if p.l < 0 {
		p.l = 0
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
