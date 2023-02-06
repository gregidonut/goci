package main

import (
	"flag"
	"fmt"
	"github.com/gregidonut/goci/step"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	pipeline := make([]step.Executable, 4)

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

	pipeline[3] = step.NewTimeoutStep(
		"git push",
		"git",
		"Git Push: SUCCESS",
		proj,
		[]string{"push", "origin", "main"},
		10*time.Second,
	)

	// prepare to handle at least one signal correctly
	sig := make(chan os.Signal, 1)

	// generic concurrency pattern for error propagation
	// and done channel to be used when business logic is
	// complete
	errCh := make(chan error)
	done := make(chan struct{})

	// pass to sig chanel variable a notification if a sigint(when
	// ctrl+c is pressed) or a sigterm(when the kill command is used)
	// is received, all other signals are ignored
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for _, s := range pipeline {
			msg, err := s.Execute()
			if err != nil {
				errCh <- err
				return
			}

			if _, err = fmt.Fprintln(w, msg); err != nil {
				errCh <- err
				return
			}
		}
		close(done)
	}()

	for {
		select {
		case rec := <-sig:
			signal.Stop(sig)
			return fmt.Errorf("%s: Exiting: %w", rec, ErrSignal)
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
}
