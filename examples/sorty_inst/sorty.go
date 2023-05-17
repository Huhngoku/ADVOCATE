/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"reflect"
	"unsafe"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
	"github.com/jfcg/sixb"
)

var MaxGor uint64 = 3

func init() {
	if !(4097 > MaxGor && MaxGor > 0 && MaxLenRec > MaxLenRecFC && MaxLenRecFC >
		2*MaxLenIns && MaxLenIns > MaxLenInsFC && MaxLenInsFC > 2*nsShort) {
		panic("sorty: check your MaxGor/MaxLen* values")
	}
}

type FloatOption int32

const (
	NaNsmall FloatOption = iota - 1
	NaNignore
	NaNlarge
)

var NaNoption = NaNlarge

func Search(n int, fn func(int) bool) int {
	l, h := 0, n

	for l < h {
		m := sixb.MeanI(l, h)

		if fn(m) {
			h = m
		} else {
			l = m + 1
		}
	}
	return l
}

type syncVar struct {
	nGor uint64
	done dedego.Chan[int]
}

func gorFull(sv *syncVar) bool {
	mg := MaxGor
	return sv.nGor >= mg
}

const (
	nsShort = 4
	nsLong  = 6
	nsConc  = 8
)

func minMaxSample(slen, n uint) (first, step, last uint) {
	step = slen / n
	n--
	span := n * step
	tail := slen - span
	if tail > n && tail>>1 > (step+1)>>1 {
		step++
		span += n
		tail -= n
	}
	first = tail >> 1
	last = first + span
	return
}

var firstFour = [8]uint32{0, 0, ^uint32(0), 0, 0, 1, 1, 0}
var stepFour = [8]uint32{0, 0, 1, 1, 0, 0, 0, 1}

func minMaxFour(slen uint32) (first, step uint32) {
	mod := slen & 7
	first = slen>>3 + firstFour[mod]
	step = slen>>2 + stepFour[mod]
	return
}

func insertionI(slc []int) {
	if unsafe.Sizeof(int(0)) == 8 {
		insertionI8(*(*[]int64)(unsafe.Pointer(&slc)))
	} else {
		insertionI4(*(*[]int32)(unsafe.Pointer(&slc)))
	}
}

const sliceBias reflect.Kind = 100

func extractSK(ar any) (slc sixb.Slice, kind reflect.Kind) {
	tipe := reflect.TypeOf(ar)
	if tipe.Kind() != reflect.Slice {
		return
	}
	tipe = tipe.Elem()
	kind = tipe.Kind()

	switch kind {

	case reflect.Uintptr, reflect.Pointer, reflect.UnsafePointer:
		kind = reflect.Uint32 + reflect.Kind(unsafe.Sizeof(uintptr(0))>>3)
	case reflect.Uint:
		kind = reflect.Uint32 + reflect.Kind(unsafe.Sizeof(uint(0))>>3)
	case reflect.Int:
		kind = reflect.Int32 + reflect.Kind(unsafe.Sizeof(int(0))>>3)

	case reflect.Slice:
		kind = sliceBias + tipe.Elem().Kind()

	case reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
	default:
		kind = reflect.Invalid
		return
	}

	v := reflect.ValueOf(ar)
	p, l := v.Pointer(), v.Len()
	slc = sixb.Slice{Data: unsafe.Pointer(p), Len: l, Cap: l}
	return
}
