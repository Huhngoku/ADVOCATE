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

func IsSortedSlice(ar any) int {
	slc, kind := extractSK(ar)
	switch kind {
	case reflect.Int32:
		i := *(*[]int32)(unsafe.Pointer(&slc))
		return isSortedI4(i)
	case reflect.Int64:
		i := *(*[]int64)(unsafe.Pointer(&slc))
		return isSortedI8(i)
	case reflect.Uint32:
		u := *(*[]uint32)(unsafe.Pointer(&slc))
		return isSortedU4(u)
	case reflect.Uint64:
		u := *(*[]uint64)(unsafe.Pointer(&slc))
		return isSortedU8(u)
	case reflect.Float32:
		f := *(*[]float32)(unsafe.Pointer(&slc))
		return isSortedF4(f)
	case reflect.Float64:
		f := *(*[]float64)(unsafe.Pointer(&slc))
		return isSortedF8(f)
	case sliceBias + reflect.Uint8:
		b := *(*[][]byte)(unsafe.Pointer(&slc))
		return isSortedB(b)
	case reflect.String:
		s := *(*[]string)(unsafe.Pointer(&slc))
		return isSortedS(s)
	}
	panic("sorty: IsSortedSlice: invalid input type")
}

func SortSlice(ar any) {
	slc, kind := extractSK(ar)
	switch kind {
	case reflect.Int32:
		i := *(*[]int32)(unsafe.Pointer(&slc))
		sortI4(i)
	case reflect.Int64:
		i := *(*[]int64)(unsafe.Pointer(&slc))
		sortI8(i)
	case reflect.Uint32:
		u := *(*[]uint32)(unsafe.Pointer(&slc))
		sortU4(u)
	case reflect.Uint64:
		u := *(*[]uint64)(unsafe.Pointer(&slc))
		sortU8(u)
	case reflect.Float32:
		f := *(*[]float32)(unsafe.Pointer(&slc))
		sortF4(f)
	case reflect.Float64:
		f := *(*[]float64)(unsafe.Pointer(&slc))
		sortF8(f)
	case sliceBias + reflect.Uint8:
		b := *(*[][]byte)(unsafe.Pointer(&slc))
		sortB(b)
	case reflect.String:
		s := *(*[]string)(unsafe.Pointer(&slc))
		sortS(s)
	default:
		panic("sorty: SortSlice: invalid input type")
	}
}
