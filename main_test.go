package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gregidonut/goci/step"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name     string
		proj     string
		wantOut  string
		wantErr  error
		setupGit bool
	}{
		{
			name: "Success",
			proj: "./testdata/tool/",
			wantOut: `Go Build: SUCCESS
Go Test: SUCCESS
Gofmt: SUCCESS
Git Push: SUCCESS
`,
			wantErr:  nil,
			setupGit: true,
		},
		{
			name:     "Fail",
			proj:     "./testdata/toolErr/",
			wantOut:  "",
			wantErr:  step.NewStepErr("go build", "", nil),
			setupGit: false,
		},
		{
			name:     "FailFormat",
			proj:     "./testdata/toolFmtErr/",
			wantOut:  "",
			wantErr:  step.NewStepErr("go fmt", "", nil),
			setupGit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupGit {
				cleanup := setupGit(t, tt.proj)
				defer cleanup()
			}

			w := bytes.Buffer{}
			err := run(tt.proj, &w)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("want error: %q, got 'nil'", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("want error: %q, got %q", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %q", err)
			}
			want := w.String()
			if want != tt.wantOut {
				t.Errorf("want: %q, got %q", tt.wantOut, want)
			}
		})
	}
}

// setupGit is responsible for setting up the git repo in the
// project directory and also the git server in a temporary directory,
// every prerequisite that the project presumably has already done before
// running 'git push'
func setupGit(t *testing.T, proj string) func() {
	t.Helper()

	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}

	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}

	remoteURI := fmt.Sprintf("file://%s", tempDir)

	gitCmdList := []struct {
		args []string
		dir  string
		env  []string
	}{
		// initialize remote repository in tempDir
		{args: []string{"init", "--bare"}, dir: tempDir, env: nil},
		// initialize git repository on projPath
		{args: []string{"init"}, dir: projPath, env: nil},
		// register remote repository(in tempDir) as remoteURI in projPath's git repo
		{args: []string{"remote", "add", "origin", remoteURI}, dir: projPath, env: nil},
		// add files in projPath to projPath's git repo's index
		{args: []string{"add", "."}, dir: projPath, env: nil},
		// commit staged in projPath repo with following environment variables
		{args: []string{"commit", "-m", "test"}, dir: projPath,
			env: []string{
				"GIT_COMMITTER_NAME=test",
				"GIT_COMMITTER_EMAIL=test@example.com",
				"GIT_AUTHOR_NAME=test",
				"GIT_AUTHOR_EMAIL=test@example.com",
			}},
	}

	for _, g := range gitCmdList {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir

		if g.env != nil {
			gitCmd.Env = append(os.Environ(), g.env...)
		}

		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		os.RemoveAll(tempDir)
		os.RemoveAll(filepath.Join(projPath, ".git"))
	}
}
