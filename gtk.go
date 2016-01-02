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
// static inline GObject *toObject(void *p) { return G_OBJECT(p); }
// static inline GtkWidget *toWidget(void *p) { return GTK_WIDGET(p); }
// static inline GtkWindow *toWindow(void *p) { return GTK_WINDOW(p); }
// static inline GtkContainer *toContainer(void *p) { return GTK_CONTAINER(p); }
// static inline GtkBox *toBox(void *p) { return GTK_BOX(p); }
// static inline GtkLabel *toLabel(void *p) { return GTK_LABEL(p); }
// static inline GtkEntry *toEntry(void *p) { return GTK_ENTRY(p); }
// static inline GtkTextView *toTextView(void *p) { return GTK_TEXT_VIEW(p); }
// static inline GtkButton *toButton(void *p) { return GTK_BUTTON(p); }
// static inline GtkEditable *toEditable(void *p) { return GTK_EDITABLE(p); }
// static GtkWidget *widget_new(GType type) { return gtk_widget_new(type, NULL); }
// static gchar *text_view_get_text(GtkTextView *tv) {
//   GtkTextIter start, end;
//   GtkTextBuffer *buffer = gtk_text_view_get_buffer(tv);
//   gtk_text_buffer_get_bounds(buffer, &start, &end);
//   return gtk_text_buffer_get_text(buffer, &start, &end, FALSE); // g_free is needed
// }
// static void text_view_set_text(GtkTextView *tv, gchar *t) {
//   GtkTextIter start, end;
//   GtkTextBuffer *buffer = gtk_text_view_get_buffer(tv);
//   gtk_text_buffer_set_text(buffer, t, strlen(t));
// }
import "C"
import (
        //"os"
        "fmt"
        "errors"
        "unsafe"
	//"sync"
        "log"
)

var (
        errNoProperty = errors.New("no property")
        errPropertyNotSet = errors.New("not set property")
        errBadValue = errors.New("bad value")
        errCantAddChild = errors.New("cant add child")
        errIncompatibleView = errors.New("incompatible view")
)

type flagbits uint
const (
        bitPackStart flagbits = 1<<iota
        bitExpend
        bitFill
)

type gtkWidgetCompatible interface {
        widget() *gtkWidget
}

type gtkWidget struct {
        glibInitiallyUnowned
        bits flagbits
}

type gtkBox struct {
        gtkWidget
        padding int
}

type gtkLabel struct {
        gtkWidget
}

type gtkEntry struct {
        gtkWidget
}

type gtkTextView struct {
        gtkWidget
}

type gtkButton struct {
        gtkWidget
}

type gtkWindow struct {
        gtkWidget
}

type gtkEditable struct {
        gtkWidget
}

func init() {
        C.gtk_init(nil, nil)
}

func castBoolValue(value interface{}) (rv bool, e error) {
        switch v := value.(type) {
        case bool:   rv = v
        case string: _, e = fmt.Sscanf(v, "%t", &rv)
        default:     e = errBadValue
        }
        return
}

func castIntValue(value interface{}) (n int, err error) {
        switch v := value.(type) {
        case byte:  n = int(v)
        case int16: n = int(v)
        case int64: n = int(v)
        case int:   n = v
        case string:
                if _, e := fmt.Sscanf(v, "%d", &n); e != nil {
                        err = e
                }
        }
        return
}

func castSizeValue(value interface{}) (sz SizeType, err error) {
        switch s := value.(type) {
        case SizeType: sz = s
        case string:
                if n, e := fmt.Sscanf(s, "%d,%d", &sz.X, &sz.Y); n == 2 && e == nil {
                        return
                } else {
                        log.Printf("SizeValue: %v %v\n", e, s)
                        err = errBadValue
                }
        default: err = errBadValue
        }
        return
}

func (w *gtkWidget) g() *C.GtkWidget {
        return C.toWidget(unsafe.Pointer(w.glibInitiallyUnowned.glibObject.g))
}

func (w *gtkWidget) widget() *gtkWidget {
        return w
}

func (w *gtkWidget) Get(name PropName) (interface{}, error) {
        switch g := w.g(); name {
        case Show:
                return C.gtk_widget_get_visible(g) != 0, nil
        case Size:
                var s SizeType
                C.gtk_widget_get_size_request(g, (*C.gint)(unsafe.Pointer(&s.X)), (*C.gint)(unsafe.Pointer(&s.Y)))
                return s, nil
        case Expend:
                return ((w.bits & bitExpend) != 0), nil
        case Fill:
                return ((w.bits & bitFill) != 0), nil
        }
        return nil, errNoProperty
}

