package xz

import (
	"bytes"
	"testing"
)

func TestRecordReadWrite(t *testing.T) {
	r := record{1234567, 10000}
	var buf bytes.Buffer
	n, err := r.writeTo(&buf)
	if err != nil {
		t.Fatalf("writeTo error %s", err)
	}
	var g record
	m, err := g.readFrom(&buf)
	if err != nil {
		t.Fatalf("readFrom error %s", err)
	}
	if m != n {
		t.Fatalf("read %d bytes; wrote %d", m, n)
	}
	if g.unpaddedSize != r.unpaddedSize {
		t.Fatalf("got unpaddedSize %d; want %d", g.unpaddedSize,
			r.unpaddedSize)
	}
	if g.uncompressedSize != r.uncompressedSize {
		t.Fatalf("got uncompressedSize %d; want %d", g.uncompressedSize,
			r.uncompressedSize)
	}
}

func TestIndexReadWrite(t *testing.T) {
	records := []record{{1234, 1}, {2345, 2}}

	var buf bytes.Buffer
	n, err := writeIndex(&buf, records)
	if err != nil {
		t.Fatalf("writeIndex error %s", err)
	}
	if n != buf.Len() {
		t.Fatalf("writeIndex returned %d; want %d", n, buf.Len())
	}

	// indicator
	c, err := buf.ReadByte()
	if err != nil {
		t.Fatalf("buf.ReadByte error %s", err)
	}
	if c != 0 {
		t.Fatalf("indicator %d; want %d", c, 0)
	}

	g, m, err := readIndexBody(&buf)
	if err != nil {
		for i, r := range g {
			t.Logf("records[%d] %v", i, r)
		}
		t.Fatalf("readIndexBody error %s", err)
	}
	if m != n-1 {
		t.Fatalf("readIndexBody returned %d; want %d", m, n-1)
	}
	for i, rec := range records {
		if g[i] != rec {
			t.Errorf("records[%d] is %v; want %v", i, g[i], rec)
		}
	}
}
