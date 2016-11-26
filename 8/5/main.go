// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 61.
//!+

// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
	"time"
)

type pointDescr struct {
	x, y int
	col  color.Color
}

func main() {
	now := time.Now()

	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 2048, 2048
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	ch := make(chan *pointDescr)
	var wg sync.WaitGroup

	go addPoint(img, ch)

	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin

		wg.Add(1)
		go func(py int) {
			defer func() {
				wg.Done()
			}()

			for px := 0; px < width; px++ {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)

				ch <- &pointDescr{px, py, mandelbrot(z)}
			}
		}(py)
	}

	wg.Wait()
	close(ch)

	png.Encode(os.Stdout, img) // NOTE: ignoring errors

	fmt.Fprintf(os.Stderr, "time elapsed: %s\n", time.Since(now))
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}

func addPoint(img *image.RGBA, ch <-chan *pointDescr) {
	if img == nil {
		return
	}

	for point := range ch {
		if point == nil {
			continue
		}

		img.Set(point.x, point.y, point.col)
	}
}
