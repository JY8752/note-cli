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

		var dirName string

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

		if n != "" {
			dirName = n
		}

		if !t && n == "" {
			dirName = uuid.NewString()
		}

		if err := os.Mkdir(fmt.Sprintf("./%s", dirName), 0744); err != nil {
			return err
		}

		fmt.Printf("Create directory. %s\n", dirName)

		filePath := fmt.Sprintf("./%s/%s.md", dirName, dirName)
		if _, err := os.OpenFile(filePath, os.O_CREATE, 0644); err != nil {
			return err
		}

		fmt.Printf("Create file. %s.md\n", dirName)

		return nil
	}
}

//go:embed templates/*
var templateFiles embed.FS

const (
	CustomTemplateFile    = "template.tmpl"
	ConfigFile            = "config.yaml"
	DefaultOutputFileName = "output.png"
)

type Ogp struct {
	Title    string
	IconPath string
	Author   string
}

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
			tmpl, err = template.ParseFiles(CustomTemplateFile)
		} else {
			tmpl, err = template.ParseFS(templateFiles, fmt.Sprintf("templates/%d.tmpl", *templateNo))
		}

		if err != nil {
			return err
		}

		var encoded string
		if *iconPath != "" {
			b, err := os.ReadFile(*iconPath)
			if err != nil {
				return err
			}
			encoded = base64.StdEncoding.EncodeToString(b)
		}

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

		page, err := rod.New().MustConnect().Page(proto.TargetCreateTarget{})
		if err != nil {
			return err
		}

		if err = page.SetDocumentContent(html); err != nil {
			return err
		}

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
