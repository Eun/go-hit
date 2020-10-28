package codebuffer

import (
	"io"

	"errors"

	"go.uber.org/atomic"
)

type PartType uint8

const (
	TextPartType PartType = iota
	CodePartType
)

type Part struct {
	Type    PartType
	Content []byte
}

const (
	notReadState int32 = iota
	readingState
	readState
)

type Iterator interface {
	Next() bool
	Value() *Part
	Error() error
}

type CodeBuffer struct {
	r           io.Reader
	startTokens []rune
	endTokens   []rune
	parts       []*Part
	state       *atomic.Int32
}

type InReadingState struct{}

func (InReadingState) Error() string {
	return "the reader was not fully consumed yet"
}

func New(r io.Reader, startTokens, endTokens []rune) *CodeBuffer {
	return &CodeBuffer{
		r:           r,
		startTokens: startTokens,
		endTokens:   endTokens,
		state:       atomic.NewInt32(notReadState),
	}
}

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
