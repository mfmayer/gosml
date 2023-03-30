package gosml

type GetListRequest struct {
	ClientID OctetString
	ServerID OctetString // optional
	Username OctetString // optional
	Password OctetString // optional
	ListName OctetString // optional
}

func GetListRequestParse(buf *Buffer) (GetListRequest, error) {
	msg := GetListRequest{}
	var err error

	if err := buf.Expect(OCTET_TYPE_LIST, 5); err != nil {
		return msg, err
	}

	if msg.ClientID, err = buf.OctetStringParse(); err != nil {
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

	if msg.ListName, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	return msg, nil
}
