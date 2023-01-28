package main

import (
	"bytes"
	"errors"
	"testing"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name    string
		proj    string
		wantOut string
		wantErr error
	}{
		{
			name:    "Success",
			proj:    "./testdata/tool/",
			wantOut: "Go Build: SUCCESS\n",
			wantErr: nil,
		},
		{
			name:    "Fail",
			proj:    "./testdata/toolErr/",
			wantOut: "",
			wantErr: &stepErr{step: "go build"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
