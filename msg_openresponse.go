package gosml

type OpenResponse struct {
	Codepage  OctetString
	ClientID  OctetString
	ReqFileID OctetString
	ServerID  OctetString
	RefTime   Time
	Version   uint8
}

func OpenResponseParse(buf *Buffer) (OpenResponse, error) {
	msg := OpenResponse{}
	var err error

	if err := buf.Expect(OCTET_TYPE_LIST, 6); err != nil {
		return msg, err
	}

	if msg.Codepage, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	if msg.ClientID, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	if msg.ReqFileID, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	if msg.ServerID, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	if msg.RefTime, err = buf.TimeParse(); err != nil {
		return msg, err
	}

	if msg.Version, err = buf.U8Parse(); err != nil {
		return msg, err
	}

	return msg, nil
}
