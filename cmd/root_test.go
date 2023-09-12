package cmd

import (
	"testing"
)

func TestCommandFlag(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "article command1",
			args:    []string{"create", "article", "-t", "-n", "test"},
			wantErr: true,
		},
		{
			name:    "article command2",
			args:    []string{"create", "article", "-t", "--name", "test"},
			wantErr: true,
		},
		{
			name:    "article command3",
			args:    []string{"create", "article", "-time", "--name", "test"},
			wantErr: true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			rootCmd.SetArgs(testcase.args)
			err := rootCmd.Execute()

			if testcase.wantErr && err == nil {
				t.Error("want error, but not error")
			}

			if !testcase.wantErr && err != nil {
				t.Errorf("don't want error, but error err: %s", err)
			}
		})
	}
}
