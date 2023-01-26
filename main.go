package main

import (
	"fmt"
	"io"
	"os/exec"
)

func main() {

}

func run(proj string, w io.Writer) error {
	// proj is the name of the project dir to run the tool against,
	// this will be passed in as a flag argument when using the tool
	if proj == "" {
		return fmt.Errorf("project directory is required")
	}

	args := []string{"build", ".", "errors"}
	cmd := exec.Command("go", args...)

	cmd.Dir = proj
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("'go build failed %s", err)
	}

	_, err := fmt.Fprintln(w, "Go Build: SUCCESS")
	return err
}
