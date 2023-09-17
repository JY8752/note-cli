package file_test

import (
	"os"
	"testing"
	"time"

	"github.com/JY8752/note-cli/internal/file"
)

func removeDir(dir string) {
	for i := 0; i < 10; i++ {
		if err := os.RemoveAll(dir); err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestExist(t *testing.T) {
	tmpdir := t.TempDir()

	t.Cleanup(func() {
		removeDir(tmpdir)
	})

	if err := os.Mkdir(tmpdir+"/dir", 0777); err != nil {
		t.Fatal(err)
	}
	if _, err := os.OpenFile(tmpdir+"/dir/test.txt", os.O_CREATE, 0777); err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		name string
		path string
		want bool
	}{
		{
			"file is exist",
			tmpdir + "/dir/test.txt",
			true,
		},
		{
			"file is not exist",
			tmpdir + "/dir/test2.txt",
			false,
		},
		{
			"directory is exist",
			tmpdir + "/dir",
			true,
		},
		{
			"directory is not exist",
			tmpdir + "/dir2",
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
