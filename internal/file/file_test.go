package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JY8752/note-cli/internal/file"
)

func TestExist(t *testing.T) {
	tmpdir := t.TempDir()

	if err := os.Mkdir(filepath.Join(tmpdir, "dir"), 0777); err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile(filepath.Join(tmpdir, "dir", "test.txt"), os.O_CREATE, 0777)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	})

	testcases := []struct {
		name string
		path string
		want bool
	}{
		{
			"file is exist",
			filepath.Join(tmpdir, "dir", "test.txt"),
			true,
		},
		{
			"file is not exist",
			filepath.Join(tmpdir, "dir", "test2.txt"),
			false,
		},
		{
			"directory is exist",
			filepath.Join(tmpdir, "dir"),
			true,
		},
		{
			"directory is not exist",
			filepath.Join(tmpdir, "dir2"),
			false,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			if act := file.Exist(testcase.path); act != testcase.want {
				t.Errorf("File availability is not expected. want: %v act: %v\n", testcase.want, act)
			}
		})
	}
}
