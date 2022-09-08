package common

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Time time.Time

var ZERO_TIME Time

func (t *Time) UnmarshalJSON(dateTimeByte []byte) (err error) {
	dateTimeStr := strings.Trim(string(dateTimeByte), `"`)
	if "0000-00-00 00:00:00" == dateTimeStr {
		*t = ZERO_TIME
		return nil
	}
	now, err := time.ParseInLocation(F_DATETIME, dateTimeStr, time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(F_DATETIME)+2)
	b = append(b, '"')
	if time.Time(t).IsZero() {
		b = append(b, []byte("0000-00-00 00:00:00")...)
	} else {
		b = time.Time(t).AppendFormat(b, F_DATETIME)
	}
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	if time.Time(t).IsZero() {
		return "0000-00-00 00:00:00"
	}
	return time.Time(t).Format(F_DATETIME)
}

func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

func (ts Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(ts)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return "0000-00-00 00:00:00", nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (ts *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*ts = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

//判断是否0000-00-00 00:00:00
func (t Time) IsZero() bool {
	return time.Time(t).IsZero()
}
