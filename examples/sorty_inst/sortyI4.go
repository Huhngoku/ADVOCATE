/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"sync/atomic"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
	"github.com/jfcg/sixb"
)

func isSortedI4(ar []int32) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

func insertionI4(slc []int32) {
	for h := 1; h < len(slc); h++ {
		l, val := h, slc[h]
		var pre int32
		goto start
	loop:
		slc[l] = pre
		l--
		if l == 0 {
			goto last
		}
	start:
		pre = slc[l-1]
		if val < pre {
			goto loop
		}
		if l == h {
			continue
		}
	last:
		slc[l] = val
	}
}

func pivotI4(slc []int32, n uint) int32 {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc]int32
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = slc[first]
		first += step
	}
	insertionI4(sample[:n])

	n >>= 1
	return sixb.MeanI4(sample[n-1], sample[n])
}

func partOneI4(slc []int32, pv int32) int {
	l, h := 0, len(slc)-1
	goto start
second:
	for {
		h--
		if h <= l {
			return l
		}
		if slc[h] <= pv {
			break
		}
	}
swap:
	slc[l], slc[h] = slc[h], slc[l]
next:
	l++
	h--
start:
	if h <= l {
		goto last
	}

	if pv <= slc[h] {
		if pv < slc[l] {
			goto second
		}
		goto next
	}
	for {
		if pv <= slc[l] {
			goto swap
		}
		l++
		if h <= l {
			return l + 1
		}
	}
last:
	if l == h && slc[h] < pv {
		l++
	}
	return l
}

func partTwoI4(slc []int32, l, h int, pv int32) int {
	l--
	if h <= l {
		return -1
	}
	goto start
second:
	for {
		h++
		if h >= len(slc) {
			return l
		}
		if slc[h] <= pv {
			break
		}
	}
swap:
	slc[l], slc[h] = slc[h], slc[l]
next:
	l--
	h++
start:
	if l < 0 {
		return h
	}
	if h >= len(slc) {
		return l
	}

	if pv <= slc[h] {
		if pv < slc[l] {
			goto second
		}
		goto next
	}
	for {
		if pv <= slc[l] {
			goto swap
		}
		l--
		if l < 0 {
			return h
		}
	}
}

func gPartOneI4(ar []int32, pv int32, ch dedego.Chan[int]) {
	ch.Send(partOneI4(ar, pv))
}

func partConI4(slc []int32, ch dedego.Chan[int]) int {

	pv := pivotI4(slc, nsConc)
	mid := len(slc) >> 1
	l, h := mid>>1, sixb.MeanI(mid, len(slc))
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{

				gPartOneI4(slc[l:h:h], pv, ch)
			}
		}()
	}()

	r := partTwoI4(slc, l, h, pv)

	k := l + ch.Receive()

	if r < mid {
		for ; 0 <= r; r-- {
			if pv < slc[r] {
				k--
				slc[r], slc[k] = slc[k], slc[r]
			}
		}
	} else {
		for ; r < len(slc); r++ {
			if slc[r] < pv {
				slc[r], slc[k] = slc[k], slc[r]
				k++
			}
		}
	}
	return k
}

func shortI4(ar []int32) {
start:
	first, step := minMaxFour(uint32(len(ar)))
	a, b, c, d := ar[first], ar[first+step], ar[first+2*step], ar[first+3*step]

	if d < b {
		d, b = b, d
	}
	if c < a {
		c, a = a, c
	}
	if d < c {
		c = d
	}
	if b < a {
		b = a
	}
	pv := sixb.MeanI4(b, c)

	k := partOneI4(ar, pv)
	var aq []int32

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:]
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortI4(aq)
		goto start
	}
isort:
	insertionI4(aq)

	if len(ar) > MaxLenIns {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto isort
	}
}

func gLongI4(ar []int32, sv *syncVar) {
	longI4(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 {
		sv.done.Send(0)

	}
}

func longI4(ar []int32, sv *syncVar) {
start:
	pv := pivotI4(ar, nsLong)
	k := partOneI4(ar, pv)
	var aq []int32

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:]
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) <= MaxLenRec {

		if len(aq) > MaxLenIns {
			shortI4(aq)
		} else {
			insertionI4(aq)
		}

		if len(ar) > MaxLenRec {
			goto start
		}
		shortI4(ar)
		return
	}

	if sv == nil || gorFull(sv) {
		longI4(aq, sv)
		goto start
	}

	atomic.AddUint64(&sv.nGor, 1)
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{
				gLongI4(ar, sv)
			}
		}()
	}()
	ar = aq
	goto start
}

func sortI4(ar []int32) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec {
			longI4(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortI4(ar)
		} else {
			insertionI4(ar)
		}
		return
	}

	sv := syncVar{1,
		dedego.NewChan[int](int(0))}
	for {

		k := partConI4(ar, sv.done)
		var aq []int32

		if k < len(ar)-k {
			aq = ar[:k:k]
			ar = ar[k:]
		} else {
			aq = ar[k:]
			ar = ar[:k:k]
		}

		if len(aq) > MaxLenRec {
			atomic.AddUint64(&sv.nGor, 1)
			func() {
				DedegoRoutineIndex := dedego.SpawnPre()
				go func() {
					dedego.SpawnPost(DedegoRoutineIndex)
					{
						gLongI4(aq, &sv)
					}
				}()
			}()

		} else if len(aq) > MaxLenIns {
			shortI4(aq)
		} else {
			insertionI4(aq)
		}

		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}

	}

	longI4(ar, &sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 {
		sv.done.Receive()

	}
}
