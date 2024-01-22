package buffer

import (
	"io"
	"math/rand"
	"testing"
)

func TestNewBufferFile(t *testing.T) {
	ce := 512
	bf := NewBufferFile()
	cg := cap(bf.Bytes())

	if cg != ce {
		t.Errorf("expected capacity: %v but got: %v", ce, cg)
	}
}

func TestNewBufferFileFromBytes(t *testing.T) {
	l := 1024
	b := make([]byte, l)
	bf := NewBufferFileFromBytes(b)

	lg := len(bf.Bytes())

	if lg != l {
		t.Errorf("expected length: %v but got: %v", l, lg)
	}
}

func TestNewBufferFileCapacity(t *testing.T) {
	ce := 1024
	bf := NewBufferFileCapacity(ce)
	cg := cap(bf.Bytes())

	if cg != ce {
		t.Errorf("expected capacity: %v but got: %v", ce, cg)
	}
}

func TestCreate(t *testing.T) {
	bf := NewBufferFile()

	if _, err := bf.Create("foo"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestOpen(t *testing.T) {
	bf := NewBufferFile()

	if _, err := bf.Open("foo"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestRead(t *testing.T) {
	testCases := []struct {
		name   string // name of test case
		capSrc int    // capacity of source buffer
		capDst int    // capacity of destination count
		cntExp int    // expected copied bytes count
		errExp error  // expected count
	}{
		{
			// regulary read without errors
			name:   "case1",
			capSrc: 4,
			capDst: 5,
			cntExp: 4,
			errExp: nil,
		},
		{
			// read to buffer with not enaugh capacity
			name:   "case2",
			capSrc: 4,
			capDst: 3,
			cntExp: 3,
			errExp: io.EOF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bufSrc := make([]byte, tc.capSrc)
			bufDst := make([]byte, tc.capDst)

			rand.Read(bufSrc) // fill source buffer with random
			bf := NewBufferFileFromBytes(bufSrc)

			cnt, err := bf.Read(bufDst)

			if tc.errExp != err {
				t.Errorf("expected error to be: %v but got: %v", tc.errExp, err)
			}

			if tc.cntExp != cnt {
				t.Errorf("expected copied bytes to be: %v but got: %v", tc.cntExp, cnt)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	testCases := []struct {
		name   string // name of test case
		capSrc int    // capacity of source buffer
		capDst int    // capacity of destination count
		errExp error  // expected count
	}{
		{
			// regulary write with enaugh capacity
			name:   "case1",
			capSrc: 4,
			capDst: 5,
			errExp: nil,
		},
		{
			// write to buffer with not enaugh capacity
			name:   "case2",
			capSrc: 5,
			capDst: 3,
			errExp: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bufSrc := make([]byte, tc.capSrc)

			rand.Read(bufSrc) // fill source buffer with random

			bf := NewBufferFileCapacity(tc.capDst)

			cnt, err := bf.Write(bufSrc)

			if err != tc.errExp {
				t.Errorf("expected error to be: %v but got: %v", tc.errExp, err)
			}

			if cnt != tc.capSrc {
				t.Errorf("expected copied bytes to be: %v but got: %v", tc.capSrc, cnt)
			}
		})
	}
}

func TestSeek(t *testing.T) {
	testCases := []struct {
		name      string // name of test case
		offset    int64  // seek to this position
		whence    int    // starting position
		offsetExp int64  // expected location after seeking
		errExp    error  // expected error
	}{
		{
			name:      "case1",
			offset:    1,
			whence:    io.SeekStart,
			offsetExp: 1,
			errExp:    nil,
		},
		{
			name:      "case2",
			offset:    1,
			whence:    io.SeekCurrent,
			offsetExp: 2,
			errExp:    nil,
		},
		{
			name:      "case3",
			offset:    1,
			whence:    io.SeekEnd,
			offsetExp: 42,
			errExp:    nil,
		},
		{
			name:      "case4",
			offset:    -44,
			whence:    io.SeekEnd,
			offsetExp: 42,
			errExp:    errSeekToNegativeLocation,
		},
	}

	// prepare buffer file with random bytes
	buf := make([]byte, 42)
	rand.Read(buf)
	bf := NewBufferFileFromBytes(buf)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			offset, err := bf.Seek(tc.offset, tc.whence)

			if offset != tc.offsetExp {
				t.Errorf("expected offset to be %d but got %d", tc.offsetExp, offset)
			}
			if err != tc.errExp {
				t.Errorf("expected error to be %v but got %v", tc.offsetExp, err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	bf := NewBufferFile()

	if err := bf.Close(); err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
}
