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

func isSortedB(ar [][]byte) int {
	for i := len(ar) - 1; i > 0; i-- {
		if sixb.BtoS(ar[i]) < sixb.BtoS(ar[i-1]) {
			return i
		}
	}
	return 0
}

func insertionB(slc [][]byte) {
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
		if sixb.BtoS(val) < sixb.BtoS(pre) {
			goto loop
		}
		if l == h {
			continue
		}
	last:
		slc[l] = val
	}
}

func pivotB(slc [][]byte, n uint) string {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc - 1]string
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = sixb.BtoS(slc[first])
		first += step
	}
	insertionS(sample[:n])

	return sample[n>>1]
}

func partOneB(slc [][]byte, pv string) int {
	l, h := 0, len(slc)-1
	goto start
second:
	for {
		h--
		if h <= l {
			return l
		}
		if sixb.BtoS(slc[h]) <= pv {
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

	if pv <= sixb.BtoS(slc[h]) {
		if pv < sixb.BtoS(slc[l]) {
			goto second
		}
		goto next
	}
	for {
		if pv <= sixb.BtoS(slc[l]) {
			goto swap
		}
		l++
		if h <= l {
			return l + 1
		}
	}
last:
	if l == h && sixb.BtoS(slc[h]) < pv {
		l++
	}
	return l
}

func partTwoB(slc [][]byte, l, h int, pv string) int {
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
		if sixb.BtoS(slc[h]) <= pv {
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

	if pv <= sixb.BtoS(slc[h]) {
		if pv < sixb.BtoS(slc[l]) {
			goto second
		}
		goto next
	}
	for {
		if pv <= sixb.BtoS(slc[l]) {
			goto swap
		}
		l--
		if l < 0 {
			return h
		}
	}
}

func gPartOneB(ar [][]byte, pv string, ch dedego.Chan[int]) {
	ch.Send(partOneB(ar, pv))
}

func partConB(slc [][]byte, ch dedego.Chan[int]) int {

	pv := pivotB(slc, nsConc-1)
	mid := len(slc) >> 1
	l, h := mid>>1, sixb.MeanI(mid, len(slc))
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{

				gPartOneB(slc[l:h:h], pv, ch)
			}
		}()
	}()

	r := partTwoB(slc, l, h, pv)

	k := l + ch.Receive()

	if r < mid {
		for ; 0 <= r; r-- {
			if pv < sixb.BtoS(slc[r]) {
				k--
				slc[r], slc[k] = slc[k], slc[r]
			}
		}
	} else {
		for ; r < len(slc); r++ {
			if sixb.BtoS(slc[r]) < pv {
				slc[r], slc[k] = slc[k], slc[r]
				k++
			}
		}
	}
	return k
}

func shortB(ar [][]byte) {
start:
	first, step, last := minMaxSample(uint(len(ar)), 3)
	f, pv, l := sixb.BtoS(ar[first]), sixb.BtoS(ar[first+step]), sixb.BtoS(ar[last])

	if pv < f {
		pv, f = f, pv
	}
	if l < pv {
		if l < f {
			pv = f
		} else {
			pv = l
		}
	}

	k := partOneB(ar, pv)
	var aq [][]byte

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:]
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenInsFC {
		shortB(aq)
		goto start
	}
isort:
	insertionB(aq)

	if len(ar) > MaxLenInsFC {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto isort
	}
}

func gLongB(ar [][]byte, sv *syncVar) {
	longB(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 {
		sv.done.Send(0)

	}
}

func longB(ar [][]byte, sv *syncVar) {
start:
	pv := pivotB(ar, nsLong-1)
	k := partOneB(ar, pv)
	var aq [][]byte

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:]
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) <= MaxLenRecFC {

		if len(aq) > MaxLenInsFC {
			shortB(aq)
		} else {
			insertionB(aq)
		}

		if len(ar) > MaxLenRecFC {
			goto start
		}
		shortB(ar)
		return
	}

	if sv == nil || gorFull(sv) {
		longB(aq, sv)
		goto start
	}

	atomic.AddUint64(&sv.nGor, 1)
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{
				gLongB(ar, sv)
			}
		}()
	}()
	ar = aq
	goto start
}

func sortB(ar [][]byte) {

	if len(ar) < 2*(MaxLenRecFC+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRecFC {
			longB(ar, nil)
		} else if len(ar) > MaxLenInsFC {
			shortB(ar)
		} else {
			insertionB(ar)
		}
		return
	}

	sv := syncVar{1,
		dedego.NewChan[int](int(0))}
	for {

		k := partConB(ar, sv.done)
		var aq [][]byte

		if k < len(ar)-k {
			aq = ar[:k:k]
			ar = ar[k:]
		} else {
			aq = ar[k:]
			ar = ar[:k:k]
		}

		if len(aq) > MaxLenRecFC {
			atomic.AddUint64(&sv.nGor, 1)
			func() {
				DedegoRoutineIndex := dedego.SpawnPre()
				go func() {
					dedego.SpawnPost(DedegoRoutineIndex)
					{
						gLongB(aq, &sv)
					}
				}()
			}()

		} else if len(aq) > MaxLenInsFC {
			shortB(aq)
		} else {
			insertionB(aq)
		}

		if len(ar) < 2*(MaxLenRecFC+1) || gorFull(&sv) {
			break
		}

	}

	longB(ar, &sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 {
		sv.done.Receive()

	}
}
