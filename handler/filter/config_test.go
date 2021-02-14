package main

import (
	"os"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
)

func TestTime_Decode(t *testing.T) {
	type config struct {
		Time Time `envconfig:"TIME"`
	}
	tokyo, _ := time.LoadLocation("asia/Tokyo")

	cases := []struct {
		Env map[string]string
		Time time.Time
	}{
		{
			Env: map[string]string{
				"TIME": "2020-01-01T01:00:00+09:00",
			},
			Time: time.Date(2020, 1, 1, 1, 0, 0, 0, tokyo),
		},
		{
			Env: map[string]string{
				"TIME": "2020-01-01T01:00:00+00:00",
			},
			Time: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range cases {
		for k, v := range tc.Env {
			os.Setenv(k, v)
		}
		var conf config
		if err := envconfig.Process("process", &conf); err != nil {
			t.Fatal(err)
		}
		if !conf.Time.ToTime().Equal(tc.Time) {
			t.Fatal("time should be equal")
		}
	}
}

func TestTime_DecodeWithProcess(t *testing.T) {
	t.Run("env key not found : struct", func(t *testing.T) {
		type config struct {
			Time Time `envconfig:"TIME_STRUCT"`
		}
		// os.Setenv("TIME_STRUCT", "2020-01-01T01:00:00+09:00")
		var conf config
		if err := envconfig.Process("process", &conf); err != nil {
			t.Fatal(err)
		}
		if !conf.Time.IsZero() {
			t.Fatal("time should be zero value")
		}
	})

	t.Run("env key not found : pointer", func(t *testing.T) {
		type config struct {
			Time *Time `envconfig:"TIME_POINTER"`
		}
		// os.Setenv("TIME_STRUCT", "2020-01-01T01:00:00+09:00")
		var conf config
		if err := envconfig.Process("process", &conf); err != nil {
			t.Fatal(err)
		}
		if !conf.Time.IsZero() {
			t.Fatal("time should be zero value")
		}
	})

	t.Run("time format invalid : struct", func(t *testing.T) {
		type config struct {
			Time Time `envconfig:"TIME_STRUCT"`
		}
		os.Setenv("TIME_STRUCT", "2020-01-01T01:00:00")
		var conf config
		if err := envconfig.Process("process", &conf); err == nil {
			t.Fatal("error should be nil")
		}
		if !conf.Time.IsZero() {
			t.Fatal("time should be zero value")
		}
	})

	t.Run("time format invalid : pointer", func(t *testing.T) {
		type config struct {
			Time *Time `envconfig:"TIME_POINTER"`
		}
		os.Setenv("TIME_POINTER", "2020-01-01T01:00:00")
		var conf config
		if err := envconfig.Process("process", &conf); err == nil {
			t.Fatal("error should be nil")
		}
		if !conf.Time.IsZero() {
			t.Fatal("time should be zero value")
		}
	})
}
