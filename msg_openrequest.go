package gosml

type OpenRequest struct {
	Codepage  OctetString // optional
	ClientID  OctetString
	ReqFileID OctetString
	ServerID  OctetString // optional
	Username  OctetString // optional
	Password  OctetString // optional
	Version   uint8       // optional
}

func OpenRequestParse(buf *Buffer) (OpenRequest, error) {
	msg := OpenRequest{}
	var err error

	if err := buf.Expect(OCTET_TYPE_LIST, 7); err != nil {
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

	if msg.Username, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	if msg.Password, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	if msg.Version, err = buf.U8Parse(); err != nil {
		return msg, err
	}

	return msg, nil
}
