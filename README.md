# cdocs - urfave/cli/v2 docs extension

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/clok/cdocs/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/clok/cdocs)](https://goreportcard.com/report/clok/cdocs)
[![Coverage Status](https://coveralls.io/repos/github/clok/cdocs/badge.svg)](https://coveralls.io/github/clok/cdocs)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/clok/cdocs?tab=overview)

This is an enhanced version of the `ToMarkdown` and `ToMan` methods for [https://github.com/urfave/cli/v2](https://github.com/urfave/cli).

Key differences are:

- Addition of a Table of Contents with working markdown links.
- `UsageText` included in generated doc files.

Examples:
- [gwsm](https://github.com/GoodwayGroup/gwsm/blob/master/docs/gwsm.md)
- [gwvault](https://github.com/GoodwayGroup/gwvault/blob/master/docs/gwvault.md)
- [gw-aws-audit](https://github.com/GoodwayGroup/gw-aws-audit/blob/master/docs/gw-aws-audit.md)

## Usage

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/clok/cdocs"
    "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "demo",
		Version:  "0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "John Doe",
				Email: "j@doe.com",
			},
		},
		HelpName:             "demo",
		Usage:                "a demo cli app",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "s3",
				Usage: "simple S3 commands",
				Subcommands: []*cli.Command{
					{
						Name:      "get",
						Usage:     "[object path] [destination path]",
						UsageText: "it's going to get an object",
						Action: func(c *cli.Context) error {
							fmt.Println("get")
							return nil
						},
					},
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version info",
				Action: func(c *cli.Context) error {
					fmt.Println("version")
					return nil
				},
			},
		},
	}

	if os.Getenv("DOCS_MD") != "" {
		docs, err := cdocs.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := cdocs.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
```