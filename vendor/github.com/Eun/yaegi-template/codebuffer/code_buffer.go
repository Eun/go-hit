// Package codebuffer is a package that provides a method to read trough a reader and separate the contents into
// code and text parts.
package codebuffer

import (
	"io"

	"github.com/pkg/errors"

	"go.uber.org/atomic"
)

// PartType represents a part type.
type PartType uint8

const (
	// TextPartType represents a text part.
	TextPartType PartType = iota
	// CodePartType represents a code part.
	CodePartType
)

// Part represents one part of the text.
type Part struct {
	Type    PartType
	Content []byte
}

const (
	notReadState int32 = iota
	readingState
	readState
)

// Iterator represents an iterator that can be used to walk trough the CodeBuffer.
type Iterator interface {
	Next() bool
	Value() *Part
	Error() error
}

// CodeBuffer is a construct to walk over a text separating code from text blocks.
type CodeBuffer struct {
	r           io.Reader
	startTokens []rune
	endTokens   []rune
	parts       []*Part
	state       *atomic.Int32
}

// InReadingState is a error that will be returned when the CodeBuffer is busy on another thread.
type InReadingState struct{}

// Error returns the error text for InReadingState.
func (InReadingState) Error() string {
	return "the reader was not fully consumed yet"
}

// New creates a new CodeBuffer with the specified reader.
func New(r io.Reader, startTokens, endTokens []rune) *CodeBuffer {
	return &CodeBuffer{
		r:           r,
		startTokens: startTokens,
		endTokens:   endTokens,
		state:       atomic.NewInt32(notReadState),
	}
}

// Iterator returns an iterator that can be used to walk trough the CodeBuffer.
func (c *CodeBuffer) Iterator() (Iterator, error) {
	switch c.state.Load() {
	case notReadState:
		c.state.Store(readingState)
		return newLiveIterator(c.state, &c.parts, c.r, c.startTokens, c.endTokens)
	case readingState:
		return nil, InReadingState{}
	case readState:
		return newCacheIterator(c.parts)
	}
	return nil, errors.New("unknown error")
}
