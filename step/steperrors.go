package step

import "fmt"

type stepErr struct {
	step  string
	msg   string
	cause error
}

// NewStepErr is a constructor for stepErr useful for testing for now
func NewStepErr(step, msg string, cause error) *stepErr {
	return &stepErr{
		step:  step,
		msg:   msg,
		cause: cause,
	}
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

func (s *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}
	return t.step == s.step
}

func (s *stepErr) Unwrap() error {
	return s.cause
}
