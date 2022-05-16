package errors

import (
	"time"
)

func now() Time {
	return Time(time.Now())
}

// Time 提高 time.Time 序列化性能
type Time time.Time

//nolint
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
	buf[10] = ' ' // or 'T'
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

	// buf[26] = ' ' // or 'Z'
	_, offset := t.Zone()
	zone := offset / 60 // convert to minutes
	if zone < 0 {
		buf[26] = byte((micros)%10) + '0'
		zone = -zone
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
