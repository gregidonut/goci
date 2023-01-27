package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
)

type stepErr struct {
	step  string
	msg   string
	cause error
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}
