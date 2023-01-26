package main

import (
	"fmt"
	"io"
)

func main() {

}

func run(proj string, w io.Writer) error {
	// proj is the name of the project dir to run the tool against,
	// this will be passed in as a flag argument when using the tool
	if proj == "" {
		return fmt.Errorf("project directory is required")
	}
	return nil
}
