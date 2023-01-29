package main

import "time"

// timeoutStep is a concrete type with the step type embedded in it
// that implements the stepExecutor interface, for steps that require
// a time out feature in their steps such as requests, over the internet
type timeoutStep struct {
	step
	timeout time.Duration
}

// newTimeoutStep is the accompanying constructor for timeout step
// as with all step types
func newTimeoutStep(name, exe, message, proj string, args []string, timeout time.Duration) timeoutStep {
	var s timeoutStep

	s.step = newStep(name, exe, message, proj, args)
	s.timeout = timeout

	return s
}
