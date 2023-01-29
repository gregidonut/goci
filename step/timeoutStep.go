package step

import (
	"context"
	"os/exec"
	"time"
)

// timeoutStep is a concrete type with the step type embedded in it
// that implements the stepExecutor interface, for steps that require
// a time out feature in their steps such as requests, over the internet
type timeoutStep struct {
	step
	timeout time.Duration
}

// NewTimeoutStep is the accompanying constructor for timeout step
// as with all step types
func NewTimeoutStep(name, exe, message, proj string, args []string, timeout time.Duration) timeoutStep {
	var s timeoutStep

	s.step = NewStep(name, exe, message, proj, args)
	s.timeout = timeout

	return s
}

// Execute is an implementation of the stepExecutor interface's requirement
// which with timeoutStep, uses the timeout field as timeout argument to a
// context.WithTimeout function.
func (s timeoutStep) Execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:  s.name,
				msg:   "failed time out",
				cause: context.DeadlineExceeded,
			}
		}

		// a different error other than context deadline
		// then return that error instead
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to Execute",
			cause: err,
		}
	}

	return s.message, nil
}
