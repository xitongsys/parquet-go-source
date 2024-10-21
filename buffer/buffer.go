package buffer

import (
	"errors"
	"io"

	"github.com/xitongsys/parquet-go/source"
)

// BufferFile allows reading parquet messages from a memory buffer.
type BufferFile struct {
	buff []byte
	loc  int
}

// DefaultCapacity is the size in bytes of a new BufferFile's backing buffer
const DefaultCapacity = 512

var errSeekToNegativeLocation = errors.New("unable to seek to a location <0")

// NewBufferFile creates new in memory parquet buffer.
func NewBufferFile() *BufferFile {
	return NewBufferFileCapacity(DefaultCapacity)
}

// NewBufferFileFromBytes creates new in memory parquet buffer from the given bytes.
// It allocates a new slice and copy the contents of s.
func NewBufferFileFromBytes(s []byte) *BufferFile {
	b := make([]byte, len(s))
	copy(b, s)
	return &BufferFile{buff: b}
}

// NewBufferFileFromBytes creates new in memory parquet buffer from the given bytes.
// It uses the provided slice as its buffer.
func NewBufferFileFromBytesNoAlloc(s []byte) *BufferFile {
	return &BufferFile{buff: s}
}

// NewBufferFileCapacity starts the returned BufferFile with the given capacity
func NewBufferFileCapacity(cap int) *BufferFile {
	return &BufferFile{buff: make([]byte, 0, cap)}
}

// NewBufferFileFromBytesZeroAlloc creates new in memory parquet buffer without memory allocation.
func NewBufferFileFromBytesZeroAlloc(s []byte) *BufferFile {
	return &BufferFile{buff: s}
}

func (bf BufferFile) Create(string) (source.ParquetFile, error) {
	return NewBufferFile(), nil
}

func (bf BufferFile) Open(string) (source.ParquetFile, error) {
	return NewBufferFileFromBytes(bf.buff), nil
}

// Seek seeks in the underlying memory buffer.
func (bf *BufferFile) Seek(offset int64, whence int) (int64, error) {
	newLoc := bf.loc
	switch whence {
	case io.SeekStart:
		newLoc = int(offset)
	case io.SeekCurrent:
		newLoc += int(offset)
	case io.SeekEnd:
		newLoc = len(bf.buff) + int(offset)
	}

	if newLoc < 0 {
		return int64(bf.loc), errSeekToNegativeLocation
	}

	if newLoc > len(bf.buff) {
		newLoc = len(bf.buff)
	}

	bf.loc = newLoc

	return int64(bf.loc), nil
}

// Read reads data form BufferFile into p.
func (bf *BufferFile) Read(p []byte) (n int, err error) {
	n = copy(p, bf.buff[bf.loc:len(bf.buff)])
	bf.loc += n

	// if copied files into buffer with not enaugh capacity
	if bf.loc < len(bf.buff) {
		return n, io.EOF
	}

	return n, nil
}

// Write writes data from p into BufferFile.
func (bf *BufferFile) Write(p []byte) (n int, err error) {
	// Do we have space?
	if available := cap(bf.buff) - bf.loc; available < len(p) {
		// How much should we expand by?
		addCap := cap(bf.buff)
		if addCap < len(p) {
			addCap = len(p)
		}

		newBuff := make([]byte, len(bf.buff), cap(bf.buff)+addCap)

		copy(newBuff, bf.buff)

		bf.buff = newBuff
	}

	// Write
	n = copy(bf.buff[bf.loc:cap(bf.buff)], p)
	bf.loc += n
	if len(bf.buff) < bf.loc {
		bf.buff = bf.buff[:bf.loc]
	}

	return n, nil
}

// Close is a no-op for a memory buffer.
func (bf BufferFile) Close() error {
	return nil
}

// Bytes returns the byte slice representing the buffer file.
func (bf BufferFile) Bytes() []byte {
	return bf.buff
}
