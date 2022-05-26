package errors

import (
	"fmt"
)

type wrapper struct {
	trace string
	err   error
	frame
}

func Wrap(err error, trace string) error {
	if err == nil {
		return nil
	}
	e := &wrapper{
		trace: trace,
		err:   err,
		frame: buildFrame(1),
	}
	return e
}

func Wrapf(err error, format string, a ...interface{}) error {
	e := &wrapper{
		trace: fmt.Sprintf(format, a...),
		err:   err,
		frame: buildFrame(1),
	}
	return e
}

func (e *wrapper) Unwrap() error {
	return e.err
}

func (e *wrapper) Error() string {
	return e.trace
}

func (e *wrapper) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e), nil
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

func (e *wrapper) fmt() fmtWrapper {
	return fmtWrapper{e.trace, e.frame.String()}
}

type fmtWrapper struct {
	trace      string
	frameCache string
}

func (f *fmtWrapper) jsonSize() int {
	return 10 + len(f.trace) + 12 + len(f.frameCache) + 2
}

func (f *fmtWrapper) textSize() int {
	return 6 + len(f.trace) + len(f.frameCache)
}

func (f *fmtWrapper) json(bs []byte) []byte {
	bs = append(bs, `{"trace":"`...)
	bs = append(bs, f.trace...)
	bs = append(bs, `","caller":"`...)
	bs = append(bs, f.frameCache...)
	bs = append(bs, `"}`...)
	return bs
}

func (f *fmtWrapper) text(bs []byte) []byte {
	bs = append(bs, f.trace...)
	bs = append(bs, ",\n    "...)
	bs = append(bs, f.frameCache...)
	return bs
}
