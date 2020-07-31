package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var local, _ = time.LoadLocation("Asia/Shanghai")

type JSONTime time.Time

func (t JSONTime) String() string {
	return time.Time(t).In(local).Format("2006-01-02 15:04:05")
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, t.String())
	return []byte(stamp), nil
}

type JSONDateTime time.Time

func (dt JSONDateTime) String() string {
	return time.Time(dt).In(local).Format("20060102")
}

func (dt JSONDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.String())
}

type TimestampTime time.Time

func (t TimestampTime) MarshalJSON() ([]byte, error) {
	val := strconv.Itoa(int(time.Time(t).Unix()))
	return []byte(val), nil
}
