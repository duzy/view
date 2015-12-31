// Copyright (c) 2015 Duzy Chan <code@duzy.info>.
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 

package gv

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// #include <stdint.h>
// #include <stdlib.h>
// #include <string.h>
// static inline GObject *ToObject(void *p) { return G_OBJECT(p); }
// static inline GtkWidget *ToWidget(void *p) { return GTK_WIDGET(p); }
// static inline GtkWindow *ToWindow(void *p) { return GTK_WINDOW(p); }
import "C"
import (
        //"os"
        //"fmt"
        //"errors"
        "unsafe"
        "reflect"
	//"sync"
)

type gtkWidget struct {
        glibInitiallyUnowned
}

type gtkWindow struct {
        gtkWidget
}

func init() {
        C.gtk_init(nil, nil)
}

func (w *gtkWidget) g() *C.GtkWidget {
        return C.ToWidget(unsafe.Pointer(w.glibInitiallyUnowned.glibObject.g))
}

func (w *gtkWidget) Get(name PropName) (Value, error) {
        return Value{ reflect.ValueOf(nil) }, nil
}

func (w *gtkWidget) Set(name PropName, value Value) error {
        switch g := w.g(); name {
        case ShowAll:
                if b := value.IsValid() && value.Bool(); b {
                        C.gtk_widget_show_all(g)
                } else {
                        C.gtk_widget_hide(g)
                }
        case Size:
        }
        return nil
}

func (w *gtkWidget) Connect(name SignalName, h interface{}) error {
        return nil
}

func (w *gtkWidget) Disconnect(name SignalName) (interface{}, error) {
        return nil, nil
}

func (w *gtkWindow) Get(name PropName) (Value, error) {
        switch g := C.ToWindow(unsafe.Pointer(w.g())); name {
        case Title:
                s := C.GoString((*C.char)(C.gtk_window_get_title(g)));
                return ValueOf(s), nil
        }
        return w.gtkWidget.Get(name)
}

func (w *gtkWindow) Set(name PropName, value Value) error {
        switch g := C.ToWindow(unsafe.Pointer(w.g())); name {
        case Title:
                if value.IsValid() && value.Kind() == reflect.String {
                        s := C.CString(value.String()); defer C.free(unsafe.Pointer(s))
                        C.gtk_window_set_title(g, (*C.gchar)(s))
                }
        }
        return w.gtkWidget.Set(name, value)
}

func newGtkView(t ViewType) View {
        switch {
        case t == ViewTopLevel:
                p := unsafe.Pointer(C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL));
                w := &gtkWindow{gtkWidget{glibInitiallyUnowned{glibObject{C.ToObject(p)}}}}
                return View(w)
        }
        return nil
}

func runGtkMain() {
        // for { C.gtk_main_iteration_do(true) }
        C.gtk_main()
}

func quitGtkMain() {
	C.gtk_main_quit()
}
