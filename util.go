// Copyright (c) 2015 Duzy Chan <code@duzy.info>.
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 

package gv

import (
        "unsafe"
        "reflect"
)

func setSliceHeader(slicePtr, dataPtr unsafe.Pointer, dataLen int) {
	h := (*reflect.SliceHeader)(slicePtr)
	h.Cap = dataLen
	h.Len = dataLen
	h.Data = uintptr(dataPtr)
}
