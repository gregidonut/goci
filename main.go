package main

import (
	"flag"
	"fmt"
	"github.com/gregidonut/goci/step"
	"io"
	"os"
)

func main() {
	// parsing the only flag tool accepts
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(proj string, w io.Writer) error {
	// proj is the name of the project dir to run the tool against,
	// this will be passed in as a flag argument when using the tool
	if proj == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}

	pipeline := make([]step.Executable, 3)

	pipeline[0] = step.NewStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)
	pipeline[1] = step.NewStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		proj,
		[]string{"test", "-v"},
	)
	pipeline[2] = step.NewExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		proj,
		[]string{"-l", "."},
	)

	for _, s := range pipeline {
		msg, err := s.Execute()
		if err != nil {
			return err
		}

		if _, err = fmt.Fprintln(w, msg); err != nil {
			return err
		}
	}

	return nil
}
