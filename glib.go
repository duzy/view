// Copyright (c) 2015 Duzy Chan <code@duzy.info>.
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 

package gv

// #cgo pkg-config: gtk+-3.0
// #include <gio/gio.h>
// #include <glib.h>
// #include <glib-object.h>
// #include <stdint.h>
// #include <stdlib.h>
// #include <string.h>
// extern void glibRemoveClosure(gpointer, GClosure *);
// extern void glibMarshal(GClosure *, GValue *, guint, GValue *, gpointer, GValue *);
// static gboolean is_value(GValue *val) { return G_IS_VALUE(val); }
// static GType value_type(GValue *val) { return G_VALUE_TYPE(val); }
// static GType value_fundamental(GType type) { return G_TYPE_FUNDAMENTAL(type); }
// static GValue *value_init(GType type) { return g_value_init(g_new0(GValue, 1), type); }
// static inline GClosure * closure_new() {
//     GClosure *closure = g_closure_new_simple(sizeof(GClosure), NULL);
//     g_closure_set_marshal(closure, (GClosureMarshal)(glibMarshal));
//     g_closure_add_finalize_notifier(closure, NULL, glibRemoveClosure);
//     return closure;
// }
import "C"
import (
        "os"
        "fmt"
        "errors"
        "unsafe"
        "reflect"
        "runtime"
	"sync"
)

// glibKind is a representation of GType.
type glibKind uint

// g_type_name().
func (t glibKind) name() string {
	return C.GoString((*C.char)(C.g_type_name(C.GType(t))))
}

// g_type_depth().
func (t glibKind) depth() uint {
	return uint(C.g_type_depth(C.GType(t)))
}

// g_type_parent().
func (t glibKind) parent() glibKind {
	return glibKind(C.g_type_parent(C.GType(t)))
}

const (
	glibKind_INVALID   glibKind = C.G_TYPE_INVALID
	glibKind_NONE      glibKind = C.G_TYPE_NONE
	glibKind_INTERFACE glibKind = C.G_TYPE_INTERFACE
	glibKind_CHAR      glibKind = C.G_TYPE_CHAR
	glibKind_UCHAR     glibKind = C.G_TYPE_UCHAR
	glibKind_BOOLEAN   glibKind = C.G_TYPE_BOOLEAN
	glibKind_INT       glibKind = C.G_TYPE_INT
	glibKind_UINT      glibKind = C.G_TYPE_UINT
	glibKind_LONG      glibKind = C.G_TYPE_LONG
	glibKind_ULONG     glibKind = C.G_TYPE_ULONG
	glibKind_INT64     glibKind = C.G_TYPE_INT64
	glibKind_UINT64    glibKind = C.G_TYPE_UINT64
	glibKind_ENUM      glibKind = C.G_TYPE_ENUM
	glibKind_FLAGS     glibKind = C.G_TYPE_FLAGS
	glibKind_FLOAT     glibKind = C.G_TYPE_FLOAT
	glibKind_DOUBLE    glibKind = C.G_TYPE_DOUBLE
	glibKind_STRING    glibKind = C.G_TYPE_STRING
	glibKind_POINTER   glibKind = C.G_TYPE_POINTER
	glibKind_BOXED     glibKind = C.G_TYPE_BOXED
	glibKind_PARAM     glibKind = C.G_TYPE_PARAM
	glibKind_OBJECT    glibKind = C.G_TYPE_OBJECT
	glibKind_VARIANT   glibKind = C.G_TYPE_VARIANT
)

type glibValue struct {
        g *C.GValue
}

type glibObject struct {
        g *C.GObject
}

type glibInitiallyUnowned struct {
        glibObject
}

type glibClosureContext struct {
	f reflect.Value
	d reflect.Value
}

type glibSignalHandle uint

var (
        errNilPtr = errors.New("cgo returned unexpected nil pointer")

	glibClosures = struct {
		sync.RWMutex
		m map[*C.GClosure]glibClosureContext
	}{
		m: make(map[*C.GClosure]glibClosureContext),
	}

	glibSignals = make(map[glibSignalHandle]*C.GClosure)

        glibMarshalers = map[glibKind] func(uintptr) (interface{}, error) {
                glibKind_INVALID:   glibMarshalInvalid,
                glibKind_NONE:      glibMarshalNone,
                glibKind_INTERFACE: glibMarshalInterface,
                glibKind_CHAR:      glibMarshalChar,
                glibKind_UCHAR:     glibMarshalUchar,
                glibKind_BOOLEAN:   glibMarshalBoolean,
                glibKind_INT:       glibMarshalInt,
                glibKind_LONG:      glibMarshalLong,
                glibKind_ENUM:      glibMarshalEnum,
                glibKind_INT64:     glibMarshalInt64,
                glibKind_UINT:      glibMarshalUint,
                glibKind_ULONG:     glibMarshalUlong,
                glibKind_FLAGS:     glibMarshalFlags,
                glibKind_UINT64:    glibMarshalUint64,
                glibKind_FLOAT:     glibMarshalFloat,
                glibKind_DOUBLE:    glibMarshalDouble,
                glibKind_STRING:    glibMarshalString,
                glibKind_POINTER:   glibMarshalPointer,
                glibKind_BOXED:     glibMarshalBoxed,
                glibKind_OBJECT:    glibMarshalObject,
                glibKind_VARIANT:   glibMarshalVariant,
        }
)

