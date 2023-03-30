package gosml

import (
	"fmt"
	"runtime"
)

var debugEnable bool

const (
	OCTET_MESSAGE_END       = 0x00
	OCTET_TYPE_FIELD        = 0x70
	OCTET_LENGTH_FIELD      = 0x0F
	OCTET_ANOTHER_TL        = 0x80
	OCTET_TYPE_OCTET_STRING = 0x00
	OCTET_TYPE_BOOLEAN      = 0x40
	OCTET_TYPE_INTEGER      = 0x50
	OCTET_TYPE_UNSIGNED     = 0x60
	OCTET_TYPE_LIST         = 0x70
	OCTET_OPTIONAL_SKIPPED  = 0x01
)

type Buffer struct {
	Bytes  []byte
	Cursor int
}

func (buf *Buffer) Debug() {
	if debugEnable {
		pc, _, _, ok := runtime.Caller(1)
		funcDetails := runtime.FuncForPC(pc)
		if ok && funcDetails != nil {
			fmt.Printf("%-22s % x\n", funcDetails.Name(), buf.Bytes[buf.Cursor:buf.Cursor+30])
		}
	}
}

func (buf *Buffer) GetCurrentByte() byte {
	return buf.Bytes[buf.Cursor]
}

func (buf *Buffer) UpdateBytesRead(delta int) {
	buf.Cursor += delta
}

func (buf *Buffer) Expect(expectedType uint8, expectedLength int) error {
	if err := buf.ExpectType(expectedType); err != nil {
		return err
	}

	if length := buf.GetNextLength(); length != expectedLength {
		return fmt.Errorf("invalid length: %d (expected %d)", length, expectedLength)
	}

	return nil
}

func (buf *Buffer) ExpectType(expectedType uint8) error {
	if typeField := buf.GetNextType(); typeField != expectedType {
		return fmt.Errorf("unexpected type %02x (expected %02x)", typeField, expectedType)
	}

	return nil
}

func (buf *Buffer) GetNextType() uint8 {
	return buf.GetCurrentByte() & OCTET_TYPE_FIELD
}

func (buf *Buffer) GetNextLength() int {
	var length uint8
	var list int

	b := buf.GetCurrentByte()

	// not a list
	if b&OCTET_TYPE_FIELD != OCTET_TYPE_LIST {
		list = -1
	}

	for {
		b := buf.GetCurrentByte()

		length = length << 4
		length = length | (b & OCTET_LENGTH_FIELD)

		if b&OCTET_ANOTHER_TL != OCTET_ANOTHER_TL {
			break
		}

		// another TL field used
		buf.UpdateBytesRead(1)

		// not a list
		if list != 0 {
			list--
		}
	}

	buf.UpdateBytesRead(1)

	return int(length) + list
}

func (buf *Buffer) OptionalIsSkipped() bool {
	if buf.GetCurrentByte() == OCTET_OPTIONAL_SKIPPED {
		buf.UpdateBytesRead(1)
		return true
	}

	return false
}
