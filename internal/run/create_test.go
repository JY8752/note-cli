package run_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/JY8752/note-cli/internal/clock"
	"github.com/JY8752/note-cli/internal/file"
	"github.com/JY8752/note-cli/internal/run"
	"gopkg.in/yaml.v3"
)

func TestCreateArticleFunc(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	testcases := []struct {
		name                   string
		timeFlag               bool
		dirName                string
		want                   string
		wantErr                bool
		createDuplicateDirName string
		setTestTime            func() time.Time
	}{
		{
			"default case",
			false,
			"",
			"test",
			false,
			"",
			nil,
		},
		{
			"specify dirName",
			false,
			"article",
			"article",
			false,
			"",
			nil,
		},
		{
			"specify time flag",
			true,
			"",
			"2023-09-09",
			false,
			"",
			func() time.Time { return time.Date(2023, 9, 9, 12, 0, 0, 0, time.Local) },
		},
		{
			"both name and time flag", // Flags are grouped so that both flags are never done together.
			true,
			"article2",
			"article2",
			false,
			"",
			nil,
		},
		{
			"when duplicate directory exist, specify directory name",
			false,
			"article3",
			"",
			true,
			"article3",
			nil,
		},
		{
			"when duplicate directory exist, timestamp directory name",
			true,
			"",
			"2023-09-10-2",
			false,
			"2023-09-10",
			func() time.Time { return time.Date(2023, 9, 10, 12, 0, 0, 0, time.Local) },
		},
	}

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			// set test time
			if testcase.setTestTime != nil {
				clock.Now = testcase.setTestTime
			}

			// create duplicate dir
			if testcase.createDuplicateDirName != "" && !file.Exist(filepath.Join(tmpDir, testcase.createDuplicateDirName)) {
				os.Mkdir(filepath.Join(tmpDir, testcase.createDuplicateDirName), 0777)
			}

			// run
			err := run.CreateArticleFunc(
				&testcase.timeFlag,
				&testcase.dirName,
				func(o *run.Options) { o.BasePath = tmpDir },
				func(o *run.Options) { o.DefaultDirName = "test" },
			)(nil, nil)

			// check error
			if testcase.wantErr {
				if err == nil {
					t.Error("want error, but not error.\n")
				} else {
					return
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}

			// asertion
			if !file.Exist(filepath.Join(tmpDir, testcase.want)) {
				t.Errorf("directory is not exist. want: %s\n", testcase.want)
			}

			if !file.Exist(filepath.Join(tmpDir, testcase.want, testcase.want+".md")) {
				t.Errorf("markdown file is not exist. want: %s\n", filepath.Join(tmpDir, testcase.want, testcase.want+".md"))
			}

			configPath := filepath.Join(tmpDir, testcase.want, "config.yaml")
			if !file.Exist(configPath) {
				t.Errorf("config file is not exist. want: %s\n", configPath)
			}

			b, err := os.ReadFile(configPath)
			if err != nil {
				t.Error(err)
			}

			var config run.Config
			if err = yaml.Unmarshal(b, &config); err != nil {
				t.Error(err)
			}

			expectConfig := run.Config{Title: "article title", Author: "your name"}
			if !reflect.DeepEqual(config, expectConfig) {
				t.Errorf("config.yaml content are not expected. expect: %v act: %v\n", expectConfig, config)
			}
		})
	}
}
