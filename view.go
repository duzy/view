// Copyright (c) 2015 Duzy Chan <code@duzy.info>.
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 

package gv

import (
        "image"
        //"fmt"
)

type PointType struct { image.Point }
type SizeType  struct { image.Point }

type PropName string
type SignalName string
type Connection uint

// A Bag is a property bag for view.
type getter interface { Get(name PropName) (interface{}, error) }
type setter interface { Set(name PropName, value interface{}) error }
type Bag interface { getter; setter }

type adder interface { Add(v View) error }
type finder interface { Find(id string) View }
type Container interface { adder; finder }

// A View is visible rectangle on the screen. Such as top level
// window, or a child view in a window or view.
type View interface {
        Bag

        Connect(name SignalName, h interface{}) (Connection, error)
        Disconnect(c Connection) (interface{}, error)
}

const (
        Id                      = "id"
        Show                    = "show"
        Size                    = "size"
        Text                    = "text"
        Pack                    = "pack"
        Spacing                 = "spacing"
        Padding                 = "padding"
        Expend                  = "expend"
        Fill                    = "fill"

        OnAccelClosuresChanged SignalName = "accel-closures-changed"      // void
        OnButtonPressEvent      = "button-press-event"          // boolean
        OnButtonReleaseEvent    = "button-release-event"        // boolean
        OnCanActivateAccel      = "can-activate-accel"          // boolean
        OnChildNotify           = "child-notify"                // void
        OnCompositedChanged     = "composited-changed"          // void
        OnConfigureEvent        = "configure-event"             // boolean
        OnDamageEvent           = "damage-event"                // boolean
        OnDeleteEvent           = "delete-event"                // boolean
        OnDestroyEvent          = "destroy-event"               // boolean
        OnDestroy               = "destroy"                     // void
        OnDirectionChanged      = "direction-changed"           // void
        OnDragBegin             = "drag-begin"                  // void
        OnDragDataDelete        = "drag-data-delete"            // void
        OnDragDataGet           = "drag-data-get"               // void
        OnDragDataReceived      = "drag-data-received"          // void
        OnDragDrop              = "drag-drop"                   // boolean
        OnDragEnd               = "drag-end"                    // void
        OnDragFailed            = "drag-failed"                 // boolean
        OnDragLeave             = "drag-leave"                  // void
        OnDragMotion            = "drag-motion"                 // boolean
        OnDraw                  = "draw"                        // boolean
        OnEnterNotifyEvent      = "enter-notify-event"          // boolean
        OnEvent                 = "event"                       // boolean
        OnEventAfter            = "event-after"                 // void
        OnFocus                 = "focus"                       // boolean
        OnFocusInEvent          = "focus-in-event"              // boolean
        OnFocusOutEvent         = "focus-out-event"             // boolean
        OnGrabBrokenEvent       = "grab-broken-event"           // boolean
        OnGrabFocus             = "grab-focus"                  // void
        OnGrabNotify            = "grab-notify"                 // void
        OnHide                  = "hide"                        // void
        OnHierarchyChanged      = "hierarchy-changed"           // void
        OnKeyPressEvent         = "key-press-event"             // boolean
        OnKeyReleaseEvent       = "key-release-event"           // boolean
        OnKeynavFailed          = "keynav-failed"               // boolean
        OnLeaveNotifyEvent      = "leave-notify-event"          // boolean
        OnMap                   = "map"                         // void
        OnMapEvent              = "map-event"                   // boolean
        OnMnemonicActivate      = "mnemonic-activate"           // boolean
        OnMotionNotifyEvent     = "motion-notify-event"         // boolean
        OnMoveFocus             = "move-focus"                  // void
        OnParentSet             = "parent-set"                  // void
        OnPopupMenu             = "popup-menu"                  // boolean
        OnPropertyNotifyEvent   = "property-notify-event"       // boolean
        OnProximityInEvent      = "proximity-in-event"          // boolean
        OnProximityOutEvent     = "proximity-out-event"         // boolean
        OnQueryTooltip          = "query-tooltip"               // boolean
        OnRealize               = "realize"                     // void
        OnScreenChanged         = "screen-changed"              // void
        OnScrollEvent           = "scroll-event"                // boolean
        OnSelectionClearEvent   = "selection-clear-event"       // boolean
        OnSelectionGet          = "selection-get"               // void
        OnSelectionNotifyEvent  = "selection-notify-event"      // boolean
        OnSelectionReceived     = "selection-received"          // void
        OnSelectionRequestEvent = "selection-request-event"     // boolean
        OnShow                  = "show"                        // void
        OnShowHelp              = "show-help"                   // boolean
        OnSizeAllocate          = "size-allocate"               // void
        OnStateChanged          = "state-changed"               // void
        OnStateFlagsChanged     = "state-flags-changed"         // void
        OnStyleSet              = "style-set"                   // void
        OnStyleUpdated          = "style-updated"               // void
        OnTouchEvent            = "touch-event"                 // boolean
        OnUnmap                 = "unmap"                       // void
        OnUnmapEvent            = "unmap-event"                 // boolean
        OnUnrealize             = "unrealize"                   // void
        OnVisibilityNotifyEvent = "visibility-notify-event"     // boolean
        OnWindowStateEvent      = "window-state-event"          // boolean	
)

func NewPoint(x, y int) PointType {
        return PointType{ image.Pt(x, y) }
}

func NewSize(w, h int) SizeType {
        return SizeType{ image.Pt(w, h) }
}

// Run interaction message loop.
func Interact() {
        runGtkMain()
}

func Quit() {
        quitGtkMain()
}
