package main

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
