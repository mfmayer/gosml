package gosml

import (
	"encoding/binary"
	"fmt"
)

const (
	TYPE_NUMBER_8  = 1
	TYPE_NUMBER_16 = 2
	TYPE_NUMBER_32 = 4
	TYPE_NUMBER_64 = 8
)

func (buf *Buffer) BooleanParse() (bool, error) {
	if buf.OptionalIsSkipped() {
		return false, nil
	}

	if err := buf.Expect(OCTET_TYPE_BOOLEAN, 1); err != nil {
		return false, err
	}

	b := buf.GetCurrentByte()
	buf.UpdateBytesRead(1)
	return b > 0, nil
}

func (buf *Buffer) U8Parse() (uint8, error) {
	num, err := buf.NumberParse(OCTET_TYPE_UNSIGNED, TYPE_NUMBER_8)
	return uint8(num), err
}

func (buf *Buffer) U16Parse() (uint16, error) {
	num, err := buf.NumberParse(OCTET_TYPE_UNSIGNED, TYPE_NUMBER_16)
	return uint16(num), err
}

func (buf *Buffer) U32Parse() (uint32, error) {
	num, err := buf.NumberParse(OCTET_TYPE_UNSIGNED, TYPE_NUMBER_32)
	return uint32(num), err
}

func (buf *Buffer) U64Parse() (uint64, error) {
	num, err := buf.NumberParse(OCTET_TYPE_UNSIGNED, TYPE_NUMBER_64)
	return uint64(num), err
}

func (buf *Buffer) I8Parse() (int8, error) {
	num, err := buf.NumberParse(OCTET_TYPE_INTEGER, TYPE_NUMBER_8)
	return int8(num), err
}

func (buf *Buffer) I16Parse() (int16, error) {
	num, err := buf.NumberParse(OCTET_TYPE_INTEGER, TYPE_NUMBER_16)
	return int16(num), err
}

func (buf *Buffer) I32Parse() (int32, error) {
	num, err := buf.NumberParse(OCTET_TYPE_INTEGER, TYPE_NUMBER_32)
	return int32(num), err
}

func (buf *Buffer) I64Parse() (int64, error) {
	num, err := buf.NumberParse(OCTET_TYPE_INTEGER, TYPE_NUMBER_64)
	return int64(num), err
}

func (buf *Buffer) NumberParse(numType uint8, maxSize int) (int64, error) {
	if skip := buf.OptionalIsSkipped(); skip {
		return 0, nil
	}

	typeField := buf.GetNextType()
	if typeField != numType {
		return 0, fmt.Errorf("unexpected type %02x (expected %02x)", typeField, numType)
	}

	length := buf.GetNextLength()
	if length < 0 || length > maxSize {
		return 0, fmt.Errorf("invalid length: %d", length)
	}

	np := make([]byte, maxSize)
	missingBytes := maxSize - length

	for i := 0; i < length; i++ {
		np[missingBytes+i] = buf.Bytes[buf.Cursor+i]
	}

	negativeInt := typeField == OCTET_TYPE_INTEGER && (typeField&128 > 0)
	if negativeInt {
		for i := 0; i < missingBytes; i++ {
			np[i] = 0xFF
		}
	}

	var num int64
	switch maxSize {
	case TYPE_NUMBER_8:
		num = int64(int8(np[0]))
	case TYPE_NUMBER_16:
		num = int64(int16(binary.BigEndian.Uint16(np)))
	case TYPE_NUMBER_32:
		num = int64(int32(binary.BigEndian.Uint32(np)))
	case TYPE_NUMBER_64:
		num = int64(binary.BigEndian.Uint64(np))
	default:
		return num, fmt.Errorf("invalid number type size %02x", maxSize)
	}

	buf.UpdateBytesRead(length)

	return num, nil
}

func (buf *Buffer) OctetStringParse() (OctetString, error) {
	if skip := buf.OptionalIsSkipped(); skip {
		return nil, nil
	}

	if err := buf.ExpectType(OCTET_TYPE_OCTET_STRING); err != nil {
		return nil, err
	}

	length := buf.GetNextLength()
	if length < 0 {
		return nil, fmt.Errorf("invalid octet string length %d", length)
	}

	str := buf.Bytes[buf.Cursor : buf.Cursor+length]
	buf.UpdateBytesRead(length)

	return str, nil
}

func (buf *Buffer) StatusParse() (int64, error) {
	/*
		if (BufOptionalIsSkipped(buf)) {
			return 0;
		}

		int max = 1;
		int type = BufGetNextType(buf);
		unsigned char byte = BufGetCurrentByte(buf);

		Status *status = StatusInit();
		status->type = type;
		switch (type) {
			case TYPEUNSIGNED:
				// get maximal size, if not all bytes are used (example: only 6 bytes for a u64)
				while (max < ((byte & LENGTHFIELD) - 1)) {
					max <<= 1;
				}

				status->data.status8 = NumberParse(buf, type, max);
				status->type |= max;
				break;
			default:
				buf->error = 1;
				break;
		}
	*/
	// TODO proper type handling

	if skip := buf.OptionalIsSkipped(); skip {
		return 0, nil
	}

	buf.Debug()

	var max uint8 = 1
	var status8 int64
	typeField := buf.GetNextType()
	statusType := typeField
	b := buf.GetCurrentByte()

	if typeField == OCTET_TYPE_UNSIGNED {
		// get maximal size, if not all bytes are used (example: only 6 bytes for a u64)
		for max < ((b & OCTET_LENGTH_FIELD) - 1) {
			max = max << 1
		}

		if _, err := buf.NumberParse(typeField, int(max)); err != nil {
			return 0, err
		}

		statusType = statusType | max
	} else {
		return 0, fmt.Errorf("unexpected type %02x (expected %02x)", typeField, OCTET_TYPE_UNSIGNED)
	}

	return status8, nil
}