func gboolean(b bool) C.gboolean {
        if b {
                return C.gboolean(1)
        }
        return C.gboolean(0)
}

func (v *glibValue) unset() {
	C.g_value_unset(v.g)
}

func glibValueInit(t glibKind) (*glibValue, error) {
        if val := C.value_init(C.GType(glibKind_POINTER)); val != nil {
                v := &glibValue{ val }
                runtime.SetFinalizer(v, (*glibValue).unset)
                return v, nil
        }
        return nil, errNilPtr
}

func glibUnmarshalValue(v interface{}) (*glibValue, error) {
	if v == nil {
		if val, err := glibValueInit(glibKind_POINTER); err != nil {
			return nil, err
		} else {
                        C.g_value_set_pointer(val.g, C.gpointer(uintptr(unsafe.Pointer(nil))))
                        return val, nil
                }
	}

	switch e := v.(type) {
	case bool:
		val, err := glibValueInit(glibKind_BOOLEAN)
		if err != nil {
			return nil, err
		}
                C.g_value_set_boolean(val.g, gboolean(e))
		return val, nil

	case int8:
		val, err := glibValueInit(glibKind_CHAR)
		if err != nil {
			return nil, err
		}
                C.g_value_set_schar(val.g, C.gint8(e))
		return val, nil

	case int64:
		val, err := glibValueInit(glibKind_INT64)
		if err != nil {
			return nil, err
		}
                C.g_value_set_int64(val.g, C.gint64(e))
		return val, nil

	case int:
		val, err := glibValueInit(glibKind_INT)
		if err != nil {
			return nil, err
		}
                C.g_value_set_int(val.g, C.gint(e))
		return val, nil

	case uint8:
		val, err := glibValueInit(glibKind_UCHAR)
		if err != nil {
			return nil, err
		}
                C.g_value_set_uchar(val.g, C.guchar(e))
		return val, nil

	case uint64:
		val, err := glibValueInit(glibKind_UINT64)
		if err != nil {
			return nil, err
		}
                C.g_value_set_uint64(val.g, C.guint64(e))
		return val, nil

	case uint:
		val, err := glibValueInit(glibKind_UINT)
		if err != nil {
			return nil, err
		}
                C.g_value_set_uint(val.g, C.guint(e))
		return val, nil

	case float32:
		val, err := glibValueInit(glibKind_FLOAT)
		if err != nil {
			return nil, err
		}
                C.g_value_set_float(val.g, C.gfloat(e))
		return val, nil

	case float64:
		val, err := glibValueInit(glibKind_DOUBLE)
		if err != nil {
			return nil, err
		}
                C.g_value_set_double(val.g, C.gdouble(e))
		return val, nil

	case string:
		val, err := glibValueInit(glibKind_STRING)
		if err != nil {
			return nil, err
		}
                cstr := C.CString(e); defer C.free(unsafe.Pointer(cstr))
                C.g_value_set_string(val.g, (*C.gchar)(cstr))
		return val, nil

	case *glibObject:
		val, err := glibValueInit(glibKind_OBJECT)
		if err != nil {
			return nil, err
		}
                C.g_value_set_instance(val.g, C.gpointer(e.g))
		return val, nil

	default:
		/* Try this since above doesn't catch constants under other types */
		glibUnmarshalReflectedValue(reflect.ValueOf(v))
	}

	return nil, errors.New("type not implemented")
}

