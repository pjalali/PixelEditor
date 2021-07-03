package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"pjalali.github.io/gorender/pixel"
)

func main() {
	cliArgs := os.Args[1:]

	var err error
	var cliArgsInt [8]int

	for i := 0; i < 8; i++ {
		cliArgsInt[i], err = strconv.Atoi(cliArgs[i+1])
		if err != nil {
			panic(err)
		}
	}

	rOffset := cliArgsInt[0]
	gOffset := cliArgsInt[1]
	bOffset := cliArgsInt[2]
	contrast := cliArgsInt[3]
	hue := cliArgsInt[4]
	sat := cliArgsInt[5]
	light := cliArgsInt[6]
	threads := cliArgsInt[7]

	test := pixel.ReadImageFromFile(cliArgs[0])

	if threads > test.Rect.Max.Y {
		panic("More threads than rows. Exiting.")
	} else if threads < 1 {
		panic("Need at least one thread. Exiting.")
	}

	startRGBContrast := time.Now()

	if rOffset != 0 || gOffset != 0 || bOffset != 0 {
		pixel.ModifyRGBParallel(test, rOffset, gOffset, bOffset, threads)
	}

	if contrast != 0 {
		pixel.ModifyContrastParallel(test, contrast, threads)
	}

	elapsedRGBContrast := time.Since(startRGBContrast)

	startHSL := time.Now()

	if hue != 0 || sat != 0 || light != 0 {
		pixel.ImageHSLModifications(test, hue, sat, light)
	}

	elapsedHSL := time.Since(startHSL)

	fmt.Println(elapsedRGBContrast, elapsedHSL, elapsedRGBContrast+elapsedHSL)

	pixel.WriteImageToFile("output.png", test)
}
