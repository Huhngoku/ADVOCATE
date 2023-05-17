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

type Lesswap func(i, k, r, s int) bool

func IsSorted(n int, lsw Lesswap) int {
	for i := n - 1; i > 0; i-- {
		if lsw(i, i-1, i, i) {
			return i
		}
	}
	return 0
}

func insertion(lsw Lesswap, lo, hi int) {
	for h := lo + 1; h <= hi; h++ {
		for l := h; lsw(l, l-1, l, l-1); {
			l--
			if l <= lo {
				break
			}
		}
	}
}

func pivot(lsw Lesswap, lo, hi int, n uint) int {

	f, s, l := minMaxSample(uint(hi+1-lo), n)
	first := lo + int(f)
	step := int(s)
	last := lo + int(l)

	for h := first + step; h <= last; h += step {
		for l := h; lsw(l, l-step, l, l-step); {
			l -= step
			if l <= first {
				break
			}
		}
	}

	lsw(first, lo, first, lo)
	lsw(hi, last, hi, last)

	return sixb.MeanI(first, last)
}

func partOne(lsw Lesswap, l, pv, h int) int {

	for ; l < h; l, h = l+1, h-1 {

		if lsw(h, pv, h, h) {
			for {
				if lsw(pv, l, h, l) {
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if lsw(pv, l, l, l) {
			for {
				h--
				if l >= h {
					return l
				}
				if lsw(h, pv, h, l) {
					break
				}
			}
		}
	}

	if l == h && h != pv && lsw(h, pv, h, h) {
		l++
	}
	return l
}

func partTwo(lsw Lesswap, lo, l, pv, h, hi int) int {

	for {
		if lsw(h, pv, h, h) {
			for {
				if lsw(pv, l, h, l) {
					break
				}
				l--
				if l < lo {
					return h
				}
			}
		} else if lsw(pv, l, l, l) {
			for {
				h++
				if h > hi {
					return l
				}
				if lsw(h, pv, h, l) {
					break
				}
			}
		}
		l--
		h++
		if l < lo {
			return h
		}
		if h > hi {
			return l
		}
	}
}

func gPartOne(lsw Lesswap, l, pv, h int, ch dedego.Chan[int]) {
	ch.Send(partOne(lsw, l, pv, h))
}

func partCon(lsw Lesswap, lo, hi int, ch dedego.Chan[int]) int {

	pv := pivot(lsw, lo, hi, nsConc-1)
	lo++
	hi--
	l, h := sixb.MeanI(lo, pv), sixb.MeanI(pv, hi)
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{

				gPartOne(lsw, l+1, pv, h-1, ch)
			}
		}()
	}()

	r := partTwo(lsw, lo, l, pv, h, hi)
	k := ch.Receive()

	if r < pv {
		for ; lo <= r; r-- {
			if lsw(pv, r, k-1, r) {
				k--
				if k == pv {
					pv = r
				}
			}
		}
	} else {
		for ; r <= hi; r++ {
			if lsw(r, pv, r, k) {
				if k == pv {
					pv = r
				}
				k++
			}
		}
	}
	return k
}

func short(lsw Lesswap, lo, hi int) {
start:
	fr, step, _ := minMaxSample(uint(hi+1-lo), 3)
	first := lo + int(fr)
	pv := first + int(step)
	last := pv + int(step)

	lsw(pv, first, pv, first)
	if lsw(last, pv, last, pv) {
		lsw(pv, first, pv, first)
	}

	lsw(first, lo, first, lo)
	lsw(hi, last, hi, last)

	l := partOne(lsw, lo+1, pv, hi-1)
	h := l - 1
	no, n := h-lo, hi-l

	if no < n {
		n, no = no, n
		l, lo = lo, l
	} else {
		h, hi = hi, h
	}

	if n >= MaxLenInsFC {
		short(lsw, l, h)
		goto start
	}

isort:
	for k := l + 1; k <= h; k++ {
		for i := k; lsw(i, i-1, i, i-1); {
			i--
			if i <= l {
				break
			}
		}
	}

	if no >= MaxLenInsFC {
		goto start
	}
	if lo != l {
		l, h = lo, hi
		goto isort
	}
}

func gLong(lsw Lesswap, lo, hi int, sv *syncVar) {
	long(lsw, lo, hi, sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) == 0 {
		sv.done.Send(0)

	}
}

func long(lsw Lesswap, lo, hi int, sv *syncVar) {
start:
	pv := pivot(lsw, lo, hi, nsLong-1)
	l := partOne(lsw, lo+1, pv, hi-1)
	h := l - 1
	no, n := h-lo, hi-l

	if no < n {
		n, no = no, n
		l, lo = lo, l
	} else {
		h, hi = hi, h
	}

	if n < MaxLenRecFC {

		if n >= MaxLenInsFC {
			short(lsw, l, h)
		} else {
			insertion(lsw, l, h)
		}

		if no >= MaxLenRecFC {
			goto start
		}
		short(lsw, lo, hi)
		return
	}

	if sv == nil || gorFull(sv) {
		long(lsw, l, h, sv)
		goto start
	}

	atomic.AddUint64(&sv.nGor, 1)
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{
				gLong(lsw, lo, hi, sv)
			}
		}()
	}()
	lo, hi = l, h
	goto start
}

func Sort(n int, lsw Lesswap) {

	n--
	if n <= 2*MaxLenRecFC || MaxGor <= 1 {

		if n >= MaxLenRecFC {
			long(lsw, 0, n, nil)
		} else if n >= MaxLenInsFC {
			short(lsw, 0, n)
		} else if n > 0 {
			insertion(lsw, 0, n)
		}
		return
	}

	sv := syncVar{1,
		dedego.NewChan[int](int(0))}
	lo, hi := 0, n
	for {

		l := partCon(lsw, lo, hi, sv.done)
		h := l - 1
		no, n := h-lo, hi-l

		if no < n {
			n, no = no, n
			l, lo = lo, l
		} else {
			h, hi = hi, h
		}

		if n >= MaxLenRecFC {
			atomic.AddUint64(&sv.nGor, 1)
			func() {
				DedegoRoutineIndex := dedego.SpawnPre()
				go func() {
					dedego.SpawnPost(DedegoRoutineIndex)
					{
						gLong(lsw, l, h, &sv)
					}
				}()
			}()

		} else if n >= MaxLenInsFC {
			short(lsw, l, h)
		} else {
			insertion(lsw, l, h)
		}

		if no <= 2*MaxLenRecFC || gorFull(&sv) {
			break
		}

	}

	long(lsw, lo, hi, &sv)

	if atomic.AddUint64(&sv.nGor, ^uint64(0)) != 0 {
		sv.done.Receive()

	}
}
