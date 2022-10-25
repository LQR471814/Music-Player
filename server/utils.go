package main

import (
	"reflect"
	"strings"
	"time"
)

type Number interface {
	int | int32 | int64 | float32 | float64
}

type Primitive interface {
	string | Number
}

func Fallback[T Primitive](values []T) T {
	var defaultValue T
	for _, value := range values {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String:
			if strings.Trim(any(value).(string), " ") != "" {
				return value
			}
		case reflect.Int:
			if any(value).(int) != 0 {
				return value
			}
		case reflect.Int32:
			if any(value).(int32) != 0 {
				return value
			}
		case reflect.Int64:
			if any(value).(int64) != 0 {
				return value
			}
		case reflect.Float32:
			if any(value).(float32) != 0 {
				return value
			}
		case reflect.Float64:
			if any(value).(float64) != 0 {
				return value
			}
		}
	}
	return defaultValue
}

func Sleep(d time.Duration) chan bool {
	channel := make(chan bool)
	go func() {
		time.Sleep(d)
		channel <- true
	}()
	return channel
}

func Remove[T any](slice []T, target func(T) bool) []T {
	values := []T{}
	for _, v := range slice {
		if target(v) {
			values = append(values, v)
		}
	}
	return values
}
