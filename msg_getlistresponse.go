package gosml

import (
	"fmt"
	"math"
)

type GetListResponse struct {
	ClientID       OctetString
	ServerID       OctetString
	ListName       OctetString
	ActSensorTime  Time
	ValList        []*ListEntry
	ListSignature  OctetString
	ActGatewayTime Time
}

type ListEntry struct {
	ObjName        OctetString
	status         int64
	valTime        Time
	Unit           uint8
	scaler         int8
	Value          Value
	ValueSignature OctetString
}

func (le *ListEntry) ObjectName() string {
	return fmt.Sprintf("%d-%d:%d.%d.%d*%d", le.ObjName[0], le.ObjName[1], le.ObjName[2], le.ObjName[3], le.ObjName[4], le.ObjName[5])
}

func (le *ListEntry) Scaler() float64 {
	scaler := int(le.scaler)
	return math.Pow10(scaler)
}

func (le *ListEntry) ValueString() string {
	switch le.Value.Typ {
	case OCTET_TYPE_OCTET_STRING:
		return fmt.Sprintf("% x", le.Value.DataBytes)
	case OCTET_TYPE_BOOLEAN:
		return fmt.Sprintf("%v", le.Value.DataBoolean)
	default:
		if ((le.Value.Typ & OCTET_TYPE_FIELD) == OCTET_TYPE_INTEGER) || ((le.Value.Typ & OCTET_TYPE_FIELD) == OCTET_TYPE_UNSIGNED) {
			value := float64(le.Value.DataInt) * le.Scaler()
			// unit := ""
			// switch entry.Unit {
			// case 0x1B:
			// 	unit = "W"
			// case 0x1E:
			// 	unit = "Wh"
			// }
			return fmt.Sprintf("%12.1f", value)
		}
	}
	return ""
}

func (le *ListEntry) Float() float64 {
	if ((le.Value.Typ & OCTET_TYPE_FIELD) == OCTET_TYPE_INTEGER) || ((le.Value.Typ & OCTET_TYPE_FIELD) == OCTET_TYPE_UNSIGNED) {
		value := float64(le.Value.DataInt) * le.Scaler()
		return value
	}
	return 0.0
}

func (le *ListEntry) String() string {
	return fmt.Sprintf("%-22s%s", le.ObjectName(), le.ValueString())
}

// func PrintListEntry(entry sml.ListEntry) {
// 	fmt.Printf("%-22s", entry)

// 	if entry.Value.Typ == sml.TYPE_OCTET_STRING {
// 		fmt.Printf("% x\n", entry.Value.DataBytes)
// 	} else if entry.Value.Typ == sml.TYPE_BOOLEAN {
// 		fmt.Printf("%v\n", entry.Value.DataBoolean)
// 	} else if ((entry.Value.Typ & sml.TYPE_FIELD) == sml.TYPE_INTEGER) || ((entry.Value.Typ & sml.TYPE_FIELD) == sml.TYPE_UNSIGNED) {
// 		value := ListEntry2Float(entry)

// 		unit := ""
// 		switch entry.Unit {
// 		case 0x1B:
// 			unit = "W"
// 		case 0x1E:
// 			unit = "Wh"
// 		}

// 		fmt.Printf("%12.1f %-3s\n", value, unit)
// 	}
// }

func GetListResponseParse(buf *Buffer) (GetListResponse, error) {
	list := GetListResponse{}
	var err error

	if err := buf.Expect(OCTET_TYPE_LIST, 7); err != nil {
		return list, err
	}

	if list.ClientID, err = buf.OctetStringParse(); err != nil {
		return list, err
	}

	if list.ServerID, err = buf.OctetStringParse(); err != nil {
		return list, err
	}

	if list.ListName, err = buf.OctetStringParse(); err != nil {
		return list, err
	}

	if list.ActSensorTime, err = buf.TimeParse(); err != nil {
		return list, err
	}

	if list.ValList, err = ListParse(buf); err != nil {
		return list, err
	}

	if list.ListSignature, err = buf.OctetStringParse(); err != nil {
		return list, err
	}

	if list.ActGatewayTime, err = buf.TimeParse(); err != nil {
		return list, err
	}

	return list, nil
}

func ListParse(buf *Buffer) ([]*ListEntry, error) {
	if buf.OptionalIsSkipped() {
		return nil, nil
	}

	buf.Debug()

	if err := buf.ExpectType(OCTET_TYPE_LIST); err != nil {
		return nil, err
	}

	list := make([]*ListEntry, 0)

	elems := buf.GetNextLength()

	for elems > 0 {
		elem, err := ListEntryParse(buf)
		if err != nil {
			return nil, err
		}
		list = append(list, elem)
		elems--
	}

	return list, nil
}

func ListEntryParse(buf *Buffer) (*ListEntry, error) {
	buf.Debug()

	elem := ListEntry{}
	var err error

	if err := buf.Expect(OCTET_TYPE_LIST, 7); err != nil {
		return &elem, err
	}

	if elem.ObjName, err = buf.OctetStringParse(); err != nil {
		return &elem, err
	}

	if elem.status, err = buf.StatusParse(); err != nil {
		return &elem, err
	}

	if elem.valTime, err = buf.TimeParse(); err != nil {
		return &elem, err
	}

	if elem.Unit, err = buf.U8Parse(); err != nil {
		return &elem, err
	}

	if elem.scaler, err = buf.I8Parse(); err != nil {
		return &elem, err
	}

	if elem.Value, err = buf.ValueParse(); err != nil {
		return &elem, err
	}

	if elem.ValueSignature, err = buf.OctetStringParse(); err != nil {
		return &elem, err
	}

	return &elem, nil
}
