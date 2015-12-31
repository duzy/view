// Copyright (c) 2015 Duzy Chan <code@duzy.info>.
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 

package gv

import (
        "reflect"
)

type Value struct {
        reflect.Value
}

type ViewType string
type PropName string
type SignalName string

// A Bag is a property bag for view.
type Bag interface {
        Get(name PropName) (Value, error)
        Set(name PropName, value Value) error
}

// A View is visible rectangle on the screen. Such as top level
// window, or a child view in a window or view.
type View interface {
        Bag

        Connect(name SignalName, h interface{}) error
        Disconnect(name SignalName) (interface{}, error)
}

const (
        ViewTopLevel ViewType   = ""

        ShowAll PropName        = "show-all"
        Size                    = "size"
        Title                   = "title"

        OnDestroy SignalName    = "destroy"
)

func ValueOf(i interface{}) Value {
        return Value{ reflect.ValueOf(i) }
}

// Create a new view.
func NewView(t ViewType) View {
        return View(newGtkView(t))
}

// Run interaction message loop.
func Interact() {
        runGtkMain()
}
