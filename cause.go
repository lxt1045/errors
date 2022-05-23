package errors

import (
	"bytes"
	"strconv"
)

//Cause 是一个带错误栈信息的error
type Cause struct {
	Code int    //业务错误码
	Msg  string //业务错误信息

	stack stack //错误的现场
}

func buildCause(code int, msg string, stack stack) (e *Cause) {
	return &Cause{
		Code:  code,
		Msg:   msg,
		stack: stack,
	}
}

func (e *Cause) GetCode() int {
	return e.Code
}

func (e *Cause) GetMsg() string {
	return e.Msg
}

func (e *Cause) Is(err error) bool {
	to, ok := err.(*Cause)
	return ok && e.Code != -1 && e.Code == to.Code
}

//Error error interface, 序列化为string, 包含调用栈
func (e *Cause) Error() string {
	buf := bytes.NewBuffer(nil)
	e.text(buf)
	return bufToString(buf)
}

// MarshalJSON json.Marshaler的方法, json.Marshal 里调用
func (e *Cause) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	e.json(buf)
	return buf.Bytes(), nil
}

func (e *Cause) json(buf *bytes.Buffer) {
	buf.WriteString(`{"code":`)
	buf.WriteString(strconv.Itoa(e.Code))
	buf.WriteString(`,"msg":"`)
	buf.WriteString(e.Msg)
	buf.WriteByte('"')
	if e.stack.npc > 0 {
		buf.WriteString(`,"stack":`)
		e.stack.json(buf)
	}
	buf.WriteByte('}')
}

func (e *Cause) text(buf *bytes.Buffer) {
	buf.WriteString(strconv.Itoa(e.Code))
	buf.WriteString(`, `)
	if e.Msg == "" {
		e.Msg = "-"
	}
	buf.WriteString(e.Msg)
	if e.stack.npc > 0 {
		buf.WriteString(";\n")
		e.stack.text(buf)
		buf.WriteString(";")
	}
}
