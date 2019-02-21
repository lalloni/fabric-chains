package collection

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type cc struct {
}

func (c *cc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return peer.Response{}
}

func (c *cc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return peer.Response{}
}

var _ shim.Chaincode = &cc{}

func TestAccessorGetBytes(t *testing.T) {
	stub := shim.NewMockStub("bla", &cc{})
	k, err := stub.CreateCompositeKey(objectType, []string{"data", "k"})
	if err != nil {
		t.Fatal(err)
	}
	stub.MockTransactionStart("a")
	stub.PutState(k, []byte("123456"))
	stub.MockTransactionEnd("a")
	tests := []struct {
		name    string
		key     []string
		want    []byte
		wantErr bool
	}{
		{"basic", []string{"k"}, []byte("123456"), false},
		{"emptydata", []string{"j"}, nil, false},
		{"emptykey", nil, nil, true},
	}
	a := &acc{
		stub:      stub,
		namespace: []string{"data"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.GetBytes(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("acc.GetBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("acc.GetBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Data struct {
	Name string
	Age  int
}

func TestAccessorGet(t *testing.T) {
	d := Data{"pedro", 20}
	ds, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	stub := shim.NewMockStub("bla", &cc{})
	k, err := stub.CreateCompositeKey(objectType, []string{"data", "k"})
	if err != nil {
		t.Fatal(err)
	}
	stub.MockTransactionStart("a")
	stub.PutState(k, ds)
	stub.MockTransactionEnd("a")
	type args struct {
		key []string
		v   interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     bool
		wantData Data
		wantErr  bool
	}{
		{"basic", args{[]string{"k"}, &Data{"blah!", -10}}, true, d, false},
		{"emptydata", args{[]string{"j"}, &Data{}}, false, Data{}, false},
	}
	a := &acc{
		stub:       stub,
		namespace:  []string{"data"},
		marshaller: JSONMarshaller(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.Get(tt.args.key, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("acc.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("acc.Get() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(&tt.wantData, tt.args.v) {
				t.Errorf("acc.Get() v = %v, want %v", tt.args.v, tt.wantData)
			}
		})
	}
}

func TestAccessorPutBytes(t *testing.T) {
	stub := shim.NewMockStub("bla", &cc{})
	k, err := stub.CreateCompositeKey(objectType, []string{"data", "k"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		key []string
		bs  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple", args{[]string{"k"}, []byte("blabla")}, false},
		{"emptykey", args{nil, []byte("blabla")}, true},
		{"emptydata", args{[]string{"k"}, nil}, false},
	}
	a := &acc{
		stub:      stub,
		namespace: []string{"data"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub.MockTransactionStart("a")
			if err := a.PutBytes(tt.args.key, tt.args.bs); (err != nil) != tt.wantErr {
				t.Errorf("acc.PutBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			stub.MockTransactionEnd("a")
			bs, err := stub.GetState(k)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(bs, tt.args.bs) {
				t.Errorf("acc.PutBytes(bs) bs = %v, stateBS = %v", tt.args.bs, bs)
			}
		})
	}
}

func TestAccessorPut(t *testing.T) {
	stub := shim.NewMockStub("bla", &cc{})
	k, err := stub.CreateCompositeKey(objectType, []string{"data", "k"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		key []string
		v   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"basic", args{[]string{"k"}, Data{"pedro", 10}}, false},
		{"emptykey", args{nil, Data{"pedro", 10}}, true},
		{"emptydata", args{[]string{"k"}, nil}, false},
	}
	m := JSONMarshaller()
	a := &acc{
		stub:       stub,
		namespace:  []string{"data"},
		marshaller: m,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub.MockTransactionStart("a")
			if err := a.Put(tt.args.key, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("acc.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
			stub.MockTransactionEnd("a")
			bs, err := stub.GetState(k)
			if err != nil {
				t.Fatal(err)
			}
			bs2, err := m.Marshal(tt.args.v)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(bs, bs2) {
				t.Errorf("acc.Put(v) bs = %v, stateBS = %v", bs2, bs)
			}
		})
	}
}

func TestAccessorIterator(t *testing.T) {
	stub := shim.NewMockStub("bla", &cc{})
	a := &acc{
		stub:      stub,
		namespace: []string{"algo"},
	}
	it, err := a.Iterator()
	if err != nil {
		t.Fatal(err)
	}
	i := it.(*iter)
	if i.stub != a.stub {
		t.Errorf("stub not set")
	}
	if !reflect.DeepEqual(i.namespace, a.namespace) {
		t.Errorf("namespace not set")
	}
}

func TestAccessorRemove(t *testing.T) {
	stub := shim.NewMockStub("bla", &cc{})
	a := &acc{
		stub:      stub,
		namespace: []string{"algo"},
	}
	stub.MockTransactionStart("tx1")
	if err := a.PutBytes([]string{"a"}, []byte("hola1")); err != nil {
		t.Fatal(err)
	}
	if err := a.PutBytes([]string{"b"}, []byte("hola2")); err != nil {
		t.Fatal(err)
	}
	if err := a.PutBytes([]string{"c"}, []byte("hola3")); err != nil {
		t.Fatal(err)
	}
	stub.MockTransactionEnd("")
	stub.MockTransactionStart("tx2")
	if err := a.Remove([]string{"b"}); err != nil {
		t.Fatal(err)
	}
	stub.MockTransactionEnd("")
	it, err := a.Iterator()
	if err != nil {
		t.Fatal(err)
	}
	bs, err := it.ToByteArrays()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(bs, [][]byte{[]byte("hola1"), []byte("hola3")}) {
		t.Errorf("acc.Remove() failed: %v", bs)
	}
}