func glibUnmarshalReflectedValue(rval reflect.Value) (*glibValue, error) {
        switch rval.Kind() {
        case reflect.Int8:
                val, err := glibValueInit(glibKind_CHAR)
                if err != nil {
                        return nil, err
                }
                C.g_value_set_schar(val.g, C.gint8(rval.Int()))
                return val, nil

        case reflect.Int16:
                return nil, errors.New("Type not implemented")

        case reflect.Int32:
                return nil, errors.New("Type not implemented")

        case reflect.Int64:
                val, err := glibValueInit(glibKind_INT64)
                if err != nil {
                        return nil, err
                }
                C.g_value_set_int64(val.g, C.gint64(rval.Int()))
                return val, nil

        case reflect.Int:
                val, err := glibValueInit(glibKind_INT)
                if err != nil {
                        return nil, err
                }
                C.g_value_set_int(val.g, C.gint(rval.Int()))
                return val, nil

        case reflect.Uintptr, reflect.Ptr:
                val, err := glibValueInit(glibKind_POINTER)
                if err != nil {
                        return nil, err
                }
                C.g_value_set_pointer(val.g, C.gpointer(uintptr(unsafe.Pointer(rval.Pointer()))))
                return val, nil
        }
	return nil, errors.New("type not implemented")
}

// glibMarshalValue converts a GValue to comparable Go type.
func glibMarshalValue(v *C.GValue) (interface{}, error) {
	if C.is_value(v) != 0 {
		return nil, errors.New("invalid GValue")
	}

	actual := C.value_type(v)
        if f, ok := glibMarshalers[glibKind(actual)]; ok {
                return f(uintptr(unsafe.Pointer(v)))
        }

	fundamental := C.value_fundamental(actual)
        if f, ok := glibMarshalers[glibKind(fundamental)]; ok {
                return f(uintptr(unsafe.Pointer(v)))
        }

	return nil, errors.New("missing marshaler")
}

func glibMarshalInvalid(uintptr) (interface{}, error) {
	return nil, errors.New("invalid type")
}

func glibMarshalNone(uintptr) (interface{}, error) {
	return nil, nil
}

func glibMarshalInterface(uintptr) (interface{}, error) {
	return nil, errors.New("interface conversion not yet implemented")
}

func glibMarshalChar(p uintptr) (interface{}, error) {
	c := C.g_value_get_schar((*C.GValue)(unsafe.Pointer(p)))
	return int8(c), nil
}

func glibMarshalUchar(p uintptr) (interface{}, error) {
	c := C.g_value_get_uchar((*C.GValue)(unsafe.Pointer(p)))
	return uint8(c), nil
}

func glibMarshalBoolean(p uintptr) (interface{}, error) {
	c := C.g_value_get_boolean((*C.GValue)(unsafe.Pointer(p)))
	return bool(c != 0), nil
}

func glibMarshalInt(p uintptr) (interface{}, error) {
	c := C.g_value_get_int((*C.GValue)(unsafe.Pointer(p)))
	return int(c), nil
}

func glibMarshalLong(p uintptr) (interface{}, error) {
	c := C.g_value_get_long((*C.GValue)(unsafe.Pointer(p)))
	return int(c), nil
}

func glibMarshalEnum(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return int(c), nil
}

func glibMarshalInt64(p uintptr) (interface{}, error) {
	c := C.g_value_get_int64((*C.GValue)(unsafe.Pointer(p)))
	return int64(c), nil
}

func glibMarshalUint(p uintptr) (interface{}, error) {
	c := C.g_value_get_uint((*C.GValue)(unsafe.Pointer(p)))
	return uint(c), nil
}

func glibMarshalUlong(p uintptr) (interface{}, error) {
	c := C.g_value_get_ulong((*C.GValue)(unsafe.Pointer(p)))
	return uint(c), nil
}

func glibMarshalFlags(p uintptr) (interface{}, error) {
	c := C.g_value_get_flags((*C.GValue)(unsafe.Pointer(p)))
	return uint(c), nil
}

func glibMarshalUint64(p uintptr) (interface{}, error) {
	c := C.g_value_get_uint64((*C.GValue)(unsafe.Pointer(p)))
	return uint64(c), nil
}

func glibMarshalFloat(p uintptr) (interface{}, error) {
	c := C.g_value_get_float((*C.GValue)(unsafe.Pointer(p)))
	return float32(c), nil
}

func glibMarshalDouble(p uintptr) (interface{}, error) {
	c := C.g_value_get_double((*C.GValue)(unsafe.Pointer(p)))
	return float64(c), nil
}

func glibMarshalString(p uintptr) (interface{}, error) {
	c := C.g_value_get_string((*C.GValue)(unsafe.Pointer(p)))
	return C.GoString((*C.char)(c)), nil
}

func glibMarshalBoxed(p uintptr) (interface{}, error) {
	c := C.g_value_get_boxed((*C.GValue)(unsafe.Pointer(p)))
	return uintptr(unsafe.Pointer(c)), nil
}

