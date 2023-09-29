package run

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"

	"github.com/JY8752/note-cli/internal/clock"
	"github.com/JY8752/note-cli/internal/file"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

type Options struct {
	BasePath       string
	DefaultDirName string
}

type Option func(*Options)

func CreateArticleFunc(timeFlag, noDirFlag *bool, name, author *string, options ...Option) RunEFunc {
	return func(cmd *cobra.Command, args []string) (err error) {
		t := *timeFlag
		n := *name

		// set option
		var op Options
		for _, option := range options {
			option(&op)
		}

		// create target name
		var targetName string

		// set timestamp in targert name
		now := clock.Now().Format("2006-01-02")
		if t {
			targetName = now

			counter := 1
			for {
				if !file.Exist(filepath.Join(op.BasePath, targetName)) {
					break
				}
				counter++
				targetName = now + "-" + strconv.Itoa(counter)
			}
		}

		// set specify name in target name
		if n != "" {
			if file.Exist(filepath.Join(op.BasePath, n)) && !*noDirFlag {
				return fmt.Errorf("already exist article directory. name: %s", n)
			}
			if file.Exist(filepath.Join(op.BasePath, n+".md")) && *noDirFlag {
				return fmt.Errorf("already exist article file. name: %s", n)
			}
			targetName = n
		}

		// random value since nothing was specified
		if !t && n == "" {
			if op.DefaultDirName != "" {
				targetName = op.DefaultDirName
			} else {
				targetName = uuid.NewString()
			}
		}

		// mkdir
		targetDir := op.BasePath
		if !*noDirFlag {
			targetDir = filepath.Join(targetDir, targetName)
			if err = os.Mkdir(targetDir, 0744); err != nil {
				return
			}

			fmt.Printf("Create directory. %s\n", targetDir)
		}

		// create markdown file
		filePath := filepath.Join(targetDir, fmt.Sprintf("%s.md", targetName))

		metadata := fmt.Sprintf(`---
title: ""
tags: []
date: "%s"
author: "%s"
---
`, now, *author)

		if err = os.WriteFile(filePath, []byte(metadata), 0644); err != nil {
			return
		}

		fmt.Printf("Create file. %s\n", filePath)

		return nil
	}
}

//go:embed templates/*
var templateFiles embed.FS

const (
	// custom template file name
	CustomTemplateFile = "template.tmpl"
	// config file
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
	Title  string   `yaml:"title"`
	Tags   []string `yaml:"tags"`
	Date   string   `yaml:"date"`
	Author string   `yaml:"author"`
}

func CreateImageFunc(templateNo *int16, iconPath, outputPath *string, options ...Option) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var (
			tmpl *template.Template
			err  error
		)

		var op Options
		for _, option := range options {
			option(&op)
		}

		if file.Exist(filepath.Join(op.BasePath, CustomTemplateFile)) {
			// use custom template html
			tmpl, err = template.ParseFiles(filepath.Join(op.BasePath, CustomTemplateFile))
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

		var config Config

		// read markdown file from current directry
		files, err := os.ReadDir(".")
		if err != nil {
			return err
		}

		var b []byte
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".md" {
				markdownPath := filepath.Join(op.BasePath, f.Name())
				markdownFile, err := os.Open(markdownPath)
				if err != nil {
					return err
				}
				defer markdownFile.Close()

				// extract metadata config
				if b, err = file.Extract(markdownFile, "---", "---"); err == nil {
					break
				}
			}
		}

		if len(b) == 0 {
			// If the config could not be loaded, look for config.yaml.
			// config.yaml is already obsolete and should be loaded for compatibility.
			if b, err = os.ReadFile(filepath.Join(op.BasePath, ConfigFile)); err != nil {
				return fmt.Errorf("not found article config. %w", err)
			}
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
			err = utils.OutputFile(filepath.Join(op.BasePath, DefaultOutputFileName), img)
		}

		if err != nil {
			return err
		}

		fmt.Println("Complete generate OGP image")

		return nil
	}
}
