package errors

import (
	"strings"
	"time"
)

func newBreadcrumb(skip, level int, trace string) breadcrumb {
	return breadcrumb{
		Trace: trace,
		Time:  now(),
		stack: newStack(1+skip, level),
	}
}

//breadcrumb 用于追踪调用栈的"面包屑",来自童话故事“汉赛尔和格莱特
type breadcrumb struct {
	Trace string //足迹
	Time  Time   //时间

	stack stack // 调用栈
}

func (b breadcrumb) json(buf *strings.Builder) *strings.Builder {
	buf.Write([]byte(`{"trace":`))
	buf.WriteByte('"')
	buf.WriteString(b.Trace)
	buf.WriteByte('"')
	buf.Write([]byte(`,"time":`))
	buf.WriteByte('"')
	buf.Write(b.Time.Bytes())
	buf.WriteByte('"')
	buf.WriteByte(',')
	b.stack.json(buf)
	buf.WriteByte('}')

	return buf
}
func (b breadcrumb) text(buf *strings.Builder) *strings.Builder {
	buf.Write(b.Time.Bytes())
	buf.Write([]byte(`, `))
	if b.Trace == "" {
		b.Trace = "-"
	}
	buf.WriteString(b.Trace)

	buf.Write([]byte(",\n"))
	b.stack.text(buf)

	return buf
}

func now() Time {
	return (Time)(time.Now())
}

//Time 提高 time.Time 序列化性能
type Time time.Time

func (tt Time) Bytes() []byte {
	t := time.Time(tt)
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 32)

	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = ' ' //'T'
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	buf[23] = byte((micros/100)%10) + '0'
	buf[24] = byte((micros/10)%10) + '0'
	buf[25] = byte((micros)%10) + '0'
	//buf[26] = ' ' //'Z'

	_, offset := t.Zone()
	zone := offset / 60 // convert to minutes
	absoffset := offset
	if zone < 0 {
		buf[26] = byte((micros)%10) + '0'
		zone = -zone
		absoffset = -absoffset
	} else {
		buf[26] = '+'
	}
	hour, minute = zone/60, zone%60
	buf[27] = byte((hour)/10) + '0'
	buf[28] = byte((hour)%10) + '0'
	buf[29] = ':'
	buf[30] = byte((minute)/10) + '0'
	buf[31] = byte((minute)%10) + '0'

	return buf
}
