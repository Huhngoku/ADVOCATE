/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"reflect"
	"unsafe"
)

func IsSortedLen(ar any) int {
	slc, kind := extractSK(ar)
	switch {
	case kind == reflect.String:
		s := *(*[]string)(unsafe.Pointer(&slc))
		return isSortedLenS(s)
	case kind >= sliceBias:
		b := *(*[][]byte)(unsafe.Pointer(&slc))
		return isSortedLenB(b)
	}
	panic("sorty: IsSortedLen: invalid input type")
}

func SortLen(ar any) {
	slc, kind := extractSK(ar)
	switch {
	case kind == reflect.String:
		s := *(*[]string)(unsafe.Pointer(&slc))
		sortLenS(s)
	case kind >= sliceBias:
		b := *(*[][]byte)(unsafe.Pointer(&slc))
		sortLenB(b)
	default:
		panic("sorty: SortLen: invalid input type")
	}
}
