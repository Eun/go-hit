package doppelgangerreader

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

// DoppelgangerFactory is a reader that mimics the behaviour of an other reader
// it can be used to read readers multiple times
type DoppelgangerFactory struct {
	source   io.Reader
	readers  []*readerInstance
	buffer   bytes.Buffer
	mu       sync.Mutex
	closedOn *int
}

// NewFactory creates a new DoppelgangerFactory with the original reader specified
func NewFactory(readerToMimic io.Reader) *DoppelgangerFactory {
	reader := &DoppelgangerFactory{
		source: readerToMimic,
	}
	return reader
}

// NewDoppelganger creates a new reader that acts like the original reader
func (factory *DoppelgangerFactory) NewDoppelganger() io.ReadCloser {
	factory.mu.Lock()
	reader := &readerInstance{
		DoppelBase: factory,
		// prefill Buffer with already collected data
		Buffer: bytes.NewBuffer(factory.buffer.Bytes()),
	}
	if factory.closedOn == nil {
		// only add to readers if there is still data to consume
		factory.readers = append(factory.readers, reader)
	}
	factory.mu.Unlock()
	return reader
}

// RemoveDoppelganger a created reader from receiving new data
func (factory *DoppelgangerFactory) RemoveDoppelganger(r io.ReadCloser) error {
	instance, ok := r.(*readerInstance)
	if !ok {
		return errors.New("not a reader instance")
	}
	factory.mu.Lock()
	defer factory.mu.Unlock()
	for i := len(factory.readers) - 1; i >= 0; i-- {
		if factory.readers[i] == instance {
			factory.readers[i].DoppelBase = nil
			factory.readers = append(factory.readers[:i], factory.readers[i+1:]...)
			return nil
		}
	}

	return errors.New("reader not found")
}

// Close the DoppelgangerFactory and stops all created Doppelgangers from receiving data
// (does not close the underlying reader)
func (factory *DoppelgangerFactory) Close() error {
	// we already closed
	if factory.closedOn != nil {
		return nil
	}
	factory.closedOn = new(int)
	*factory.closedOn = factory.buffer.Len()

	// remove all readers because everything has been consumed
	factory.readers = nil

	return nil
}

func (factory *DoppelgangerFactory) read(caller *readerInstance, p []byte) (int, error) {
	if factory.closedOn != nil {
		return 0, io.EOF
	}
	n, err := factory.source.Read(p)

	if n > 0 {
		// fill my own Buffer if we have data
		factory.buffer.Write(p[:n])
	}

	if err != nil {
		return n, err
	}

	for i := len(factory.readers) - 1; i >= 0; i-- {
		if factory.readers[i] != caller {
			factory.readers[i].Buffer.Write(p[:n])
		}
	}
	return n, nil
}

type readerInstance struct {
	DoppelBase *DoppelgangerFactory
	Buffer     *bytes.Buffer
}

func (r *readerInstance) Read(p []byte) (n int, err error) {
	if r.DoppelBase == nil {
		return 0, io.EOF
	}
	r.DoppelBase.mu.Lock()
	if r.Buffer.Len() > 0 {
		n, err = r.Buffer.Read(p)
	} else {
		n, err = r.DoppelBase.read(r, p)
		if err != nil {
			r.DoppelBase.Close()
		}
	}
	r.DoppelBase.mu.Unlock()
	return n, err
}

func (r *readerInstance) Close() error {
	// if the factory is already closed
	// we dont need to remove
	if r.DoppelBase.closedOn == nil {
		return r.DoppelBase.RemoveDoppelganger(r)
	}
	return nil
}
