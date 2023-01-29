package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// exceptionStep is a step type that has a different execute() method logic
// in specific cases where the output of the actual executing command's output
// needs to be handled, like the gofmt
type exceptionStep struct {
	step
}

func newExceptionStep(name, exe, message, proj string, args []string) exceptionStep {
	var s exceptionStep

	s.step = newStep(name, exe, message, proj, args)
	return s
}

func (s exceptionStep) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)

	w := &bytes.Buffer{}
	cmd.Stdout = w

	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	// if len of bytes in stdout is more than 0 then print the
	// message from stdout
	if w.Len() > 0 {
		return "", &stepErr{
			step:  s.name,
			msg:   fmt.Sprintf("invalid format: %s", w.String()),
			cause: nil,
		}
	}

	return s.message, nil
}
