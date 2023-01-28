package main

import "os/exec"

// step represents the information of the CI step that is currently executing
type step struct {
	name    string
	exe     string
	args    []string
	message string
	proj    string
}

// newStep is a constructor function for step
func newStep(name, exe, message, proj string, args []string) step {
	return step{
		name:    name,
		exe:     exe,
		args:    args,
		message: message,
		proj:    proj,
	}
}

// execute is a step method that will run the command defined in the step
func (s step) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}
