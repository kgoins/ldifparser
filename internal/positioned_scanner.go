package internal

import (
	"bufio"
	"io"
)

type PositionedScanner struct {
	pos      int64
	buffSize int
	scanner  *bufio.Scanner
}

func NewPositionedScanner(inputStream io.Reader, buffSize ...int) *PositionedScanner {
	positionedScanner := PositionedScanner{
		pos:     0,
		scanner: bufio.NewScanner(inputStream),
	}

	lineSize := 256
	if len(buffSize) > 0 {
		lineSize = buffSize[0]
	}
	positionedScanner.buffSize = lineSize

	lineBuf := make([]byte, lineSize)
	positionedScanner.Buffer(lineBuf, lineSize)

	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		positionedScanner.pos += int64(advance)
		return
	}
	positionedScanner.scanner.Split(scanLines)

	return &positionedScanner
}

func (ps PositionedScanner) Fork(input io.ReadSeeker) (*PositionedScanner, error) {
	if _, err := input.Seek(ps.Position(), 0); err != nil {
		return nil, err
	}

	newScanner := NewPositionedScanner(input)
	newScanner.pos = ps.Position()
	return newScanner, nil
}

func (ps PositionedScanner) Buffer(buffer []byte, max int) {
	ps.scanner.Buffer(buffer, max)
}

func (ps *PositionedScanner) Scan() bool {
	return ps.scanner.Scan()
}

func (ps PositionedScanner) Position() int64 {
	return ps.pos
}

func (ps PositionedScanner) Text() string {
	return ps.scanner.Text()
}

func (ps PositionedScanner) Err() error {
	return ps.scanner.Err()
}