func (w *gtkWidget) Set(name PropName, value interface{}) error {
        switch g := w.g(); name {
        case Show:
                if bv, e := castBoolValue(value); e == nil {
                        if bv {
                                C.gtk_widget_show(g)
                        } else {
                                C.gtk_widget_hide(g)
                        }
                        return nil
                } else {
                        return e
                }
        case Size:
                if sz, e := castSizeValue(value); e == nil {
                        C.gtk_widget_set_size_request(g, (C.gint)(sz.X), (C.gint)(sz.Y))
                        return nil
                } else {
                        return e
                }
                /*
        case Parent:
                if v, ok := value.(gtkWidgetCompatible); ok {
                        if C.gtk_widget_get_parent(g) == nil {
                                C.gtk_widget_set_parent(g, v.widget().g())
                        } else {
                                C.gtk_widget_reparent(g, v.widget().g())
                        }
                        //log.Printf("gtk_widget_set_parent: %v\n", v.widget().g())
                } else {
                        return errPropertyNotSet
                } */
        case Expend:
                if bv, e := castBoolValue(value); e == nil {
                        if bv {
                                w.bits |= bitExpend;
                        } else {
                                w.bits &= ^bitExpend;
                        }
                        return nil
                } else {
                        return e
                }
        case Fill:
                if bv, e := castBoolValue(value); e == nil {
                        if bv {
                                w.bits |= bitFill;
                        } else {
                                w.bits &= ^bitFill;
                        }
                        return nil
                } else {
                        return e
                }
        }
        return errNoProperty
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

func (w *gtkBox) Add(v View) error {
        if wc, ok := v.(gtkWidgetCompatible); ok {
                g, c := C.toBox(unsafe.Pointer(w.g())), wc.widget()
                expend := gboolean((c.bits & bitExpend) != 0)
                fill := gboolean((c.bits & bitFill) != 0)
                if (w.bits & bitPackStart) != 0 {
                        C.gtk_box_pack_start(g, c.g(), expend, fill, C.guint(w.padding))
                } else {
                        C.gtk_box_pack_end(g, c.g(), expend, fill, C.guint(w.padding))
                }
                return nil
        }
        return errIncompatibleView
}

func (w *gtkBox) Find(id string) View {
        return nil
}

func (w *gtkBox) Get(name PropName) (interface{}, error) {
        switch g := C.toBox(unsafe.Pointer(w.g())); name {
        case Pack:
                if (w.bits & bitPackStart) != 0 {
                        return "start", nil
                } else {
                        return "end", nil
                }
        case Spacing:
                return int(C.gtk_box_get_spacing(g)), nil
        case Padding:
                return w.padding, nil
        }
        return w.gtkWidget.Get(name)
}

func (w *gtkBox) Set(name PropName, value interface{}) error {
        switch g := C.toBox(unsafe.Pointer(w.g())); name {
        case Pack:
                if s, ok := value.(string); ok {
                        if s == "start" {
                                w.bits |= bitPackStart
                        } else {
                                w.bits &= ^bitPackStart
                        }
                        return nil
                }
        case Spacing:
                if n, e := castIntValue(value); e == nil {
                        C.gtk_box_set_spacing(g, C.gint(n)); return nil
                } else {
                        return e
                }
        case Padding:
                if n, e := castIntValue(value); e == nil {
                        w.padding = n; return nil
                } else {
                        return e
                }
        }
        return w.gtkWidget.Set(name, value)
}

func (w *gtkLabel) Get(name PropName) (interface{}, error) {
        switch g := C.toLabel(unsafe.Pointer(w.g())); name {
        case Text:
                return C.GoString((*C.char)(C.gtk_label_get_text(g))), nil
        }
        return w.gtkWidget.Get(name)
}

func (w *gtkLabel) Set(name PropName, value interface{}) error {
        switch g := C.toLabel(unsafe.Pointer(w.g())); name {
        case Text:
                if s, ok := value.(string); ok {
                        cs := C.CString(s); defer C.free(unsafe.Pointer(cs))
                        C.gtk_label_set_text(g, (*C.gchar)(cs))
                        return nil
                } else {
                        return errBadValue
                }
        }
        return w.gtkWidget.Set(name, value)
}

func (w *gtkEntry) Get(name PropName) (interface{}, error) {
        switch g := C.toEntry(unsafe.Pointer(w.g())); name {
        case Text:
                return C.GoString((*C.char)(C.gtk_entry_get_text(g))), nil
        }
        return w.gtkWidget.Get(name)
}

func (w *gtkEntry) Set(name PropName, value interface{}) error {
        switch g := C.toEntry(unsafe.Pointer(w.g())); name {
        case Text:
                if s, ok := value.(string); ok {
                        cs := C.CString(s); defer C.free(unsafe.Pointer(cs))
                        C.gtk_entry_set_text(g, (*C.gchar)(cs))
                        return nil
                } else {
                        return errBadValue
                }
        }
        return w.gtkWidget.Set(name, value)
}

func (w *gtkButton) Get(name PropName) (interface{}, error) {
        switch g := C.toButton(unsafe.Pointer(w.g())); name {
        case Text:
                return C.GoString((*C.char)(C.gtk_button_get_label(g))), nil
        }
        return w.gtkWidget.Get(name)
}

func (w *gtkButton) Set(name PropName, value interface{}) error {
        switch g := C.toButton(unsafe.Pointer(w.g())); name {
        case Text:
                if s, ok := value.(string); ok {
                        cs := C.CString(s); defer C.free(unsafe.Pointer(cs))
                        C.gtk_button_set_label(g, (*C.gchar)(cs))
                        return nil
                } else {
                        return errBadValue
                }
        }
        return w.gtkWidget.Set(name, value)
}

func (w *gtkTextView) Get(name PropName) (interface{}, error) {
        //log.Printf("gtkTextView.get")
        switch g := C.toTextView(unsafe.Pointer(w.g())); name {
        case Text:
                text := C.text_view_get_text(g); defer C.g_free(C.gpointer(text))
                return C.GoString((*C.char)(text)), nil
        }
        return w.gtkWidget.Get(name)
}

func (w *gtkTextView) Set(name PropName, value interface{}) error {
        //log.Printf("gtkTextView.set")
        switch g := C.toTextView(unsafe.Pointer(w.g())); name {
        case Text:
                if s, ok := value.(string); ok {
                        cs := C.CString(s); defer C.free(unsafe.Pointer(cs))
                        C.text_view_set_text(g, (*C.gchar)(cs))
                        return nil
                } else {
                        return errBadValue
                }
        }
        return w.gtkWidget.Set(name, value)
}

func (w *gtkWindow) Add(v View) error {
        if c, ok := v.(gtkWidgetCompatible); ok {
                g := C.toContainer(unsafe.Pointer(w.g()))
                C.gtk_container_add(g, c.widget().g())
                return nil
        }
        return errIncompatibleView
}

func (w *gtkWindow) Find(id string) View {
        return nil
}

func (w *gtkWindow) Get(name PropName) (interface{}, error) {
        switch g := C.toWindow(unsafe.Pointer(w.g())); name {
        case Text:
                return C.GoString((*C.char)(C.gtk_window_get_title(g))), nil
        }
        return w.gtkWidget.Get(name)
}

func (w *gtkWindow) Set(name PropName, value interface{}) error {
        switch g := C.toWindow(unsafe.Pointer(w.g())); name {
        case Text:
                if s, ok := value.(string); ok {
                        cs := C.CString(s); defer C.free(unsafe.Pointer(cs))
                        C.gtk_window_set_title(g, (*C.gchar)(cs))
                }
        }
        return w.gtkWidget.Set(name, value)
}

func newGtkWindow() View {
        //p := unsafe.Pointer(C.widget_new(C.gtk_widget_get_type()))
        p := unsafe.Pointer(C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)) // C.GTK_WINDOW_POPUP
        if p == nil {
                log.Fatalf("gtk: cant create window")
                return nil
        }
        w := &gtkWindow{gtkWidget{glibInitiallyUnowned{glibObject{C.toObject(p)}}, 0}}
        return View(w)
}

