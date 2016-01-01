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

func (w *gtkWidget) Get(name PropName) (ValueType, error) {
        return ValueType{ reflect.ValueOf(nil) }, nil
}

func (w *gtkWidget) Set(name PropName, value ValueType) error {
        switch g := w.g(); name {
        case ShowAll:
                if value.IsValid() && value.Bool() {
                        C.gtk_widget_show_all(g)
                } else {
                        C.gtk_widget_hide(g)
                }
        case Size:
                if value.IsValid() && value.CanInterface() {
                        if s, ok := value.Interface().(SizeType); ok {
                                C.gtk_widget_set_size_request(g, (C.gint)(s.X), (C.gint)(s.Y))
                        }
                }
        }
        return nil
}

func (w *gtkWidget) Connect(name SignalName, h interface{}) (Connection, error) {
        o := w.glibInitiallyUnowned.glibObject
        s, e := o.connect(string(name), false, h, nil)
        if e != nil {
                return 0, e
        }
        return Connection(s), nil
}

func (w *gtkWidget) Disconnect(c Connection) (interface{}, error) {
        o := w.glibInitiallyUnowned.glibObject
        return o.disconnect(glibSignalHandle(c)), nil
}

func (w *gtkWindow) Get(name PropName) (ValueType, error) {
        switch g := C.ToWindow(unsafe.Pointer(w.g())); name {
        case Title:
                s := C.GoString((*C.char)(C.gtk_window_get_title(g)));
                return ValueOf(s), nil
        }
        return w.gtkWidget.Get(name)
}

func (w *gtkWindow) Set(name PropName, value ValueType) error {
        switch g := C.ToWindow(unsafe.Pointer(w.g())); name {
        case Title:
                if value.IsValid() && value.Kind() == reflect.String {
                        s := C.CString(value.String()); defer C.free(unsafe.Pointer(s))
                        C.gtk_window_set_title(g, (*C.gchar)(s))
                }
        }
        return w.gtkWidget.Set(name, value)
}

func newGtkView(t ViewClass) View {
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
