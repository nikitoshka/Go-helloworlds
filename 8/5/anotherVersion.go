package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"strconv"
	"time"
)

func main() {
	now := time.Now()

	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 4096, 4096
	)

	var gocount int
	if len(os.Args) == 1 {
		gocount = 4
	} else {
		if c, err := strconv.Atoi(os.Args[1]); err == nil {
			gocount = c
		} else {
			gocount = 4
		}
	}

	ch := make(chan bool)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < gocount; i++ {
		go func(i int) {
			var upbound int
			if i != gocount-1 {
				upbound = (i + 1) * height / gocount
			} else {
				upbound = height
			}
			for py := i * height / gocount; py < upbound; py++ {
				y := float64(py)/height*(ymax-ymin) + ymin
				for px := 0; px < width; px++ {
					x := float64(px)/width*(xmax-xmin) + xmin
					z := complex(x, y)
					// Image point (px, py) represents complex value z.
					img.Set(px, py, mandelbrot(z))
				}
			}
			ch <- true
		}(i)
	}

	for i := 0; i < gocount; i++ {
		<-ch
	}

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

// Some other interesting functions:

func acos(z complex128) color.Color {
	v := cmplx.Acos(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}

func sqrt(z complex128) color.Color {
	v := cmplx.Sqrt(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{128, blue, red}
}

// f(x) = x^4 - 1
//
// z' = z - f(z)/f'(z)
//    = z - (z^4 - 1) / (4 * z^3)
//    = z - (z - 1/z^3) / 4
func newton(z complex128) color.Color {
	const iterations = 37
	const contrast = 7
	for i := uint8(0); i < iterations; i++ {
		z -= (z - 1/(z*z*z)) / 4
		if cmplx.Abs(z*z*z*z-1) < 1e-6 {
			return color.Gray{255 - contrast*i}
		}
	}
	return color.Black
}