func glibMarshalPointer(p uintptr) (interface{}, error) {
	c := C.g_value_get_pointer((*C.GValue)(unsafe.Pointer(p)))
	return unsafe.Pointer(c), nil
}

func glibMarshalObject(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	return &glibObject{ (*C.GObject)(c) }, nil
}

func glibMarshalVariant(p uintptr) (interface{}, error) {
	return nil, errors.New("variant conversion not yet implemented")
}

//export glibMarshal
func glibMarshal(closure *C.GClosure, retValue *C.GValue, nParams C.guint, params *C.GValue,
	invocationHint C.gpointer, marshalData *C.GValue) {

	glibClosures.RLock()
	c := glibClosures.m[closure]
	glibClosures.RUnlock()

        nGlibParams := int(nParams)
        nAllParams := nGlibParams
        nCallParams := c.f.Type().NumIn()

        if c.d.IsValid() {
                nAllParams++
        }

        if nAllParams < nCallParams {
		fmt.Fprintf(os.Stderr, "too many closure args: have %d, max allowed %d\n",
			nCallParams, nAllParams)
                return
        }

        values := []C.GValue{}
        setSliceHeader(unsafe.Pointer(&values), unsafe.Pointer(params), nCallParams)

        args := make([]reflect.Value, 0, nCallParams)
        for i := 0; i < nCallParams && i < nGlibParams; i++ {
		if v, e := glibMarshalValue(&values[i]); e != nil {
			fmt.Fprintf(os.Stderr, "no suitable Go value for arg %d: %v\n", i, e)
			return
		} else {
                        args = append(args, reflect.ValueOf(v).Convert(c.f.Type().In(i)))
                }
        }

	// If non-nil user data was passed in and not all args have been set.
        if c.d.IsValid() && len(args) < cap(args) {
                // Get and set the reflect.Value directly from the GValue.
		args = append(args, c.d.Convert(c.f.Type().In(nCallParams-1)))
        }

	// Call closure with args. If the callback returns one or more
	// values, save the GValue equivalent of the first.
        if rv := c.f.Call(args); retValue != nil && 0 < len(rv) {
 		if v, e := glibUnmarshalValue(rv[0].Interface()); e != nil {
			fmt.Fprintf(os.Stderr, "cannot unmarshal callback return value: %v", e)
		} else {
			*retValue = *v.g
		}
        }
}

//export glibRemoveClosure
func glibRemoveClosure(_ C.gpointer, closure *C.GClosure) {
	glibClosures.RLock()
	delete(glibClosures.m, closure)
	glibClosures.RUnlock()
}

func glibNewClosure(f interface{}, data ...interface{}) (*C.GClosure, error) {
        fv := reflect.ValueOf(f)

        if fv.Type().Kind() != reflect.Func {
                return nil, errors.New("value is not a func")
        }
        
        c := glibClosureContext{ f:fv }
        if 0 < len(data) {
                c.d = reflect.ValueOf(data)
        }
        
        closure :=  C.closure_new()
	glibClosures.RLock()
	glibClosures.m[closure] = c
	glibClosures.RUnlock()
        return closure, nil
}

func (obj *glibObject) connect(signal string, after bool, f interface{}, data ...interface{}) (glibSignalHandle, error) {
	if 1 < len(data) {
		return 0, errors.New("user data len must be 0 or 1")
	}

        s := C.CString(signal); defer C.free(unsafe.Pointer(s))

        closure, err := glibNewClosure(f, data)
	if err != nil {
		return 0, err
	}

	handle := glibSignalHandle(C.g_signal_connect_closure(C.gpointer(obj.g),
		(*C.gchar)(s), closure, gboolean(after)))

	// Map the signal handle to the closure.
	glibSignals[handle] = closure

	return handle, nil
}

func (o *glibObject) disconnect(handle glibSignalHandle) (f interface{}) {
        if closure, ok := glibSignals[handle]; ok {
                C.g_signal_handler_disconnect(C.gpointer(o.g), C.gulong(handle))
                C.g_closure_invalidate(closure)

                //glibRemoveClosure(nil, closure) // delete(glibClosures.m, closure)
                glibClosures.RLock()
                f = glibClosures.m[closure].f
                delete(glibClosures.m, closure)
                glibClosures.RUnlock()

                delete(glibSignals, handle)
        }
        return
}

func (o *glibObject) block(handle glibSignalHandle) {
	C.g_signal_handler_block(C.gpointer(o.g), C.gulong(handle))
}

func (o *glibObject) unblock(handle glibSignalHandle) {
	C.g_signal_handler_unblock(C.gpointer(o.g), C.gulong(handle))
}
