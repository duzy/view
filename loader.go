// Copyright (c) 2015 Duzy Chan <code@duzy.info>.
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 

package gv

import (
        "os"
        "io"
        "log"
        "fmt"
        "encoding/xml"
        "bytes"
        "strings"
        "errors"
)

type stack struct {
        next *stack
        name xml.Name
        view View
}

type builder struct {
        decoder *xml.Decoder
        top *stack
}

var (
        viewCreators = map[string] func(a []xml.Attr) View {
                "window": createWindow,
                "static": createStatic,
                "editable": createEditable,
                "pushable": createPushable,
                "horizontal": createHorizontal,
                "vertical": createVertical,
                "h": createHorizontal,
                "v": createVertical,
        }
)

// build recursively create views for a token 
func (b *builder) build() (v View, err error) {
        // We don't use Unmarshal, because we need to go through the
        // elements and create things over the traversal.
        for {
                t, e := b.decoder.Token()
                if e != nil {
                        if e == io.EOF {
                                break
                        }
                        log.Fatalf("invalid token: %v\n", e)
                        err = e
                }

                log.Printf("token: %T: %v\n", t, t)

                switch t := t.(type) {
                case xml.StartElement:
                        if e = b.push(t); e != nil {
                                return nil, e
                        }
                case xml.EndElement:
                        v = b.top.view
                        if e = b.pop(t); e != nil {
                                return nil, e
                        }
                }
        }

        log.Fatal("not fully built")
        return nil, errors.New("partially built")
}

func (b *builder) push(t xml.StartElement) error {
        create, ok := viewCreators[t.Name.Local]
        if !ok {
                log.Fatalf("unknown view %v\n", t.Name.Local)
                return errors.New("unknown view " + t.Name.Local)
        }

        v := create(t.Attr)
        if v == nil {
                log.Fatalf("cant create view %v\n", t.Name.Local)
                return errors.New("unknown view " + t.Name.Local)
        }

        if b.top != nil {
                if c, ok := b.top.view.(adder); ok {
                        if e := c.Add(v); e != nil {
                                log.Fatalf("%v: %v %v\n", b.top.name.Local, e, t.Name.Local)
                                return e
                        }
                }
        }

        b.top = &stack{ next:b.top, name:t.Name, view:v }
        return nil
}

func (b *builder) pop(t xml.EndElement) error {
        if b.top == nil {
                log.Fatalf("empty view stack: %v\n", t.Name.Local)
                return errors.New("bad view stack " + t.Name.Local)
        }

        b.top = b.top.next
        return nil
}

func applyViewAttr(v View, a []xml.Attr) View {
        hasShow := false

        for _, i := range a {
                if i.Name.Space == "-" { continue }
                if e := v.Set(PropName(i.Name.Local), i.Value); e != nil {
                        log.Printf("attr: %v %v\n", e, i.Name.Local)
                }
                if i.Name.Local == string(Show) {
                        hasShow = true
                }
        }

        if !hasShow {
                v.Set(Show, true)
        }

        return v
}

func createWindow(a []xml.Attr) View {
        return applyViewAttr(newGtkWindow(), a)
}

func createView(a []xml.Attr) View {
        const (
                horizontal = 0
                vertical = 1
                unknown
        )

        t, tt := unknown, horizontal
        for _, i := range a {
                switch i.Name.Local {
                case "vertical": if t == unknown { t = vertical }; fallthrough
                case "v": if t == unknown { t = vertical }; fallthrough
                case "horizontal": fallthrough
                case "h":
                        if t == unknown { t = horizontal }

                        bv := false
                        if _, e := fmt.Sscanf(i.Value, "%t", &bv); e == nil {
                                if bv { tt = t }
                        } else {
                                log.Printf("attr: %v %v\n", e, i.Name.Local)
                        }
                }
        }

        return applyViewAttr(newGtkBox(tt), a)
}

func createStatic(a []xml.Attr) View {
        return applyViewAttr(newGtkLabel(), a)
}

func createEditable(a []xml.Attr) View {
        return applyViewAttr(newGtkEntry(), a)
}

func createPushable(a []xml.Attr) View {
        return applyViewAttr(newGtkButton(), a)
}

func createHorizontal(a []xml.Attr) View {
        return applyViewAttr(newGtkBox(0), a)
}

func createVertical(a []xml.Attr) View {
        return applyViewAttr(newGtkBox(1), a)
}

// Load loads views from a reader.
func Load(in io.Reader) (View, error) {
        buf := new(bytes.Buffer)

        if _, e := io.Copy(buf, in); e != nil {
                return nil, e
        }

        b := &builder{ xml.NewDecoder(buf), nil }
        return b.build()
}

// LoadString loads views from XML string.
func LoadString(s string) (View, error) {
        return Load(strings.NewReader(s))
}

// LoadFile loads views from XML file.
func LoadFile(name string) (View, error) {
        f, e := os.Open(name)
        if e != nil {
                return nil, e
        }

        v, e := Load(f)
        if e != nil || v == nil {
                return nil, e
        }

        if i, e := v.Get(Text); e == nil {
                if s, ok := i.(string); ok && s == "" {
                        v.Set(Text, name)
                }
        }

        return v, nil
}
