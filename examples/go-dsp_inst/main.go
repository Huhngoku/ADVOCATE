/*
 * Copyright (c) 2011 Matt Jibson <matt.jibson@gmail.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
	"github.com/mjibson/go-dsp/dsputils"
)

const (
	sqrt2_2 = math.Sqrt2 / 2
)

type fftTest struct {
	in  []float64
	out []complex128
}

var fftTests = []fftTest{

	{
		[]float64{1},
		[]complex128{complex(1, 0)},
	},
	{
		[]float64{1, 0},
		[]complex128{complex(1, 0), complex(1, 0)},
	},
	{
		[]float64{1, 0, 0, 0},
		[]complex128{complex(1, 0), complex(1, 0), complex(1, 0), complex(1, 0)},
	},
	{
		[]float64{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0)},
	},

	{
		[]float64{0, 1},
		[]complex128{complex(1, 0), complex(-1, 0)},
	},
	{
		[]float64{0, 1, 0, 0},
		[]complex128{complex(1, 0), complex(0, -1), complex(-1, 0), complex(0, 1)},
	},
	{
		[]float64{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{
			complex(1, 0),
			complex(sqrt2_2, -sqrt2_2),
			complex(0, -1),
			complex(-sqrt2_2, -sqrt2_2),
			complex(-1, 0),
			complex(-sqrt2_2, sqrt2_2),
			complex(0, 1),
			complex(sqrt2_2, sqrt2_2)},
	},

	{
		[]float64{1, 2, 3, 4},
		[]complex128{
			complex(10, 0),
			complex(-2, 2),
			complex(-2, 0),
			complex(-2, -2)},
	},
	{
		[]float64{1, 3, 5, 7},
		[]complex128{
			complex(16, 0),
			complex(-4, 4),
			complex(-4, 0),
			complex(-4, -4)},
	},
	{
		[]float64{1, 2, 3, 4, 5, 6, 7, 8},
		[]complex128{
			complex(36, 0),
			complex(-4, 9.65685425),
			complex(-4, 4),
			complex(-4, 1.65685425),
			complex(-4, 0),
			complex(-4, -1.65685425),
			complex(-4, -4),
			complex(-4, -9.65685425)},
	},

	{
		[]float64{1, 0, 0, 0, 0},
		[]complex128{
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0)},
	},
	{
		[]float64{1, 2, 3},
		[]complex128{
			complex(6, 0),
			complex(-1.5, 0.8660254),
			complex(-1.5, -0.8660254)},
	},
	{
		[]float64{1, 1, 1},
		[]complex128{
			complex(3, 0),
			complex(0, 0),
			complex(0, 0)},
	},
}

type fft2Test struct {
	in  [][]float64
	out [][]complex128
}

var fft2Tests = []fft2Test{
	{
		[][]float64{{1, 2, 3}, {3, 4, 5}},
		[][]complex128{
			{complex(18, 0), complex(-3, 1.73205081), complex(-3, -1.73205081)},
			{complex(-6, 0), complex(0, 0), complex(0, 0)}},
	},
	{
		[][]float64{{0.1, 0.2, 0.3, 0.4, 0.5}, {1, 2, 3, 4, 5}, {3, 2, 1, 0, -1}},
		[][]complex128{
			{complex(21.5, 0), complex(-0.25, 0.34409548), complex(-0.25, 0.08122992), complex(-0.25, -0.08122992), complex(-0.25, -0.34409548)},
			{complex(-8.5, -8.66025404), complex(5.70990854, 4.6742225), complex(1.15694356, 4.41135694), complex(-1.65694356, 4.24889709), complex(-6.20990854, 3.98603154)},
			{complex(-8.5, 8.66025404), complex(-6.20990854, -3.98603154), complex(-1.65694356, -4.24889709), complex(1.15694356, -4.41135694), complex(5.70990854, -4.6742225)}},
	},
}

type fftnTest struct {
	in  []float64
	dim []int
	out []complex128
}

var fftnTests = []fftnTest{
	{
		[]float64{4, 2, 3, 8, 5, 6, 7, 2, 13, 24, 13, 17},
		[]int{2, 2, 3},
		[]complex128{
			complex(104, 0), complex(12.5, 14.72243186), complex(12.5, -14.72243186),
			complex(-42, 0), complex(-10.5, 6.06217783), complex(-10.5, -6.06217783),

			complex(-48, 0), complex(-4.5, -11.25833025), complex(-4.5, 11.25833025),
			complex(22, 0), complex(8.5, -6.06217783), complex(8.5, 6.06217783)},
	},
}

type reverseBitsTest struct {
	in  uint
	sz  uint
	out uint
}

var reverseBitsTests = []reverseBitsTest{
	{0, 1, 0},
	{1, 2, 2},
	{1, 4, 8},
	{2, 4, 4},
	{3, 4, 12},
}

func main() {
	dedego.Init(20)
	defer dedego.RunAnalyzer()
	defer time.Sleep(time.Millisecond)
	numSamples := 8

	x := func(n int) float64 {
		wave0 := math.Sin(2.0 * math.Pi * float64(n) / 8.0)
		wave1 := 0.5 * math.Sin(2*math.Pi*float64(n)/4.0+3.0*math.Pi/4.0)
		return wave0 + wave1
	}

	a := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		a[i] = x(i)
	}

	X := FFTReal(a)

	for i := 0; i < numSamples; i++ {
		r, θ := cmplx.Polar(X[i])
		θ *= 360.0 / (2 * math.Pi)
		if dsputils.Float64Equal(r, 0) {
			θ = 0
		}
		fmt.Printf("X(%d) = %.1f ∠ %.1f°\n", i, r, θ)
	}

}

var dedegoFetchOrder = make(map[int]int)
