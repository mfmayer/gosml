package gosml

type CloseResponse CloseRequest

func CloseResponseParse(buf *Buffer) (CloseResponse, error) {
	msg := CloseResponse{}
	var err error

	if err := buf.Expect(OCTET_TYPE_LIST, 1); err != nil {
		return msg, err
	}

	if msg.GlobalSignature, err = buf.OctetStringParse(); err != nil {
		return msg, err
	}

	return msg, nil
}
