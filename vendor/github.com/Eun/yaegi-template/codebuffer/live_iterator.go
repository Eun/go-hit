package codebuffer

import (
	"io"

	"bufio"
	"unicode/utf8"

	"bytes"

	"fmt"

	"unicode"

	"io/ioutil"

	"go.uber.org/atomic"
)

type liveIterator struct {
	state                   *atomic.Int32
	parts                   *[]*Part
	reader                  *bufio.Reader
	startSequence           []rune
	endSequence             []rune
	err                     error
	inCodeBlock             bool
	currentPart             *Part
	hasNext                 bool
	stripLeadingWhiteSpaces bool
}

func newLiveIterator(state *atomic.Int32, parts *[]*Part, reader io.Reader, startSequence, endSequence []rune) (Iterator, error) {
	return &liveIterator{
		state:         state,
		parts:         parts,
		reader:        bufio.NewReaderSize(reader, utf8.UTFMax),
		startSequence: startSequence,
		endSequence:   endSequence,
		hasNext:       true,
	}, nil
}

func (i *liveIterator) Next() bool {
	// early exit
	if !i.hasNext {
		i.currentPart = nil
		return false
	}

	var read func() (*Part, bool, error)
	if i.inCodeBlock {
		read = i.readCodeBlock
	} else {
		read = i.readTextBlock
	}

	part, stopAfterThisPart, err := read()
	added := i.addPart(part)
	if err != nil {
		return i.stopProcessing(err)
	}

	if stopAfterThisPart {
		// should we stop after this part and we did not add anything? => stop processing
		if !added {
			return i.stopProcessing(err)
		}
		return true
	}

	// we should continue, but we added no part => go to the next item
	if !added {
		return i.Next()
	}
	// we have this part
	return true
}

func (i *liveIterator) stopProcessing(err error) bool {
	i.state.Store(readState)
	i.hasNext = false
	i.currentPart = nil
	i.err = err
	return false
}

//nolint:gocognit // allow more complex code here
func (i *liveIterator) readTextBlock() (*Part, bool, error) {
	sequenceSize := len(i.startSequence)
	if sequenceSize == 0 {
		// shortcut, also a special case, if there is no sequence present treat everything as code
		p, err := ioutil.ReadAll(i.reader)
		if err != nil {
			if err != io.EOF {
				return nil, true, err
			}
			err = nil
		}
		return constructCodePath(p), true, err
	}
	pos := 0
	seqBuffer := make([]rune, sequenceSize)

	var contentBuffer bytes.Buffer

	writeSeqBuffer := func() error {
		if pos <= 0 {
			return nil
		}
		if err := writeRunes(&contentBuffer, seqBuffer[:pos]); err != nil {
			return err
		}
		pos = 0
		return nil
	}

	for {
		r, rsize, err := readRune(i.reader)
		if err != nil {
			return nil, true, err
		}

		if r == unicode.ReplacementChar {
			if err := writeSeqBuffer(); err != nil {
				return nil, true, err
			}
			stripLeadingWhiteSpaces := i.stripLeadingWhiteSpaces
			i.stripLeadingWhiteSpaces = false // reset strip leading whitespaces
			return constructTextPart(contentBuffer.Bytes(), stripLeadingWhiteSpaces), true, nil
		}

		if r == i.startSequence[pos] { //nolint:nestif // moving this block into a function would make this more complex
			// store the rune if this is not the sequence
			seqBuffer[pos] = r
			pos++
			if pos != sequenceSize {
				continue
			}
			// we found all runes of the sequence
			content := contentBuffer.Bytes()

			// test if the next rune is a "-" indicating we should strip previous white spaces
			r, _, err = readRune(i.reader)
			if err != nil {
				return nil, true, err
			}

			if r == '-' {
				content = bytes.TrimRightFunc(content, unicode.IsSpace)
			} else { //nolint:elseif,gocritic // for better readability keep a nested if
				// its not an "-"
				if err := i.reader.UnreadRune(); err != nil {
					return nil, true, err
				}
			}

			i.inCodeBlock = true
			stripLeadingWhiteSpaces := i.stripLeadingWhiteSpaces
			i.stripLeadingWhiteSpaces = false // reset strip leading whitespaces
			return constructTextPart(content, stripLeadingWhiteSpaces), false, nil
		}
		if err := writeSeqBuffer(); err != nil {
			return nil, true, err
		}

		wsize, err := contentBuffer.WriteRune(r)
		if err != nil {
			return nil, true, err
		}
		if wsize != rsize {
			return nil, true, fmt.Errorf("expected to write %d bytes, written %d", rsize, wsize)
		}
	}
}

