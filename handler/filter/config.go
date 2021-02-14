package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Time time.Time

var _ envconfig.Decoder = &Time{}

func (t *Time) Decode(value string) error {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}
	*t = Time(parsed)
	return nil
}

func (t *Time) ToTime() time.Time {
	return time.Time(*t)
}

func (t *Time) IsZero() bool {
	return t.ToTime().IsZero()
}