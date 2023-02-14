package util

import (
	"time"
)

type JsonTime struct {
	time.Time
}

func (j *JsonTime) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	(*j).Time = t
	return err
}

func (j *JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(j.Local().Format(`"2006-01-02 15:04:05"`)), nil
}
