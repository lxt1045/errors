package errors

import (
	"bytes"
	"fmt"
)

type wrapper struct {
	trace string
	err   error
	frame
}

func (e *wrapper) Error() string {
	return e.trace
}

func (e *wrapper) Unwrap() error {
	return e.err
}

func (e *wrapper) json(buf *bytes.Buffer) {
	buf.WriteString(`{"trace":"`)
	buf.WriteString(e.trace)
	buf.WriteString(`","caller":`)
	buf.WriteByte('"')
	buf.WriteString(e.String())
	buf.WriteString(`"}`)
}

func (e *wrapper) text(buf *bytes.Buffer) {
	if e.trace == "" {
		e.trace = "-"
	}
	buf.WriteString(e.trace)
	buf.WriteString(",\n    ")
	buf.WriteString(e.String())
}

func (e *wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			s.Write(MarshalText(e))
			return
		}
		fallthrough
	case 's':
		s.Write([]byte(e.Error()))
		return
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// MarshalJSON json.Marshaler的方法, json.Marshal 里调用
func (e *wrapper) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e), nil
}

func Wrap(err error, trace string) error {
	if err == nil {
		return nil
	}
	e := &wrapper{
		trace: trace,
		err:   err,
		frame: NewFrame(1),
	}
	return e
}

func Wrapf(err error, format string, a ...interface{}) error {
	e := &wrapper{
		trace: fmt.Sprintf(format, a...),
		err:   err,
		frame: NewFrame(1),
	}
	return e
}
