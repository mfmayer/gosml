package gosml

import (
	"errors"
	"fmt"
)

const (
	MESSAGE_OPEN_REQUEST                = 0x00000100
	MESSAGE_OPEN_RESPONSE               = 0x00000101
	MESSAGE_CLOSE_REQUEST               = 0x00000200
	MESSAGE_CLOSE_RESPONSE              = 0x00000201
	MESSAGE_GET_PROFILE_PACK_REQUEST    = 0x00000300
	MESSAGE_GET_PROFILE_PACK_RESPONSE   = 0x00000301
	MESSAGE_GET_PROFILE_LIST_REQUEST    = 0x00000400
	MESSAGE_GET_PROFILE_LIST_RESPONSE   = 0x00000401
	MESSAGE_GET_PROC_PARAMETER_REQUEST  = 0x00000500
	MESSAGE_GET_PROC_PARAMETER_RESPONSE = 0x00000501
	MESSAGE_SET_PROC_PARAMETER_REQUEST  = 0x00000600
	MESSAGE_SET_PROC_PARAMETER_RESPONSE = 0x00000601 // This doesn't exist in the spec
	MESSAGE_GET_LIST_REQUEST            = 0x00000700
	MESSAGE_GET_LIST_RESPONSE           = 0x00000701
	MESSAGE_ATTENTION_RESPONSE          = 0x0000FF01
)

type Message struct {
	TransactionID OctetString
	GroupID       uint8
	AbortOnError  uint8
	MessageBody   MessageBody
	Crc           uint16
}

type MessageBody struct {
	Tag  uint32
	Data MessageBodyData
}

type MessageBodyData interface{}

func MessageParse(buf *Buffer, validate ...bool) (*Message, error) {
	// debug(buf, "MessageParse")

	msg := &Message{}
	var err error

	crcStart := buf.Cursor

	if err := buf.Expect(OCTET_TYPE_LIST, 6); err != nil {
		return msg, err
	}

	if msg.TransactionID, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	if msg.GroupID, err = buf.U8Parse(); err != nil {
		return msg, err
	}

	if msg.AbortOnError, err = buf.U8Parse(); err != nil {
		return msg, err
	}

	if msg.MessageBody, err = MessageBodyParse(buf); err != nil {
		return msg, err
	}

	crcEnd := buf.Cursor

	if msg.Crc, err = buf.U16Parse(); err != nil {
		return msg, err
	}

	if len(validate) > 0 && validate[0] {
		//		fmt.Println(buf.Cursor)
		crc := crc16Calculate(buf.Bytes[crcStart:crcEnd], crcEnd-crcStart)
		//		fmt.Printf("%04x-%04x\n", crc, msg.Crc)

		if crc != msg.Crc {
			err := errors.New("crc error")
			return msg, err
		}
	}

	if buf.GetCurrentByte() == OCTET_MESSAGE_END {
		buf.UpdateBytesRead(1)
	}

	return msg, nil
}

func MessageBodyParse(buf *Buffer) (MessageBody, error) {
	body := MessageBody{}
	var err error

	if err := buf.Expect(OCTET_TYPE_LIST, 2); err != nil {
		return body, err
	}

	if body.Tag, err = buf.U32Parse(); err != nil {
		return body, err
	}

	switch body.Tag {
	case MESSAGE_OPEN_REQUEST:
		body.Data, err = OpenRequestParse(buf)
		return body, err
	case MESSAGE_OPEN_RESPONSE:
		body.Data, err = OpenResponseParse(buf)
		return body, err
	case MESSAGE_CLOSE_REQUEST:
		body.Data, err = CloseRequestParse(buf)
		return body, err
	case MESSAGE_CLOSE_RESPONSE:
		body.Data, err = CloseResponseParse(buf)
		return body, err
	case MESSAGE_GET_PROFILE_PACK_REQUEST:
		return body, fmt.Errorf("unimplemented message type MESSAGE_GET_PROFILE_PACK_REQUEST")
		// msgBody->data = GetProfilePackRequestParse(buf);
	case MESSAGE_GET_PROFILE_PACK_RESPONSE:
		return body, fmt.Errorf("unimplemented message type MESSAGE_GET_PROFILE_PACK_RESPONSE")
		// msgBody->data = GetProfilePackResponseParse(buf);
	case MESSAGE_GET_PROFILE_LIST_REQUEST:
		return body, fmt.Errorf("unimplemented message type MESSAGE_GET_PROFILE_LIST_REQUEST")
		// msgBody->data = GetProfileListRequestParse(buf);
	case MESSAGE_GET_PROFILE_LIST_RESPONSE:
		return body, fmt.Errorf("unimplemented message type MESSAGE_GET_PROFILE_LIST_RESPONSE")
		// msgBody->data = GetProfileListResponseParse(buf);
	case MESSAGE_GET_PROC_PARAMETER_REQUEST:
		return body, fmt.Errorf("unimplemented message type MESSAGE_GET_PROC_PARAMETER_REQUEST")
		// msgBody->data = GetProcParameterRequestParse(buf);
	case MESSAGE_GET_PROC_PARAMETER_RESPONSE:
		return body, fmt.Errorf("unimplemented message type MESSAGE_GET_PROC_PARAMETER_RESPONSE")
		// msgBody->data = GetProcParameterResponseParse(buf);
	case MESSAGE_SET_PROC_PARAMETER_REQUEST:
		return body, fmt.Errorf("unimplemented message type MESSAGE_SET_PROC_PARAMETER_REQUEST")
		// msgBody->data = SetProcParameterRequestParse(buf);
	case MESSAGE_GET_LIST_REQUEST:
		body.Data, err = GetListRequestParse(buf)
		return body, err
	case MESSAGE_GET_LIST_RESPONSE:
		body.Data, err = GetListResponseParse(buf)
		return body, err
	case MESSAGE_ATTENTION_RESPONSE:
		return body, fmt.Errorf("unimplemented message type MESSAGE_ATTENTION_RESPONSE")
		// msgBody->data = AttentionResponseParse(buf);
	}

	return body, fmt.Errorf("invalid message type: % x", body.Tag)
}
