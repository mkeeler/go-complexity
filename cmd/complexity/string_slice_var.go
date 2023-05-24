package main

import (
	"strings"
)

type StringSliceValue []string

func (s *StringSliceValue) String() string {
	return strings.Join(*s, ",")
}

func (s *StringSliceValue) Set(value string) error {
	if *s == nil {
		*s = make([]string, 0, 1)
	}

	*s = append(*s, value)
	return nil
}
