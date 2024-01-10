package route

import (
	"io"
	"log"
	"testing"

	"github.com/pandodao/mtg/mtgpack"
)

func TestPathEncode(t *testing.T) {
	v := Paths{
		{
			Weight: 50,
			Pairs:  []uint16{1},
		},
		{
			Weight: 50,
			Pairs:  []uint16{1},
		},
	}

	enc := mtgpack.NewEncoder()
	if err := enc.EncodeValue(v); err != nil {
		t.Fatal(err)
	}

	data := enc.Bytes()
	t.Log(data)

	dec := mtgpack.NewDecoder(data)
	dec.Reader = &wrappedReader{r: dec.Reader}

	var paths Paths
	if err := dec.DecodeValue(&paths); err != nil {
		t.Fatal(err)
	}

	if want, got := v.String(), paths.String(); want != got {
		t.Errorf("expect %s but got %s", want, got)
	}
}

type wrappedReader struct {
	r io.Reader
}

func (r *wrappedReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	log.Println("read:", p)
	return n, err
}
