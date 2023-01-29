package step

import (
	"os/exec"
)

type Executor interface {
	Execute() (string, error)
}

// step represents the information of the CI step that is currently executing
type step struct {
	name    string
	exe     string
	args    []string
	message string
	proj    string
}

// NewStep is a constructor function for step
func NewStep(name, exe, message, proj string, args []string) step {
	return step{
		name:    name,
		exe:     exe,
		args:    args,
		message: message,
		proj:    proj,
	}
}

// Execute is a step method that will run the command defined in the step
func (s step) Execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to Execute",
			cause: err,
		}
	}

	return s.message, nil
}
