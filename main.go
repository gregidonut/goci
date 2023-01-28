package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// parsing the only flag tool accepts
	proj := flag.String("p", "", "Project directory")
	flag.Parsed()

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

	pipeline := make([]step, 1)

	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)

	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		if _, err = fmt.Fprintln(w, msg); err != nil {
			return err
		}
	}

	return nil
}
