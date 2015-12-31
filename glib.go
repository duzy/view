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
// extern void glibMarshal(GClosure *, GValue *, guint, GValue *, gpointer, GValue *);
// static gboolean is_value(GValue *val) { return G_IS_VALUE(val); }
// static GType value_type(GValue *val) { return G_VALUE_TYPE(val); }
// static GType value_fundamental(GType type) { return G_TYPE_FUNDAMENTAL(type); }
// static GValue *value_init(GType type) { return g_value_init(g_new0(GValue, 1), type); }
// static inline GClosure * closure_new() {
//     GClosure *closure = g_closure_new_simple(sizeof(GClosure), NULL);
//     g_closure_set_marshal(closure, (GClosureMarshal)(glibMarshal));
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

// glibType is a representation of GType.
type glibType uint

// g_type_name().
func (t glibType) name() string {
	return C.GoString((*C.char)(C.g_type_name(C.GType(t))))
}

// g_type_depth().
func (t glibType) depth() uint {
	return uint(C.g_type_depth(C.GType(t)))
}

// g_type_parent().
func (t glibType) parent() glibType {
	return glibType(C.g_type_parent(C.GType(t)))
}

const (
	glibType_INVALID   glibType = C.G_TYPE_INVALID
	glibType_NONE      glibType = C.G_TYPE_NONE
	glibType_INTERFACE glibType = C.G_TYPE_INTERFACE
	glibType_CHAR      glibType = C.G_TYPE_CHAR
	glibType_UCHAR     glibType = C.G_TYPE_UCHAR
	glibType_BOOLEAN   glibType = C.G_TYPE_BOOLEAN
	glibType_INT       glibType = C.G_TYPE_INT
	glibType_UINT      glibType = C.G_TYPE_UINT
	glibType_LONG      glibType = C.G_TYPE_LONG
	glibType_ULONG     glibType = C.G_TYPE_ULONG
	glibType_INT64     glibType = C.G_TYPE_INT64
	glibType_UINT64    glibType = C.G_TYPE_UINT64
	glibType_ENUM      glibType = C.G_TYPE_ENUM
	glibType_FLAGS     glibType = C.G_TYPE_FLAGS
	glibType_FLOAT     glibType = C.G_TYPE_FLOAT
	glibType_DOUBLE    glibType = C.G_TYPE_DOUBLE
	glibType_STRING    glibType = C.G_TYPE_STRING
	glibType_POINTER   glibType = C.G_TYPE_POINTER
	glibType_BOXED     glibType = C.G_TYPE_BOXED
	glibType_PARAM     glibType = C.G_TYPE_PARAM
	glibType_OBJECT    glibType = C.G_TYPE_OBJECT
	glibType_VARIANT   glibType = C.G_TYPE_VARIANT
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

var (
        errNilPtr = errors.New("cgo returned unexpected nil pointer")

	glibClosures = struct {
		sync.RWMutex
		m map[*C.GClosure]glibClosureContext
	}{
		m: make(map[*C.GClosure]glibClosureContext),
	}

        glibMarshalers = map[glibType] func(uintptr) (interface{}, error) {
                glibType_INVALID:   glibMarshalInvalid,
                glibType_NONE:      glibMarshalNone,
                glibType_INTERFACE: glibMarshalInterface,
                glibType_CHAR:      glibMarshalChar,
                glibType_UCHAR:     glibMarshalUchar,
                glibType_BOOLEAN:   glibMarshalBoolean,
                glibType_INT:       glibMarshalInt,
                glibType_LONG:      glibMarshalLong,
                glibType_ENUM:      glibMarshalEnum,
                glibType_INT64:     glibMarshalInt64,
                glibType_UINT:      glibMarshalUint,
                glibType_ULONG:     glibMarshalUlong,
                glibType_FLAGS:     glibMarshalFlags,
                glibType_UINT64:    glibMarshalUint64,
                glibType_FLOAT:     glibMarshalFloat,
                glibType_DOUBLE:    glibMarshalDouble,
                glibType_STRING:    glibMarshalString,
                glibType_POINTER:   glibMarshalPointer,
                glibType_BOXED:     glibMarshalBoxed,
                glibType_OBJECT:    glibMarshalObject,
                glibType_VARIANT:   glibMarshalVariant,
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

func glibValueInit(t glibType) (*glibValue, error) {
        if val := C.value_init(C.GType(glibType_POINTER)); val != nil {
                v := &glibValue{ val }
                runtime.SetFinalizer(v, (*glibValue).unset)
                return v, nil
        }
        return nil, errNilPtr
}

func glibUnmarshalValue(v interface{}) (*glibValue, error) {
	if v == nil {
		if val, err := glibValueInit(glibType_POINTER); err != nil {
			return nil, err
		} else {
                        C.g_value_set_pointer(val.g, C.gpointer(uintptr(unsafe.Pointer(nil))))
                        return val, nil
                }
	}

	switch e := v.(type) {
	case bool:
		val, err := glibValueInit(glibType_BOOLEAN)
		if err != nil {
			return nil, err
		}
                C.g_value_set_boolean(val.g, gboolean(e))
		return val, nil

	case int8:
		val, err := glibValueInit(glibType_CHAR)
		if err != nil {
			return nil, err
		}
                C.g_value_set_schar(val.g, C.gint8(e))
		return val, nil

	case int64:
		val, err := glibValueInit(glibType_INT64)
		if err != nil {
			return nil, err
		}
                C.g_value_set_int64(val.g, C.gint64(e))
		return val, nil

	case int:
		val, err := glibValueInit(glibType_INT)
		if err != nil {
			return nil, err
		}
                C.g_value_set_int(val.g, C.gint(e))
		return val, nil

	case uint8:
		val, err := glibValueInit(glibType_UCHAR)
		if err != nil {
			return nil, err
		}
                C.g_value_set_uchar(val.g, C.guchar(e))
		return val, nil

	case uint64:
		val, err := glibValueInit(glibType_UINT64)
		if err != nil {
			return nil, err
		}
                C.g_value_set_uint64(val.g, C.guint64(e))
		return val, nil

	case uint:
		val, err := glibValueInit(glibType_UINT)
		if err != nil {
			return nil, err
		}
                C.g_value_set_uint(val.g, C.guint(e))
		return val, nil

	case float32:
		val, err := glibValueInit(glibType_FLOAT)
		if err != nil {
			return nil, err
		}
                C.g_value_set_float(val.g, C.gfloat(e))
		return val, nil

	case float64:
		val, err := glibValueInit(glibType_DOUBLE)
		if err != nil {
			return nil, err
		}
                C.g_value_set_double(val.g, C.gdouble(e))
		return val, nil

	case string:
		val, err := glibValueInit(glibType_STRING)
		if err != nil {
			return nil, err
		}
                cstr := C.CString(e); defer C.free(unsafe.Pointer(cstr))
                C.g_value_set_string(val.g, (*C.gchar)(cstr))
		return val, nil

	case *glibObject:
		val, err := glibValueInit(glibType_OBJECT)
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
                val, err := glibValueInit(glibType_CHAR)
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
                val, err := glibValueInit(glibType_INT64)
                if err != nil {
                        return nil, err
                }
                C.g_value_set_int64(val.g, C.gint64(rval.Int()))
                return val, nil

        case reflect.Int:
                val, err := glibValueInit(glibType_INT)
                if err != nil {
                        return nil, err
                }
                C.g_value_set_int(val.g, C.gint(rval.Int()))
                return val, nil

        case reflect.Uintptr, reflect.Ptr:
                val, err := glibValueInit(glibType_POINTER)
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
        if f, ok := glibMarshalers[glibType(actual)]; ok {
                return f(uintptr(unsafe.Pointer(v)))
        }

	fundamental := C.value_fundamental(actual)
        if f, ok := glibMarshalers[glibType(fundamental)]; ok {
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

func glibConnect(signal string, after bool, f interface{}, data ...interface{}) error {
        fv := reflect.ValueOf(f)

        if fv.Type().Kind() != reflect.Func {
                return errors.New("value is not a func")
        }
        
        s := C.CString(signal); defer C.free(unsafe.Pointer(s))

        c := glibClosureContext{ f:fv }
        if 0 < len(data) {
                c.d = reflect.ValueOf(data)
        }
        
        closure :=  C.closure_new()
	glibClosures.RLock()
	glibClosures.m[closure] = c
	glibClosures.RUnlock()

        return nil
}
