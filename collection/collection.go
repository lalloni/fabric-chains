package collection

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
)

const objectType = "_COLITE"

type Collection interface {
	Accessor(shim.ChaincodeStubInterface) Accessor
}

type Accessor interface {
	Get(key []string, v interface{}) (bool, error)
	GetBytes(key []string) ([]byte, error)
	Put(key []string, v interface{}) error
	PutBytes(key []string, bs []byte) error
	Remove(key []string) error
	Iterator() (Iterator, error)
}

func New(key []string, m Marshaller) Collection {
	return &coll{
		namespace:  key,
		marshaller: m,
	}
}

type coll struct {
	namespace  []string
	marshaller Marshaller
}

var _ Collection = &coll{}

func (c *coll) Accessor(stub shim.ChaincodeStubInterface) Accessor {
	return &acc{
		stub:       stub,
		namespace:  c.namespace,
		marshaller: c.marshaller,
	}
}

type acc struct {
	stub       shim.ChaincodeStubInterface
	namespace  []string
	marshaller Marshaller
}

var _ Accessor = &acc{}

func (a *acc) elementKey(key []string) (string, error) {
	if len(key) < 1 {
		return "", errors.New("element key can not be empty")
	}
	k, err := a.stub.CreateCompositeKey(objectType, append(a.namespace, key...))
	if err != nil {
		return "", errors.Wrapf(err, "composing key from %v", key)
	}
	return k, nil
}

func (a *acc) GetBytes(key []string) ([]byte, error) {
	k, err := a.elementKey(key)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	bs, err := a.stub.GetState(k)
	if err != nil {
		return nil, errors.Wrapf(err, "getting state for %v", key)
	}
	return bs, nil
}

func (a *acc) Get(key []string, v interface{}) (bool, error) {
	bs, err := a.GetBytes(key)
	if err != nil {
		return false, errors.Wrap(err, "getting element bytes")
	}
	if bs == nil {
		return false, nil
	}
	if a.marshaller == nil {
		return false, errors.New("marshaller not set")
	}
	err = a.marshaller.Unmarshal(bs, v)
	if err != nil {
		return false, errors.Wrap(err, "unmarshalling element")
	}
	return true, nil
}

func (a *acc) PutBytes(key []string, bs []byte) error {
	k, err := a.elementKey(key)
	if err != nil {
		return errors.WithStack(err)
	}
	err = a.stub.PutState(k, bs)
	if err != nil {
		return errors.Wrapf(err, "putting state for %v", key)
	}
	return nil
}

func (a *acc) Put(key []string, v interface{}) error {
	if a.marshaller == nil {
		return errors.New("marshaller not set")
	}
	bs, err := a.marshaller.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "marshalling element")
	}
	err = a.PutBytes(key, bs)
	if err != nil {
		return errors.Wrap(err, "putting element bytes")
	}
	return nil
}

func (a *acc) Remove(key []string) error {
	k, err := a.elementKey(key)
	if err != nil {
		return errors.WithStack(err)
	}
	err = a.stub.DelState(k)
	if err != nil {
		return errors.Wrapf(err, "removing element %v", key)
	}
	return nil
}

func (a *acc) Iterator() (Iterator, error) {
	return &iter{
		stub:      a.stub,
		namespace: a.namespace,
	}, nil
}
