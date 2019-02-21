package collection

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Marshaller interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(bs []byte, v interface{}) error
}

func JSONMarshaller() Marshaller {
	return &jsonMarshaller{}
}

type jsonMarshaller struct{}

var _ Marshaller = &jsonMarshaller{}

func (m *jsonMarshaller) Marshal(v interface{}) (bs []byte, err error) {
	if mr, ok := v.(json.Marshaler); ok {
		bs, err = mr.MarshalJSON()
	} else {
		bs, err = json.Marshal(v)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return bs, nil
}

func (m *jsonMarshaller) Unmarshal(bs []byte, v interface{}) (err error) {
	if mr, ok := v.(json.Unmarshaler); ok {
		err = mr.UnmarshalJSON(bs)
	} else {
		err = json.Unmarshal(bs, v)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
