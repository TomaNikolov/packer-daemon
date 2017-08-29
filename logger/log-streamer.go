package logger

import (
	"bytes"
	"io"
	"strings"

	"github.com/tomanikolov/packer-daemon/constants"
)

// Logstreamer ...
type Logstreamer struct {
	Logger    *Logger
	buf       *bytes.Buffer
	readLines string
	prefix    string
	// if true, saves output in memory
	record  bool
	persist string
}

// NewLogstreamer ...
func NewLogstreamer(logger *Logger, prefix string, record bool) *Logstreamer {
	return &Logstreamer{
		Logger:  logger,
		buf:     bytes.NewBuffer([]byte("")),
		prefix:  prefix,
		record:  record,
		persist: "",
	}
}

// Write ...
func (l *Logstreamer) Write(p []byte) (n int, err error) {
	if n, err = l.buf.Write(p); err != nil {
		return
	}

	err = l.OutputLines()
	return
}

// Close ...
func (l *Logstreamer) Close() error {
	l.Flush()
	l.buf = bytes.NewBuffer([]byte(""))
	return nil
}

// Flush ...
func (l *Logstreamer) Flush() error {
	var p []byte
	if _, err := l.buf.Read(p); err != nil {
		return err
	}

	l.out(string(p))
	return nil
}

// OutputLines ...
func (l *Logstreamer) OutputLines() (err error) {
	for {
		line, err := l.buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if line == "\n" {
			continue
		}

		l.readLines += line
		l.out(line)
	}

	return nil
}

// ResetReadLines ...
func (l *Logstreamer) ResetReadLines() {
	l.readLines = ""
}

// ReadLines ...
func (l *Logstreamer) ReadLines() string {
	return l.readLines
}

// FlushRecord ...
func (l *Logstreamer) FlushRecord() string {
	buffer := l.persist
	l.persist = ""
	return buffer
}

func (l *Logstreamer) out(str string) (err error) {
	if l.record == true {
		l.persist = l.persist + str
	}

	if l.prefix == constants.Stdout {
		l.Logger.Log(strings.Trim(str, "\n"))
	} else if l.prefix == constants.Stderr {
		l.Logger.LogError(strings.Trim(str, "\n"))
	} else {
		str = l.prefix + str
	}

	return nil
}
