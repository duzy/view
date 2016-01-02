// Copyright (c) 2015 Duzy Chan <code@duzy.info>.
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 

package gv_test

import (
        "log"
        "testing"
        . "github.com/duzy/gv"
)

var strExampleXml = `<?xml version="1.0" encoding="UTF-8"?>
<window text="Restuarant Scraper" size="800,600">
  <vertical spacing="10" padding="5">
    <horizontal padding="5">
      <static text="static 1" show="true" />
      <static text="static 2" expend="true" show="true" />
      <static text="static 3" show="true" />
    </horizontal>
    <h padding="5"><static text="label: " /><editable text="editable" fill="true" /></h>
    <h expend="true" />
    <horizontal padding="5" pack="end">
      <pushable text="pushable 1" show="true" />
      <pushable text="pushable 2" show="true" fill="true" />
      <pushable text="pushable 3" show="true" fill="true" />
    </horizontal>
  </vertical>
</window>
`

var strExampleLayoutXml = `<?xml version="1.0" encoding="UTF-8"?>
<layout>
  <text />
  <text />
  <text />
</layout>
`

func ExampleLoader() {
        v, e := LoadFile("example.xml")
        if e != nil {
                log.Fatalf("LoadFile: failed loading example.xml: %v", e)
                return
        }
        if v == nil {
                log.Fatalf("LoadFile: view is nil")
                return
        }

        v.Set(ShowAll, ValueOf(true))
        v.Connect(OnDestroy, Quit)

        Interact()
}

func TestExamples(t *testing.T) {
        ExampleLoader()
}
