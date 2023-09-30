package file_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JY8752/note-cli/internal/file"
	"github.com/MakeNowJust/heredoc/v2"
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

func TestExtract(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		content string
		start   string
		end     string
		want    string
		wantErr bool
	}{
		"simple: double quotation": {
			content: `"test" xxxxxx`,
			start:   "\"",
			end:     "\"",
			want:    "test",
		},
		"default: from --- to ---": {
			content: heredoc.Doc(`---
				title: ""
				tags: []
				date: "2023-09-30"
				author: "Junichi.Y"
				---
			`),
			start: "---",
			end:   "---",
			want: "\n" + heredoc.Doc(`
				title: ""
				tags: []
				date: "2023-09-30"
				author: "Junichi.Y"
			`),
		},
		"error: missing start key": {
			content: heredoc.Doc(`title: ""
				tags: []
				date: "2023-09-30"
				author: "Junichi.Y"
				---
			`),
			start:   "---",
			end:     "---",
			wantErr: true,
		},
		"error: missing end key": {
			content: heredoc.Doc(`---
				title: ""
				tags: []
				date: "2023-09-30"
				author: "Junichi.Y"
			`),
			start:   "---",
			end:     "---",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.content)
			b, err := file.Extract(r, tt.start, tt.end)

			if tt.wantErr && err == nil {
				t.Fatal("expect error, but not error")
			}
			if tt.wantErr && err != nil {
				t.Log(err)
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(b, []byte(tt.want)) {
				t.Errorf("extract contents are not expected. act: %s exp: %s", string(b), tt.want)
			}
		})
	}
}
