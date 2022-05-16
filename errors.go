package errors

import (
	"bytes"
	"fmt"
	"strconv"
)

//Err 是一个带错误栈信息的error
type Err struct {
	Code int    //业务错误码
	Msg  string //业务错误信息

	cause  stack   //错误的现场
	traces []frame //Warp()组成的路径
}

func (e Err) GetCode() int {
	return e.Code
}
func (e Err) GetMessage() string {
	return e.Msg
}

func (e *Err) Wrap(trace string) Error {
	e.traces = append(e.traces, buildFrame(1, trace))
	return e
}
func (e *Err) Wrapf(format string, a ...interface{}) Error {
	e.traces = append(e.traces, buildFrame(1, fmt.Sprintf(format, a...)))
	return e
}
func (e *Err) wrap(trace string) Error {
	e.traces = append(e.traces, buildFrame(2, trace))
	return e
}
func (e *Err) Unwrap() error {
	if len(e.traces) <= 0 {
		return nil
	}
	e.traces = e.traces[:len(e.traces)-1]
	return e
}
func (e *Err) Is(err error) bool {
	to, ok := err.(Error)
	return ok && e.Code != -1 && e.Code == to.GetCode()
}

//Error error interface, 序列化为string, 包含调用栈
func (e *Err) Error() string {
	return e.String()
}

func (e *Err) String() string {
	buf := bytes.NewBuffer(nil)
	e.text(buf)
	return bufToString(buf)
}

// MarshalJSON json.Marshaler的方法, json.Marshal 里调用
func (e *Err) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	e.json(buf)
	return buf.Bytes(), nil
}

func (e *Err) json(buf *bytes.Buffer) {
	buf.Write([]byte(`{"code":`))
	buf.WriteString(strconv.Itoa(e.Code))
	buf.Write([]byte(`,"message":`))
	buf.WriteByte('"')
	buf.WriteString(e.Msg)
	buf.WriteByte('"')
	if e.cause.nrpc > 0 {
		buf.Write([]byte(`,"stack":`))
		e.cause.json(buf)
		if len(e.traces) > 0 {
			buf.Write([]byte(`,"traces":[`))
			for i, b := range e.traces {
				if i != 0 {
					buf.WriteByte(',')
				}
				b.json(buf)
			}
			buf.WriteByte(']')
		}
	}
	buf.WriteByte('}')
}

func (e *Err) text(buf *bytes.Buffer) {
	buf.WriteString(strconv.Itoa(e.Code))
	buf.Write([]byte(`, `))
	if e.Msg == "" {
		e.Msg = "-"
	}
	buf.WriteString(e.Msg)
	buf.Write([]byte(";"))

	if e.cause.nrpc > 0 {
		buf.Write([]byte("\n"))
		e.cause.text(buf)
		buf.Write([]byte(";\n"))

		for _, b := range e.traces {
			b.text(buf)
			buf.Write([]byte(";\n"))
		}
	}
}
