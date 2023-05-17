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

func isSortedS(ar []string) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

func insertionS(slc []string) {
	for h := 1; h < len(slc); h++ {
		l, val := h, slc[h]
		var pre string
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

func pivotS(slc []string, n uint) string {

	first, step, _ := minMaxSample(uint(len(slc)), n)

	var sample [nsConc - 1]string
	for i := int(n - 1); i >= 0; i-- {
		sample[i] = slc[first]
		first += step
	}
	insertionS(sample[:n])

	return sample[n>>1]
}

func partOneS(slc []string, pv string) int {
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

func partTwoS(slc []string, l, h int, pv string) int {
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

func gPartOneS(ar []string, pv string, ch dedego.Chan[int]) {
	ch.Send(partOneS(ar, pv))
}

func partConS(slc []string, ch dedego.Chan[int]) int {

	pv := pivotS(slc, nsConc-1)
	mid := len(slc) >> 1
	l, h := mid>>1, sixb.MeanI(mid, len(slc))
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{

				gPartOneS(slc[l:h:h], pv, ch)
			}
		}()
	}()

	r := partTwoS(slc, l, h, pv)

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

func shortS(ar []string) {
start:
	first, step, last := minMaxSample(uint(len(ar)), 3)
	f, pv, l := ar[first], ar[first+step], ar[last]

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

	k := partOneS(ar, pv)
	var aq []string

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:]
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenInsFC {
		shortS(aq)
		goto start
	}
isort:
	insertionS(aq)

	if len(ar) > MaxLenInsFC {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto isort
	}
}

func gLongS(ar []string, sv *syncVar) {
	longS(ar, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 {
		sv.done.Send(0)

	}
}

func longS(ar []string, sv *syncVar) {
start:
	pv := pivotS(ar, nsLong-1)
	k := partOneS(ar, pv)
	var aq []string

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:]
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) <= MaxLenRecFC {

		if len(aq) > MaxLenInsFC {
			shortS(aq)
		} else {
			insertionS(aq)
		}

		if len(ar) > MaxLenRecFC {
			goto start
		}
		shortS(ar)
		return
	}

	if sv == nil || gorFull(sv) {
		longS(aq, sv)
		goto start
	}

	atomic.AddUint64(&sv.nGor, 1)
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{
				gLongS(ar, sv)
			}
		}()
	}()
	ar = aq
	goto start
}

func sortS(ar []string) {

	if len(ar) < 2*(MaxLenRecFC+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRecFC {
			longS(ar, nil)
		} else if len(ar) > MaxLenInsFC {
			shortS(ar)
		} else {
			insertionS(ar)
		}
		return
	}

	sv := syncVar{1,
		dedego.NewChan[int](int(0))}
	for {

		k := partConS(ar, sv.done)
		var aq []string

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
						gLongS(aq, &sv)
					}
				}()
			}()

		} else if len(aq) > MaxLenInsFC {
			shortS(aq)
		} else {
			insertionS(aq)
		}

		if len(ar) < 2*(MaxLenRecFC+1) || gorFull(&sv) {
			break
		}

	}

	longS(ar, &sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 {
		sv.done.Receive()

	}
}
