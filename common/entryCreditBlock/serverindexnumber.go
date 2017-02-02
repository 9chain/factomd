// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package entryCreditBlock

import (
	"bytes"
	"fmt"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

const (
	ServerIndexNumberSize = 1
)

type ServerIndexNumber struct {
	ServerIndexNumber uint8
}

var _ interfaces.Printable = (*ServerIndexNumber)(nil)
var _ interfaces.BinaryMarshallable = (*ServerIndexNumber)(nil)
var _ interfaces.ShortInterpretable = (*ServerIndexNumber)(nil)
var _ interfaces.IECBlockEntry = (*ServerIndexNumber)(nil)

func (e *ServerIndexNumber) String() string {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberString.Observe(float64(time.Now().UnixNano() - callTime))	
	var out primitives.Buffer
	out.WriteString(fmt.Sprintf(" %-20s\n", "ServerIndexNumber"))
	out.WriteString(fmt.Sprintf("   %-20s %d\n", "Number", e.ServerIndexNumber))
	return (string)(out.DeepCopyBytes())
}

func (e *ServerIndexNumber) Hash() interfaces.IHash {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberHash.Observe(float64(time.Now().UnixNano() - callTime))	
	bin, err := e.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return primitives.Sha(bin)
}

func (e *ServerIndexNumber) GetHash() interfaces.IHash {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberGetHash.Observe(float64(time.Now().UnixNano() - callTime))	
	return e.Hash()
}

func (a *ServerIndexNumber) GetEntryHash() interfaces.IHash {
	return nil
}

func (e *ServerIndexNumber) GetSigHash() interfaces.IHash {
	return nil
}

func (b *ServerIndexNumber) IsInterpretable() bool {
	return true
}

func (b *ServerIndexNumber) Interpret() string {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberInterpret.Observe(float64(time.Now().UnixNano() - callTime))	
	return fmt.Sprintf("ServerIndexNumber %v", b.ServerIndexNumber)
}

func NewServerIndexNumber() *ServerIndexNumber {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberNewServerIndexNumber.Observe(float64(time.Now().UnixNano() - callTime))	
	return new(ServerIndexNumber)
}

func NewServerIndexNumber2(number uint8) *ServerIndexNumber {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberNewServerIndexNumber2.Observe(float64(time.Now().UnixNano() - callTime))	
	sin := new(ServerIndexNumber)
	sin.ServerIndexNumber = number
	return sin
}

func (s *ServerIndexNumber) ECID() byte {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberECID.Observe(float64(time.Now().UnixNano() - callTime))	
	return ECIDServerIndexNumber
}

func (s *ServerIndexNumber) MarshalBinary() ([]byte, error) {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberMarshalBinary.Observe(float64(time.Now().UnixNano() - callTime))	
	buf := new(primitives.Buffer)
	buf.WriteByte(s.ServerIndexNumber)
	return buf.DeepCopyBytes(), nil
}

func (s *ServerIndexNumber) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberUnmarshalBinaryData.Observe(float64(time.Now().UnixNano() - callTime))	
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling ServerIndexNumber: %v", r)
		}
	}()

	buf := primitives.NewBuffer(data)
	var c byte
	if c, err = buf.ReadByte(); err != nil {
		return
	} else {
		s.ServerIndexNumber = c
	}
	newData = buf.DeepCopyBytes()
	return
}

func (s *ServerIndexNumber) UnmarshalBinary(data []byte) (err error) {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberUnmarshalBinary.Observe(float64(time.Now().UnixNano() - callTime))	
	_, err = s.UnmarshalBinaryData(data)
	return
}

func (e *ServerIndexNumber) JSONByte() ([]byte, error) {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberJSONByte.Observe(float64(time.Now().UnixNano() - callTime))	
	return primitives.EncodeJSON(e)
}

func (e *ServerIndexNumber) JSONString() (string, error) {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberJSONString.Observe(float64(time.Now().UnixNano() - callTime))	
	return primitives.EncodeJSONString(e)
}

func (e *ServerIndexNumber) JSONBuffer(b *bytes.Buffer) error {
	callTime := time.Now().UnixNano()
	defer entryCreditBlockServerIndexNumberJSONBuffer.Observe(float64(time.Now().UnixNano() - callTime))	
	return primitives.EncodeJSONToBuffer(e, b)
}

func (e *ServerIndexNumber) GetTimestamp() interfaces.Timestamp {
	return nil
}