//nolint:gocognit  // allow more complex code here
func (i *liveIterator) readCodeBlock() (*Part, bool, error) {
	sequenceSize := len(i.endSequence)
	if sequenceSize == 0 {
		// shortcut
		p, err := ioutil.ReadAll(i.reader)
		if err != nil {
			if err != io.EOF {
				return nil, true, err
			}
			err = nil
		}
		return constructCodePath(p), true, err
	}
	pos := 0
	seqBuffer := make([]rune, sequenceSize)
	var contentBuffer bytes.Buffer
	previousRune := unicode.ReplacementChar

	writeSeqBuffer := func() error {
		if pos < 0 {
			return nil
		}
		if err := writeRunes(&contentBuffer, seqBuffer[:pos]); err != nil {
			return err
		}
		pos = 0
		return nil
	}

	for {
		r, rsize, err := readRune(i.reader)
		if err != nil {
			return nil, true, err
		}

		if r == unicode.ReplacementChar {
			if err := writeSeqBuffer(); err != nil {
				return nil, true, err
			}
			return constructCodePath(contentBuffer.Bytes()), true, nil
		}

		if r == i.endSequence[pos] {
			// we found the current needed rune
			seqBuffer[pos] = r
			pos++
			if pos != sequenceSize {
				continue
			}
			if previousRune == '-' {
				if n := contentBuffer.Len(); n > 0 {
					// remove the -
					contentBuffer.Truncate(contentBuffer.Len() - 1)
				}
				i.stripLeadingWhiteSpaces = true
			}

			i.inCodeBlock = false
			return constructCodePath(contentBuffer.Bytes()), false, nil
		}

		if err := writeSeqBuffer(); err != nil {
			return nil, true, err
		}

		wsize, err := contentBuffer.WriteRune(r)
		if err != nil {
			return nil, true, err
		}
		if wsize != rsize {
			return nil, true, fmt.Errorf("expected to write %d bytes, written %d", rsize, wsize)
		}

		previousRune = r
	}
}

func (i *liveIterator) addPart(p *Part) bool {
	if p == nil || len(p.Content) == 0 {
		i.currentPart = nil
		return false
	}
	*i.parts = append(*i.parts, p)
	i.currentPart = p
	return true
}

func constructTextPart(content []byte, trimLeadingSpaces bool) *Part {
	if trimLeadingSpaces {
		content = bytes.TrimLeftFunc(content, unicode.IsSpace)
	}

	if len(content) == 0 {
		return nil
	}
	return &Part{
		Type:    TextPartType,
		Content: content,
	}
}

func constructCodePath(content []byte) *Part {
	if len(content) == 0 {
		return nil
	}
	return &Part{
		Type:    CodePartType,
		Content: content,
	}
}

func (i *liveIterator) Value() *Part {
	return i.currentPart
}

func (i *liveIterator) Error() error {
	return nil
}

type runeWriter interface {
	WriteRune(rune) (int, error)
}

func writeRunes(w runeWriter, runes []rune) error {
	for _, r := range runes {
		if _, err := w.WriteRune(r); err != nil {
			return err
		}
	}
	return nil
}

type runeReader interface {
	ReadRune() (rune, int, error)
}

func readRune(rd runeReader) (r rune, size int, err error) {
	r, size, err = rd.ReadRune()
	if err == io.EOF {
		r = unicode.ReplacementChar
		size = 0
		err = nil
	}
	return r, size, err
}
