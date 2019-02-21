package collection

import (
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestIteratorGetNextBytes(t *testing.T) {
	stub := shim.NewMockStub("bla", &cc{})
	stub.MockTransactionStart("a")
	c := New([]string{"data"}, nil)
	if err := c.Accessor(stub).PutBytes([]string{"k"}, []byte("blah!")); err != nil {
		t.Fatal(err)
	}
	if err := c.Accessor(stub).PutBytes([]string{"j"}, []byte("bleh!")); err != nil {
		t.Fatal(err)
	}
	stub.MockTransactionEnd("a")
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{"basic/1", []byte("bleh!"), false},
		{"basic/2", []byte("blah!"), false},
	}
	i := &iter{
		stub:      stub,
		namespace: []string{"data"},
		state:     nil,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.NextBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("iter.NextBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("iter.NextBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIteratorToByteArrays(t *testing.T) {
	stub := shim.NewMockStub("bla", &cc{})
	stub.MockTransactionStart("a")
	c := New([]string{"data"}, nil)
	if err := c.Accessor(stub).PutBytes([]string{"k"}, []byte("blah!")); err != nil {
		t.Fatal(err)
	}
	if err := c.Accessor(stub).PutBytes([]string{"j"}, []byte("bleh!")); err != nil {
		t.Fatal(err)
	}
	stub.MockTransactionEnd("a")
	it := &iter{
		stub: stub,
	}
	bs, err := it.ToByteArrays()
	if err != nil {
		t.Fatal(err)
	}
	want := [][]byte{[]byte("bleh!"), []byte("blah!")}
	if !reflect.DeepEqual(bs, want) {
		t.Errorf("iter.ToByteArrays() = %v, want %v", bs, want)
	}
}

func TestIteratorClose(t *testing.T) {
	stub := shim.NewMockStub("bla", &cc{})
	istub := shim.NewMockStateRangeQueryIterator(stub, "a", "b")
	it := &iter{
		state: istub,
	}
	err := it.Close()
	if err != nil {
		t.Errorf("failed closing: %v", err)
	}
	if !istub.Closed {
		t.Error("underlying iterator was not close")
	}
}
