// Copyright (c) 2015 Duzy Chan <code@duzy.info>.
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 

package gv

import (
        "reflect"
        "image"
)

type ValueType struct { reflect.Value }
type PointType struct { image.Point }
type SizeType  struct { image.Point }

type ViewClass string
type PropName string
type SignalName string
type Connection uint

// A Bag is a property bag for view.
type Bag interface {
        Get(name PropName) (ValueType, error)
        Set(name PropName, value ValueType) error
}

// A View is visible rectangle on the screen. Such as top level
// window, or a child view in a window or view.
type View interface {
        Bag

        Connect(name SignalName, h interface{}) (Connection, error)
        Disconnect(c Connection) (interface{}, error)
}

const (
        ViewTopLevel ViewClass   = ""

        ShowAll PropName        = "show-all"
        Size                    = "size"
        Title                   = "title"

        OnDestroy SignalName    = "destroy"
)

func ValueOf(i interface{}) ValueType {
        return ValueType{ reflect.ValueOf(i) }
}

func NewPoint(x, y int) PointType {
        return PointType{ image.Pt(x, y) }
}

func NewPointValue(x, y int) ValueType {
        return ValueOf(NewPoint(x, y))
}

func NewSize(w, h int) SizeType {
        return SizeType{ image.Pt(w, h) }
}

func NewSizeValue(w, h int) ValueType {
        return ValueOf(NewSize(w, h))
}

// Create a new view.
func NewView(t ViewClass) View {
        return View(newGtkView(t))
}

// Run interaction message loop.
func Interact() {
        runGtkMain()
}

func Quit() {
        quitGtkMain()
}
