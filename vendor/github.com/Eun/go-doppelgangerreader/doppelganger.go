package doppelgangerreader

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

// DoppelgangerFactory is a reader that mimics the behaviour of an other reader
// it can be used to read readers multiple times
type DoppelgangerFactory interface {
	NewDoppelganger() io.ReadCloser
	RemoveDoppelganger(r io.ReadCloser) error
	Close() error
}

// NewFactory creates a new DoppelgangerFactory with the original reader specified
// if the reader is already a Doppelganger it will return the original factory
func NewFactory(readerToMimic io.Reader) DoppelgangerFactory {
	factory := GetFactory(readerToMimic)
	if factory != nil {
		return &nestedDoppelgangerFactory{
			parent: factory,
		}
	}
	return &doppelgangerFactory{
		source: readerToMimic,
	}
}

type doppelgangerFactory struct {
	source   io.Reader
	readers  []*readerInstance
	buffer   bytes.Buffer
	mu       sync.Mutex
	closedOn *int
}

// NewDoppelganger creates a new reader that acts like the original reader
func (factory *doppelgangerFactory) NewDoppelganger() io.ReadCloser {
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
func (factory *doppelgangerFactory) RemoveDoppelganger(r io.ReadCloser) error {
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
func (factory *doppelgangerFactory) Close() error {
	// this is a public function so make sure we lock
	factory.mu.Lock()
	err := factory.close()
	factory.mu.Unlock()
	return err
}

// close the DoppelgangerFactory and stops all created Doppelgangers from receiving data
// (does not close the underlying reader)
func (factory *doppelgangerFactory) close() error {
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

func (factory *doppelgangerFactory) read(caller *readerInstance, p []byte) (int, error) {
	if factory.closedOn != nil {
		return 0, io.EOF
	}
	if factory.source == nil {
		return 0, NilReaderError{}
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
	DoppelBase *doppelgangerFactory
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
			r.DoppelBase.close()
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

// NilReaderError will be reported if the provided reader is nil
type NilReaderError struct{}

// Error returns the error message
func (NilReaderError) Error() string {
	return "Reader to mimic is nil"
}

// IsNilReaderError returns true if the specified error is a NilReaderError
func IsNilReaderError(e error) bool {
	_, ok := e.(NilReaderError)
	return ok
}

// GetFactory returns the DoppelgangerFactory if the reader is a Doppelganger
func GetFactory(reader io.Reader) DoppelgangerFactory {
	if v, ok := reader.(*readerInstance); ok {
		return v.DoppelBase
	}
	return nil
}

type nestedDoppelgangerFactory struct {
	parent  DoppelgangerFactory
	readers []io.ReadCloser
}

func (factory *nestedDoppelgangerFactory) NewDoppelganger() io.ReadCloser {
	r := factory.parent.NewDoppelganger()
	factory.readers = append(factory.readers, r)
	return r
}

func (factory *nestedDoppelgangerFactory) RemoveDoppelganger(r io.ReadCloser) error {
	return factory.parent.RemoveDoppelganger(r)
}

func (factory *nestedDoppelgangerFactory) Close() error {
	for i := len(factory.readers) - 1; i >= 0; i-- {
		if err := factory.RemoveDoppelganger(factory.readers[i]); err != nil {
			return err
		}
	}
	factory.readers = nil
	return nil
}
