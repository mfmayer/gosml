package gosml

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const (
	maxFileSize = 512
)

var (
	escSeq = []byte{0x1b, 0x1b, 0x1b, 0x1b}
	endSeq = []byte{0x1b, 0x1b, 0x1b, 0x1b, 0x1a}
)

// ErrUnrecognizedSequence means that a sequence was found but its end was not found.
// E.g. a new sequence started before the end of current sequence was found.
var ErrUnrecognizedSequence = errors.New("unrecognized sequence")

// ErrSequenceTooLong means that the max length of a sequence has been reached before
// end of sequence has been detected.
var ErrSequenceTooLong = errors.New("max sequence length exceeded")

type OctetString []byte

type Time uint32

type Value struct {
	Typ         uint8
	DataBytes   OctetString
	DataBoolean bool
	DataInt     int64
}

func readChunk(r *bufio.Reader, buf []byte) error {
	bytes, err := r.Read(buf)
	if err != nil {
		return err
	}

	if bytes < len(buf) {
		// fmt.Printf("ReadChunk %d -> %d\n", len(buf), bytes)
		return errors.New("premature eof")
	}

	// success - no error
	return nil
}

// readFile reads from buffered reader until next SML file has been completely read and returns
// full SML file as byte slice which can then be parsed with FileParse to get its messages
func readFile(r *bufio.Reader) ([]byte, error) {
	buf := make([]byte, maxFileSize)

	var len int
	var err error

	// find escape sequence/begin 1B 1B 1B 1B 01 01 01 01
	for len < 8 {
		if buf[len], err = r.ReadByte(); err != nil {
			return nil, err
		}

		if (buf[len] == 0x1b && len < 4) || (buf[len] == 0x01 && len >= 4) {
			len++
		} else {
			len = 0
		}
	}

	// found start sequence
	for len+8 < maxFileSize {
		if err = readChunk(r, buf[len:len+4]); err != nil {
			return nil, err
		}

		// find escape sequence
		if bytes.Equal(buf[len:len+4], escSeq) {
			len += 4

			// read end sequence
			if err = readChunk(r, buf[len:len+4]); err != nil {
				return nil, err
			}

			if buf[len] == 0x1a {
				// found end sequence
				len += 4
				return buf[:len], nil
			}

			// don't read other escaped sequences yet
			return nil, ErrUnrecognizedSequence
		}

		// continue reading
		len += 4
	}

	return nil, ErrSequenceTooLong
}

// parseFile parses SML file provided as byte slice
func parseFile(fileBytes []byte) ([]*Message, error) {
	buf := &Buffer{
		Bytes:  fileBytes,
		Cursor: 0,
	}

	messages := make([]*Message, 0)

	for buf.Cursor < len(buf.Bytes) {
		if buf.GetCurrentByte() == OCTET_MESSAGE_END {
			// reading trailing zeroed bytes
			buf.UpdateBytesRead(1)
			continue
		}

		msg, err := MessageParse(buf, true)
		if err != nil {
			return messages, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

type obisGroupCallback struct {
	callbacks   []func(message *ListEntry)
	childGroups map[byte]*obisGroupCallback
}

func newObisGroupCallback() *obisGroupCallback {
	return &obisGroupCallback{
		childGroups: map[byte]*obisGroupCallback{},
	}
}

func (oc *obisGroupCallback) addCallback(subCode OctetString, callback func(message *ListEntry)) {
	if len(subCode) == 0 {
		oc.callbacks = append(oc.callbacks, callback)
	} else {
		subGroupCallback, ok := oc.childGroups[subCode[0]]
		if !ok {
			subGroupCallback = newObisGroupCallback()
			oc.childGroups[subCode[0]] = subGroupCallback
		}
		subGroupCallback.addCallback(subCode[1:], callback)
	}
}

func (oc *obisGroupCallback) call(obisCode OctetString, listEntry *ListEntry) {
	// call registered callbacks
	for _, callback := range oc.callbacks {
		callback(listEntry)
	}
	// check if additional registered handlers exist for remaining obis groups
	subOc, ok := oc.childGroups[obisCode[0]]
	if ok {
		subOc.call(obisCode[1:], listEntry)
	}
}

type options struct {
	topLevelCallback *obisGroupCallback
}

type ReadOption func(*options)

func WithObisCallback(obisCode OctetString, callback func(message *ListEntry)) ReadOption {
	return func(o *options) {
		if o.topLevelCallback == nil {
			o.topLevelCallback = newObisGroupCallback()
		}
		o.topLevelCallback.addCallback(obisCode, callback)
	}
}

// Read reads and parses sml file from given buffered reader.
// If sml file is not recognized ErrUnrecognizedSequence is returned.
// If sml file is too long ErrSequenceTooLong is returned.
// If file is successfully read and parsed slice of found messages is returned
func Read(r *bufio.Reader, opts ...ReadOption) error {
	options := &options{}
	for _, opt := range opts {
		opt(options)
	}
loop:
	for {
		var fileBytes []byte
		fileBytes, err := readFile(r)
		switch {
		case err == io.EOF:
			break loop
		case err == ErrSequenceTooLong || err == ErrUnrecognizedSequence:
			continue
		case err != nil:
			return err
		}
		// parse without escaped begin and end sequences
		fileMessages, err := parseFile(fileBytes[8 : len(fileBytes)-8])
		if err != nil {
			return err
		}
		// handle messages according to given options
		for _, msg := range fileMessages {
			if options.topLevelCallback != nil {
				if msg.MessageBody.Tag == MESSAGE_GET_LIST_RESPONSE {
					list, ok := msg.MessageBody.Data.(GetListResponse)
					if !ok {
						continue
					}
					for _, elem := range list.ValList {
						if len(elem.ObjName) > 0 {
							options.topLevelCallback.call(elem.ObjName, elem)
						}
						// fmt.Printf("%v\n", elem)
					}
				}
			}
		}
	}
	return nil
}
