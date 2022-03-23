package colourConversions

import "math"

type HSLPoint struct {
	H, S, L float64
}

type RGBPoint struct {
	R, G, B uint8
}

// Formula from https://www.niwa.nu/2013/05/math-behind-colorspace-conversions-rgb-hsl/
func RGBtoHSL(p RGBPoint) HSLPoint {
	var h, s, l float64

	r := float64(p.R) / 255
	g := float64(p.G) / 255
	b := float64(p.B) / 255

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

	return HSLPoint{h, s, l}
}

// Formula from: https://www.niwa.nu/2013/05/math-behind-colorspace-conversions-rgb-hsl/
func HSLToRGB(p HSLPoint) RGBPoint {
	var h, s, l float64
	h = p.H / 360.0
	s = p.S / 100.0
	l = p.L / 100.0

	if s == 0 {
		v := l * 255
		return RGBPoint{uint8(v + 0.5), uint8(v + 0.5), uint8(v + 0.5)}
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

	return RGBPoint{uint8(r + 0.5), uint8(g + 0.5), uint8(b + 0.5)}
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