func (buf *Buffer) TimeParse() (Time, error) {
	/*
		if (BufOptionalIsSkipped(buf)) {
			return 0;
		}

		Time *tme = TimeInit();

		if (BufGetNextType(buf) != TYPELIST) {
			buf->error = 1;
			goto error;
		}

		if (BufGetNextLength(buf) != 2) {
			buf->error = 1;
			goto error;
		}

		tme->tag = U8Parse(buf);
		if (BufHasErrors(buf)) goto error;

		int type = BufGetNextType(buf);
		switch (type) {
		case TYPEUNSIGNED:
			tme->data.timestamp = U32Parse(buf);
			if (BufHasErrors(buf)) goto error;
			break;
		case TYPELIST:
			// Some meters (e.g. FROETEC Multiflex ZG22) giving not one uint32
			// as timestamp, but a list of 3 values.
			// Ignoring these values, so that parsing does not fail.
			BufGetNextLength(buf); // should we check the length here?
			u32 *t1 = U32Parse(buf);
			if (BufHasErrors(buf)) goto error;
			i16 *t2 = I16Parse(buf);
			if (BufHasErrors(buf)) goto error;
			i16 *t3 = I16Parse(buf);
			if (BufHasErrors(buf)) goto error;
			fprintf(stderr,
				"libsml: error: Time as list[3]: ignoring value[0]=%u value[1]=%d value[2]=%d\n",
				*t1, *t2, *t3);
			break;
		default:
			goto error;
		}
	*/
	// TODO return proper timestamps

	if skip := buf.OptionalIsSkipped(); skip {
		return 0, nil
	}

	if err := buf.Expect(OCTET_TYPE_LIST, 2); err != nil {
		return 0, err
	}

	// time.tag
	if _, err := buf.U8Parse(); err != nil {
		return 0, err
	}

	var timestamp uint32
	var err error

	typeField := buf.GetNextType()
	switch typeField {
	case OCTET_TYPE_UNSIGNED:
		if timestamp, err = buf.U32Parse(); err != nil {
			return 0, err
		}
	case OCTET_TYPE_LIST:
		// Some meters (e.g. FROETEC Multiflex ZG22) giving not one uint32
		// as timestamp, but a list of 3 values.
		// Ignoring these values, so that parsing does not fail.
		buf.GetNextLength() // should we check the length here?

		if _, err := buf.U32Parse(); err != nil {
			return 0, err
		}
		if _, err := buf.I16Parse(); err != nil {
			return 0, err
		}
		if _, err := buf.I16Parse(); err != nil {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("invalid time format %02x", typeField)
	}

	return Time(timestamp), nil
}

func (buf *Buffer) ValueParse() (Value, error) {
	/*
		if (BufOptionalIsSkipped(buf)) {
			return 0;
		}

		int max = 1;
		int type = BufGetNextType(buf);
		unsigned char byte = BufGetCurrentByte(buf);

		Value *value = ValueInit();
		value->type = type;

		switch (type) {
			case TYPEOCTETSTRING:
				value->data.bytes = OctetStringParse(buf);
				break;
			case TYPEBOOLEAN:
				value->data.boolean = BooleanParse(buf);
				break;
			case TYPEUNSIGNED:
			case TYPEINTEGER:
				// get maximal size, if not all bytes are used (example: only 6 bytes for a u64)
				while (max < ((byte & LENGTHFIELD) - 1)) {
					max <<= 1;
				}

				value->data.uint8 = NumberParse(buf, type, max);
				value->type |= max;
				break;
			default:
				buf->error = 1;
				break;
		}
	*/
	value := Value{}

	if buf.OptionalIsSkipped() {
		return value, nil
	}

	typeField := buf.GetNextType()
	b := buf.GetCurrentByte()

	max := 1
	value.Typ = typeField

	var err error
	switch typeField {
	case OCTET_TYPE_OCTET_STRING:
		value.DataBytes, err = buf.OctetStringParse()
		if err != nil {
			return value, err
		}
	case OCTET_TYPE_BOOLEAN:
		value.DataBoolean, err = buf.BooleanParse()
		if err != nil {
			return value, err
		}
	case OCTET_TYPE_UNSIGNED:
		// get maximal size, if not all bytes are used (example: only 6 bytes for a u64)
		for max < int((b&OCTET_LENGTH_FIELD)-1) {
			max = max << 1
		}

		value.DataInt, err = buf.NumberParse(typeField, max)
		if err != nil {
			return value, err
		}

		value.Typ = value.Typ | uint8(max)
	case OCTET_TYPE_INTEGER:
		// get maximal size, if not all bytes are used (example: only 6 bytes for a u64)
		for max < int((b&OCTET_LENGTH_FIELD)-1) {
			max = max << 1
		}

		value.DataInt, err = buf.NumberParse(typeField, max)
		if err != nil {
			return value, err
		}

		value.Typ = value.Typ | uint8(max)
	default:
		return value, fmt.Errorf("unexpected type %02x", typeField)
	}

	return value, nil
}
