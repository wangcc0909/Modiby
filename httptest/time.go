package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TimeDuration time.Duration

func (d TimeDuration) String() string {
	return time.Duration(d).String()
}

func (d TimeDuration) MarshalJson() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *TimeDuration) UnmarshalJson(b []byte) error {
	str := string(b)
	if str != "" && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		*d = TimeDuration(val)
		return nil
	}
	t, err := time.ParseDuration(str)
	if err == nil {
		*d = TimeDuration(t)
		return nil
	}
	return fmt.Errorf("invalid duration type %T, value: '%s'", b, b)
}
