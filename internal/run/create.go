package run

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/JY8752/note-cli/internal/file"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

func CreateArticleFunc(timeFlag *bool, name *string) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		t := *timeFlag
		n := *name

		// create directory name
		var dirName string

		// set timestamp in directory name
		if t {
			dirName = time.Now().Format("2006-01-02")

			counter := 1
			for {
				if !file.Exist(dirName) {
					break
				}
				counter++
				dirName = time.Now().Format("2006-01-02") + "-" + strconv.Itoa(counter)
			}
		}

		// set specify directory name
		if n != "" {
			dirName = n
		}

		// random value since nothing was specified
		if !t && n == "" {
			dirName = uuid.NewString()
		}

		// mkdir
		if err := os.Mkdir(fmt.Sprintf("./%s", dirName), 0744); err != nil {
			return err
		}

		fmt.Printf("Create directory. %s\n", dirName)

		// create markdown file
		filePath := fmt.Sprintf("./%s/%s.md", dirName, dirName)
		if _, err := os.OpenFile(filePath, os.O_CREATE, 0644); err != nil {
			return err
		}

		fmt.Printf("Create file. %s.md\n", dirName)

		// create config.yaml
		configFilePath := fmt.Sprintf("./%s/%s", dirName, ConfigFile)
		if err := os.WriteFile(configFilePath, []byte("title: article title\nauthor: your name"), 0644); err != nil {
			return err
		}

		return nil
	}
}

//go:embed templates/*
var templateFiles embed.FS

const (
	// custom template file name
	CustomTemplateFile = "template.tmpl"
	// config file name
	ConfigFile = "config.yaml"
	// output file name
	DefaultOutputFileName = "output.png"
)

// Information on images to be generated
type Ogp struct {
	Title    string
	IconPath string
	Author   string
}

// config schema
type Config struct {
	Title  string `yaml:"title"`
	Author string `yaml:"author"`
}

func CreateImageFunc(templateNo *int16, iconPath, outputPath *string) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var (
			tmpl *template.Template
			err  error
		)

		if file.Exist(CustomTemplateFile) {
			// use custom template html
			tmpl, err = template.ParseFiles(CustomTemplateFile)
		} else {
			// use template html
			tmpl, err = template.ParseFS(templateFiles, fmt.Sprintf("templates/%d.tmpl", *templateNo))
		}

		if err != nil {
			return err
		}

		// if use icon, base64 encode icon
		var encoded string
		if *iconPath != "" {
			b, err := os.ReadFile(*iconPath)
			if err != nil {
				return err
			}
			encoded = base64.StdEncoding.EncodeToString(b)
		}

		// read config yaml
		var config Config
		b, err := os.ReadFile(ConfigFile)
		if err != nil {
			return err
		}
		yaml.Unmarshal(b, &config)

		var buf bytes.Buffer
		tmpl.Execute(&buf, Ogp{
			Title:    config.Title,
			Author:   config.Author,
			IconPath: encoded,
		})

		html := buf.String()

		// Open tabs in headless browser
		page, err := rod.New().MustConnect().Page(proto.TargetCreateTarget{})
		if err != nil {
			return err
		}

		// set template html
		if err = page.SetDocumentContent(html); err != nil {
			return err
		}

		// take screenshot
		img, err := page.MustWaitStable().Screenshot(true, &proto.PageCaptureScreenshot{
			Format: proto.PageCaptureScreenshotFormatPng,
			Clip: &proto.PageViewport{
				X:      0,
				Y:      0,
				Width:  1200,
				Height: 600,
				Scale:  1,
			},
			FromSurface: true,
		})

		if err != nil {
			return err
		}

		// output
		if *outputPath != "" {
			err = utils.OutputFile(*outputPath, img)
		} else {
			err = utils.OutputFile(DefaultOutputFileName, img)
		}

		if err != nil {
			return err
		}

		fmt.Println("Complete generate OGP image")

		return nil
	}
}
