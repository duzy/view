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
        //"fmt"
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
        text *bytes.Buffer
        top *stack
}

type signalCreator interface {
        createSignal(name string) error
}

var (
        viewCreators = map[string] func(a []xml.Attr) View {
                "window": createWindow,
                "static": createStatic,
                "editable": createEditable,
                "edit": createTextView,
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

                //log.Printf("token: %T: %v\n", t, t)

                switch t := t.(type) {
                case xml.StartElement:
                        if e = b.push(t); e != nil {
                                return nil, e
                        }
                case xml.EndElement:
                        if b.top != nil { v = b.top.view }
                        if e = b.pop(t); e != nil {
                                return nil, e
                        }
                case xml.CharData:
                        if _, e := b.text.Write([]byte(t)); e != nil {
                                return nil, e
                        }
                }
        }

        if v == nil {
                log.Fatal("partially built")
                return nil, errors.New("partially built")
        }
        return v, nil
}

func (b *builder) push(t xml.StartElement) error {
        if t.Name.Local == "signal" {
                if b.top == nil {
                        log.Fatalf("no view for signal")
                        return errors.New("no view for signal")
                }

                if e := createSignal(b.top.view, t.Attr); e != nil {
                        log.Fatalf("no view for signal")
                        return e
                }

                return nil
        }

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
                if id := getAttrByName(t.Attr, "", "id"); id != nil {
                        for s := b.top; s != nil; s = s.next {
                                if f, ok := s.view.(Finder); ok {
                                        if e := f.insert(id.Value, v); e != nil {
                                                log.Fatalf("%v.%v: %v\n", b.top.name.Local, t.Name.Local, e)
                                                return e
                                        }
                                }
                        }
                }
        }

        b.top = &stack{ next:b.top, name:t.Name, view:v }
        return nil
}

func (b *builder) pop(t xml.EndElement) error {
        if t.Name.Local == "signal" {
                return nil
        }

        if b.top == nil {
                log.Fatalf("empty view stack: %v\n", t.Name.Local)
                return errors.New("bad view stack " + t.Name.Local)
        }

        if s := strings.TrimSpace(b.text.String()); s != "" {
                b.top.view.Set(Text, prettify(s))
        }

        b.text.Reset()
        b.top = b.top.next
        return nil
}

func prettify(s string) string {
        // TODO: prettify text?
        return s
}

func getAttrByName(a []xml.Attr, space, local string) *xml.Attr {
        for _, i := range a {
                if i.Name.Space == space && i.Name.Local == local {
                        return &i
                }
        }
        return nil
}

func applyViewAttr(v View, a []xml.Attr) View {
        hasShow := false

        for _, i := range a {
                if i.Name.Space == "-" || i.Name.Local == "id" { continue }
                if e := v.Set(PropName(i.Name.Local), i.Value); e != nil {
                        log.Fatalf("attribute: %v %v\n", e, i.Name.Local)
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

func createSignal(v View, a []xml.Attr) error {
        name := getAttrByName(a, "", "name")
        if name == nil {
                log.Fatalf("signal: no name property")
                return errors.New("signal: no name property")
        }

        if sc, ok := v.(signalCreator); ok {
                return sc.createSignal(name.Value)
        }

        log.Fatalf("view cant have signal")
        return errors.New("view cant have signal")
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
                        if bv, e := castBoolValue(i.Value); e == nil {
                                if bv { tt = t }
                        } else {
                                log.Fatalf("%v: not boolean: %v (%v)\n", i.Name.Local, i.Value, e)
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

func createTextView(a []xml.Attr) View {
        return applyViewAttr(newGtkTextView(), a)
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

        b := &builder{ xml.NewDecoder(buf), new(bytes.Buffer) , nil }
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
