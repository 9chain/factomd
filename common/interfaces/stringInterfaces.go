// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

import (
	"bytes"
)

//Interface for printing structures into JSON
type JSONable interface {
	JSONByte() ([]byte, error)
	JSONString() (string, error)
	JSONBuffer(b *bytes.Buffer) error
}

//Interface for both JSON and Spew
type Printable interface {
	JSONable
	String() string
}

//Interface for short, reoccuring data structures to interpret themselves into human-friendly form
type ShortInterpretable interface {
	IsInterpretable() bool //Whether the structure can interpret itself
	Interpret() string     //Turns the data encoded int he structure into human-friendly string
}
