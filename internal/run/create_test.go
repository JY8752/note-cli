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

	clock.Now = func() time.Time { return time.Date(2023, 9, 9, 12, 0, 0, 0, time.Local) }
	nowStr := "2023-09-09"

	testcases := []struct {
		name                         string
		timeFlag                     bool
		noDirFlag                    bool
		targetName                   string
		author                       string
		want                         string
		wantErr                      bool
		createDuplicateDirOrFileName string
	}{
		{
			name:                         "default case",
			timeFlag:                     false,
			noDirFlag:                    false,
			targetName:                   "",
			author:                       "author",
			want:                         "test",
			wantErr:                      false,
			createDuplicateDirOrFileName: "",
		},
		{
			name:                         "specify dirName",
			timeFlag:                     false,
			noDirFlag:                    false,
			targetName:                   "article",
			author:                       "author",
			want:                         "article",
			wantErr:                      false,
			createDuplicateDirOrFileName: "",
		},
		{
			name:                         "specify time flag",
			timeFlag:                     true,
			noDirFlag:                    false,
			targetName:                   "",
			author:                       "author",
			want:                         "2023-09-09",
			wantErr:                      false,
			createDuplicateDirOrFileName: "",
		},
		{
			name:                         "both name and time flag", // Flags are grouped so that both flags are never done together.
			timeFlag:                     true,
			noDirFlag:                    false,
			targetName:                   "article2",
			author:                       "author",
			want:                         "article2",
			wantErr:                      false,
			createDuplicateDirOrFileName: "",
		},
		{
			name:                         "when duplicate directory exist, specify directory name",
			timeFlag:                     false,
			noDirFlag:                    false,
			targetName:                   "article3",
			author:                       "author",
			want:                         "",
			wantErr:                      true,
			createDuplicateDirOrFileName: "article3",
		},
		{
			name:                         "when duplicate directory exist, timestamp directory name",
			timeFlag:                     true,
			noDirFlag:                    false,
			targetName:                   "",
			author:                       "author",
			want:                         "2023-09-09-2",
			wantErr:                      false,
			createDuplicateDirOrFileName: "2023-09-09",
		},
		{
			name:                         "default, when no directory",
			timeFlag:                     false,
			noDirFlag:                    true,
			targetName:                   "",
			author:                       "author",
			want:                         "test",
			wantErr:                      false,
			createDuplicateDirOrFileName: "",
		},
		{
			name:                         "timestamp, when no directory",
			timeFlag:                     true,
			noDirFlag:                    true,
			targetName:                   "",
			author:                       "author",
			want:                         "2023-09-09",
			wantErr:                      false,
			createDuplicateDirOrFileName: "",
		},
		{
			name:                         "specific file name, when no directory",
			timeFlag:                     false,
			noDirFlag:                    true,
			targetName:                   "article",
			author:                       "author",
			want:                         "article",
			wantErr:                      false,
			createDuplicateDirOrFileName: "",
		},
	}

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()

			// create duplicate dir
			if testcase.createDuplicateDirOrFileName != "" {
				if testcase.noDirFlag {
					os.Create(filepath.Join(tmpDir, testcase.createDuplicateDirOrFileName))
				}
				os.Mkdir(filepath.Join(tmpDir, testcase.createDuplicateDirOrFileName), 0777)
			}

			// run
			err := run.CreateArticleFunc(
				&testcase.timeFlag,
				&testcase.noDirFlag,
				&testcase.targetName,
				&testcase.author,
				func(o *run.Options) { o.BasePath = tmpDir },
				func(o *run.Options) { o.DefaultDirName = "test" },
			)(nil, nil)

			// check error
			if testcase.wantErr {
				if err == nil {
					t.Fatal("want error, but not error.\n")
				} else {
					return
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
			}

			// assertion
			if !testcase.noDirFlag && !file.Exist(filepath.Join(tmpDir, testcase.want)) {
				t.Fatalf("directory is not exist. want: %s\n", testcase.want)
			}

			var markdownPath string
			if testcase.noDirFlag {
				markdownPath = filepath.Join(tmpDir, testcase.want+".md")
			} else {
				markdownPath = filepath.Join(tmpDir, testcase.want, testcase.want+".md")
			}
			if !file.Exist(markdownPath) {
				t.Fatalf("markdown file is not exist. want: %s\n", markdownPath)
			}

			b, err := os.ReadFile(markdownPath)
			if err != nil {
				t.Fatal(err)
			}

			var config run.Config
			if err = yaml.Unmarshal(b, &config); err != nil {
				t.Fatal(err)
			}

			expectConfig := run.Config{Date: nowStr, Author: testcase.author, Tags: []string{}}
			if !reflect.DeepEqual(config, expectConfig) {
				t.Errorf("config content are not expected. expect: %v act: %v\n", expectConfig, config)
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
