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

func isSortedLenB(ar [][]byte) int {
	for i := len(ar) - 1; i > 0; i-- {
		if len(ar[i]) < len(ar[i-1]) {
			return i
		}
	}
	return 0
}

func insertionLenB(slc [][]byte) {
	for h := 1; h < len(slc); h++ {
		l, val := h, slc[h]
		var pre []byte
		goto start
	loop:
		slc[l] = pre
		l--
		if l == 0 {
			goto last
		}
	start:
		pre = slc[l-1]
		if len(val) < len(pre) {
			goto loop
		}
		if l == h {
			continue
		}
	last:
		slc[l] = val
	}
}

func pivotLenB(slc [][]byte, n uint) int {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc]int
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = len(slc[first])
		first += step
	}
	insertionI(sample[:n])

	n >>= 1
	return sixb.MeanI(sample[n-1], sample[n])
}

func partOneLenB(slc [][]byte, pv int) int {
	l, h := 0, len(slc)-1
	goto start
second:
	for {
		h--
		if h <= l {
			return l
		}
		if len(slc[h]) <= pv {
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

	if pv <= len(slc[h]) {
		if pv < len(slc[l]) {
			goto second
		}
		goto next
	}
	for {
		if pv <= len(slc[l]) {
			goto swap
		}
		l++
		if h <= l {
			return l + 1
		}
	}
last:
	if l == h && len(slc[h]) < pv {
		l++
	}
	return l
}

func partTwoLenB(slc [][]byte, l, h int, pv int) int {
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
		if len(slc[h]) <= pv {
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

	if pv <= len(slc[h]) {
		if pv < len(slc[l]) {
			goto second
		}
		goto next
	}
	for {
		if pv <= len(slc[l]) {
			goto swap
		}
		l--
		if l < 0 {
			return h
		}
	}
}

func gPartOneLenB(ar [][]byte, pv int, ch dedego.Chan[int]) {
	ch.Send(partOneLenB(ar, pv))
}

func partConLenB(slc [][]byte, ch dedego.Chan[int]) int {

	pv := pivotLenB(slc, nsConc)
	mid := len(slc) >> 1
	l, h := mid>>1, sixb.MeanI(mid, len(slc))
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{

				gPartOneLenB(slc[l:h:h], pv, ch)
			}
		}()
	}()

	r := partTwoLenB(slc, l, h, pv)

	k := l + ch.Receive()

	if r < mid {
		for ; 0 <= r; r-- {
			if pv < len(slc[r]) {
				k--
				slc[r], slc[k] = slc[k], slc[r]
			}
		}
	} else {
		for ; r < len(slc); r++ {
			if len(slc[r]) < pv {
				slc[r], slc[k] = slc[k], slc[r]
				k++
			}
		}
	}
	return k
}

func shortLenB(ar [][]byte) {
start:
	first, step := minMaxFour(uint32(len(ar)))
	a, b := len(ar[first]), len(ar[first+step])
	c, d := len(ar[first+2*step]), len(ar[first+3*step])

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
	pv := sixb.MeanI(b, c)

	k := partOneLenB(ar, pv)
	var aq [][]byte

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:]
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortLenB(aq)
		goto start
	}
isort:
	insertionLenB(aq)

	if len(ar) > MaxLenIns {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto isort
	}
}

func gLongLenB(ar [][]byte, sv *syncVar) {
	longLenB(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 {
		sv.done.Send(0)

	}
}

func longLenB(ar [][]byte, sv *syncVar) {
start:
	pv := pivotLenB(ar, nsLong)
	k := partOneLenB(ar, pv)
	var aq [][]byte

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:]
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) <= MaxLenRec {

		if len(aq) > MaxLenIns {
			shortLenB(aq)
		} else {
			insertionLenB(aq)
		}

		if len(ar) > MaxLenRec {
			goto start
		}
		shortLenB(ar)
		return
	}

	if sv == nil || gorFull(sv) {
		longLenB(aq, sv)
		goto start
	}

	atomic.AddUint64(&sv.nGor, 1)
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{
				gLongLenB(ar, sv)
			}
		}()
	}()
	ar = aq
	goto start
}

func sortLenB(ar [][]byte) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec {
			longLenB(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortLenB(ar)
		} else {
			insertionLenB(ar)
		}
		return
	}

	sv := syncVar{1,
		dedego.NewChan[int](int(0))}
	for {

		k := partConLenB(ar, sv.done)
		var aq [][]byte

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
						gLongLenB(aq, &sv)
					}
				}()
			}()

		} else if len(aq) > MaxLenIns {
			shortLenB(aq)
		} else {
			insertionLenB(aq)
		}

		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}

	}

	longLenB(ar, &sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 {
		sv.done.Receive()

	}
}
