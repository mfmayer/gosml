package gosml

type CloseRequest struct {
	GlobalSignature OctetString
}

func CloseRequestParse(buf *Buffer) (CloseRequest, error) {
	msg := CloseRequest{}
	var err error

	if err := buf.Expect(OCTET_TYPE_LIST, 1); err != nil {
		return msg, err
	}

	if msg.GlobalSignature, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	return msg, nil
}