func newGtkBox(orientation int) View {
        // C.GTK_ORIENTATION_HORIZONTAL==0, C.GTK_ORIENTATION_VERTICAL==1
        p := unsafe.Pointer(C.gtk_box_new(C.GtkOrientation(orientation), C.gint(0)))
        if p == nil {
                log.Fatalf("gtk: cant create box")
                return nil
        }
        w := &gtkBox{gtkWidget{glibInitiallyUnowned{glibObject{C.toObject(p)}}, 0}, 0}
        w.bits |= bitPackStart
        return View(w)
}

func newGtkLabel() View {
        p := unsafe.Pointer(C.widget_new(C.gtk_label_get_type()))
        if p == nil {
                log.Fatalf("gtk: cant create label")
                return nil
        }
        w := &gtkLabel{gtkWidget{glibInitiallyUnowned{glibObject{C.toObject(p)}}, 0}}
        return View(w)
}

func newGtkEntry() View {
        p := unsafe.Pointer(C.widget_new(C.gtk_entry_get_type()))
        if p == nil {
                log.Fatalf("gtk: cant create entry")
                return nil
        }
        w := &gtkEntry{gtkWidget{glibInitiallyUnowned{glibObject{C.toObject(p)}}, 0}}
        return View(w)
}

func newGtkTextView() View {
        p := unsafe.Pointer(C.widget_new(C.gtk_text_view_get_type()))
        if p == nil {
                log.Fatalf("gtk: cant create TextView")
                return nil
        }
        w := &gtkTextView{gtkWidget{glibInitiallyUnowned{glibObject{C.toObject(p)}}, 0}}
        return View(w)
}

func newGtkButton() View {
        p := unsafe.Pointer(C.widget_new(C.gtk_button_get_type()))
        if p == nil {
                log.Fatalf("gtk: cant create button")
                return nil
        }
        w := &gtkButton{gtkWidget{glibInitiallyUnowned{glibObject{C.toObject(p)}}, 0}}
        return View(w)
}

func newGtkEditable() View {
        p := unsafe.Pointer(C.widget_new(C.gtk_editable_get_type()))
        if p == nil {
                log.Fatalf("gtk: cant create editable")
                return nil
        }
        //w := &gtkEditable{glibInitiallyUnowned{glibObject{C.toObject(p)}}}
        return nil //View(w)
}

func runGtkMain() {
        // for { C.gtk_main_iteration_do(true) }
        C.gtk_main()
}

func quitGtkMain() {
	C.gtk_main_quit()
}
