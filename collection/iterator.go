package collection

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
)

type Iterator interface {
	HasNext() (bool, error)
	Next(v interface{}) (bool, error)
	NextBytes() ([]byte, error)
	ToArray(vs []interface{}) error
	ToByteArrays() ([][]byte, error)
	Close() error
}

type iter struct {
	stub      shim.ChaincodeStubInterface
	namespace []string
	state     shim.StateQueryIteratorInterface
}

var _ Iterator = &iter{}

func (i *iter) NextBytes() ([]byte, error) {
	if i.state == nil {
		err := i.fetch()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	bs, err := i.nextBytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return bs, nil
}

func (i *iter) Next(v interface{}) (bool, error) {
	return false, errors.New("not implemented")
}

func (i *iter) HasNext() (bool, error) {
	if i.state == nil {
		err := i.fetch()
		if err != nil {
			return false, errors.WithStack(err)
		}
	}
	return i.state.HasNext(), nil
}

func (i *iter) ToArray(vs []interface{}) error {
	return errors.New("not implemented")
}

func (i *iter) ToByteArrays() ([][]byte, error) {
	rs := [][]byte{}
	for {
		hn, err := i.HasNext()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if !hn {
			break
		}
		bs, err := i.NextBytes()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		rs = append(rs, bs)
	}
	return rs, nil
}

func (i *iter) Close() error {
	if i.state == nil {
		return nil
	}
	err := i.state.Close()
	if err != nil {
		return errors.Wrap(err, "closing iterator")
	}
	return nil
}

func (i *iter) fetch() error {
	state, err := i.stub.GetStateByPartialCompositeKey(objectType, i.namespace)
	if err != nil {
		return errors.Wrap(err, "fetching collection items")
	}
	i.state = state
	return nil
}

func (i *iter) nextBytes() ([]byte, error) {
	kv, err := i.state.Next()
	if err != nil {
		return nil, errors.Wrap(err, "getting next item")
	}
	return kv.GetValue(), nil
}
