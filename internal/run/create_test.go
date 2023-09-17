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
			// t.Parallel()

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

			// assertion
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

// func copy(src, dst string) error {
// 	srcFile, err := os.Open(src)
// 	if err != nil {
// 		return fmt.Errorf("failed open src file. src: %s %w", src, err)
// 	}

// 	defer func() {
// 		err = srcFile.Close()
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}()

// 	dstFile, err := os.Create(dst)
// 	if err != nil {
// 		return fmt.Errorf("failed open dst file. dst: %s %w", dst, err)
// 	}

// 	defer func() {
// 		err = dstFile.Close()
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}()

// 	if _, err = io.Copy(dstFile, srcFile); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func loadImage(path string) (image.Image, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	img, _, err := image.Decode(file)
// 	return img, err
// }

// func imagesEqual(img1, img2 image.Image) bool {
// 	if img1.Bounds() != img2.Bounds() {
// 		return false
// 	}

// 	for y := 0; y < img1.Bounds().Dy(); y++ {
// 		for x := 0; x < img1.Bounds().Dx(); x++ {
// 			if img1.At(x, y) != img2.At(x, y) {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }

// func TestCreateImageFunc(t *testing.T) {
// 	t.Parallel()

// 	tmpDir := t.TempDir()
// 	tmpDir2 := t.TempDir()

// 	// tmpDir
// 	// set icon file
// 	if err := copy("../../testdata/icon.png", filepath.Join(tmpDir, "icon.png")); err != nil {
// 		t.Fatal(err)
// 	}

// 	// set config file
// 	if err := copy("../../testdata/config.yaml", filepath.Join(tmpDir, "config.yaml")); err != nil {
// 		t.Fatal(err)
// 	}

// 	// tmpDir2
// 	// set icon file
// 	if err := copy("../../testdata/icon.png", filepath.Join(tmpDir2, "icon.png")); err != nil {
// 		t.Fatal(err)
// 	}

// 	// set config file
// 	if err := copy("../../testdata/config.yaml", filepath.Join(tmpDir2, "config.yaml")); err != nil {
// 		t.Fatal(err)
// 	}

// 	// set custom template file
// 	if err := copy("../../testdata/template.tmpl", filepath.Join(tmpDir2, "template.tmpl")); err != nil {
// 		t.Fatal(err)
// 	}

// 	testcases := []struct {
// 		name          string
// 		templateNo    int16
// 		iconPath      string
// 		outputPath    string
// 		useTmpDir     string
// 		actImagePath  string
// 		wantImagePath string
// 	}{
// 		{
// 			name:          "default",
// 			templateNo:    1,
// 			useTmpDir:     tmpDir,
// 			outputPath:    filepath.Join(tmpDir, "output1.png"),
// 			actImagePath:  filepath.Join(tmpDir, "output1.png"),
// 			wantImagePath: "../../testdata/output1.png",
// 		},
// 		{
// 			name:          "use icon flag",
// 			templateNo:    1,
// 			useTmpDir:     tmpDir,
// 			iconPath:      "../../testdata/icon.png",
// 			outputPath:    filepath.Join(tmpDir, "output2.png"),
// 			actImagePath:  filepath.Join(tmpDir, "output2.png"),
// 			wantImagePath: "../../testdata/output3.png",
// 		},
// 		{
// 			name:          "use output flag",
// 			templateNo:    1,
// 			useTmpDir:     tmpDir,
// 			outputPath:    filepath.Join(tmpDir, "output", "output.png"),
// 			actImagePath:  filepath.Join(tmpDir, "output", "output.png"),
// 			wantImagePath: "../../testdata/output1.png",
// 		},
// 		{
// 			name:          "use custom template",
// 			useTmpDir:     tmpDir2,
// 			actImagePath:  filepath.Join(tmpDir2, "output.png"),
// 			wantImagePath: "../../testdata/output2.png",
// 		},
// 		{
// 			name:          "use template2",
// 			templateNo:    2,
// 			useTmpDir:     tmpDir,
// 			outputPath:    filepath.Join(tmpDir, "output3.png"),
// 			actImagePath:  filepath.Join(tmpDir, "output3.png"),
// 			wantImagePath: "../../testdata/output4.png",
// 		},
// 		{
// 			name:          "use template3",
// 			templateNo:    3,
// 			useTmpDir:     tmpDir,
// 			outputPath:    filepath.Join(tmpDir, "output4.png"),
// 			actImagePath:  filepath.Join(tmpDir, "output4.png"),
// 			wantImagePath: "../../testdata/output5.png",
// 		},
// 	}

// 	for _, testcase := range testcases {
// 		testcase := testcase
// 		t.Run(testcase.name, func(t *testing.T) {
// 			// t.Parallel()
// 			err := run.CreateImageFunc(
// 				&testcase.templateNo,
// 				&testcase.iconPath,
// 				&testcase.outputPath,
// 				func(o *run.Options) { o.BasePath = testcase.useTmpDir },
// 			)(nil, nil)

// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// act, err := os.ReadFile(testcase.actImagePath)
// 			act, err := loadImage(testcase.actImagePath)
// 			if err != nil {
// 				t.Errorf("failed open act image file %v", err)
// 			}

// 			// want, err := os.ReadFile(testcase.wantImagePath)
// 			want, err := loadImage(testcase.wantImagePath)
// 			if err != nil {
// 				t.Errorf("failed open want image file %v", err)
// 			}

// 			// if !bytes.Equal(act, want) {
// 			// 	t.Errorf("not equal actPath: %s wantPath: %s", testcase.actImagePath, testcase.wantImagePath)
// 			// }
// 			if !imagesEqual(act, want) {
// 				t.Errorf("not equal actPath: %s wantPath: %s", testcase.actImagePath, testcase.wantImagePath)
// 			}
// 		})
// 	}
// }
