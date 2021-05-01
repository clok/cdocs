package cdocs

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var (
	lastExitCode = 0
	fakeOsExiter = func(rc int) {
		lastExitCode = rc
	}
	fakeErrWriter = &bytes.Buffer{}
)

func init() {
	cli.OsExiter = fakeOsExiter
	cli.ErrWriter = fakeErrWriter
}

const multiline = `
View the diff of the local ansible-vault encrypted Kubenetes Secret file
against a given dotenv file on a pod within a namespace.

The local file will use the contents of the 'data.<accsessor flag>' block.
This defaults to 'data..env'.

Supported ansible-vault encryption version: $ANSIBLE_VAULT;1.1;AES256

Example file structure of decrypted file:

---
apiVersion: v1
kind: Secret
type: Opaque
data:
  .env: <BASE64 ENCODED STRING>

It will then grab contents of the dotenv filr on a Pod in a given Namespace.

This defaults to inspecting the '$PWD/.env on' when executing a 'cat' command.
This method uses '/bin/bash -c' as the base command to perform inspection.
`

const singleline = "A single line of UsageText"

func stub(c *cli.Context) error {
	fmt.Println("test")
	return nil
}

func testApp() *cli.App {
	app := &cli.App{
		Name:     "test-app",
		Version:  "v1.0.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "A bc",
				Email: "a@b.c",
			},
			{
				Name: "Just a Name",
			},
		},
		Copyright:            "(c) 2020 Yolo",
		HelpName:             "test-app",
		Usage:                "interact with config map and secret manager variables",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "Commands to interact with environment variables, both local and on cluster.",
				Subcommands: []*cli.Command{
					{
						Name:    "diff",
						Aliases: []string{"d"},
						Usage:   "Print out detailed diff reports comparing local and running Pod",
						Subcommands: []*cli.Command{
							{
								Name:      "namespace",
								Aliases:   []string{"ns"},
								Usage:     "View diff of local vs. namespace",
								UsageText: multiline,
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "secrets",
										Aliases:  []string{"s"},
										Usage:    "Path to secrets.yml",
										Required: false,
										Value:    ".docker/secrets.yml",
									},
									&cli.StringFlag{
										Name:     "configmap",
										Aliases:  []string{"c"},
										Usage:    "Path to configmap.yaml",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "namespace",
										Aliases:  []string{"n"},
										Usage:    "Kube Namespace to list Pods from for inspection",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "cmd",
										Usage:    "Command to inspect",
										Required: false,
										Value:    "node",
									},
									&cli.StringFlag{
										Name:     "filter-prefix",
										Aliases:  []string{"f"},
										Usage:    "List of prefixes (csv) used to filter values from display. Set to `\"\"` to remove any filters.",
										Required: false,
										Value:    "npm_,KUBERNETES_,API_PORT",
									},
									&cli.StringFlag{
										Name:     "exclude",
										Usage:    "List (csv) of specific env vars to exclude values from display. Set to `\"\"` to remove any exclusions.",
										Required: false,
										Value:    "PATH,SHLVL,HOSTNAME",
									},
									&cli.StringFlag{
										Name:  "secret-suffix",
										Usage: "Suffix used to find ENV variables that denote the Secret Manager Secrets to lookup",
										Value: "_NAME",
									},
								},
								Action: stub,
							},
							{
								Name:      "ansible",
								Aliases:   []string{"legacy"},
								Usage:     "View diff of local (ansible encrypted) vs. namespace",
								UsageText: multiline,
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "vault-password-file",
										Usage:    "vault password file `VAULT_PASSWORD_FILE`",
										Required: false,
									},
									&cli.StringFlag{
										Name:     "encrypted-env-file",
										Aliases:  []string{"e"},
										Usage:    "Path to encrypted Kube Secret file",
										Required: true,
									},
									&cli.StringFlag{
										Name:    "accessor",
										Aliases: []string{"a"},
										Usage:   "Accessor key to pull data out of Data block.",
										Value:   ".env",
									},
									&cli.StringFlag{
										Name:     "namespace",
										Aliases:  []string{"n"},
										Usage:    "Kube Namespace list Pods from for inspection",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "dotenv",
										Usage:    "Path to `.env` file on Pod",
										Required: false,
										Value:    "$PWD/.env",
									},
								},
								Action: stub,
							},
						},
					},
					{
						Name:    "view",
						Aliases: []string{"v"},
						Usage:   "View configured environment for either local or running on a Pod",
						Subcommands: []*cli.Command{
							{
								Name:      "configmap",
								Aliases:   []string{"c"},
								Usage:     "View env values based on local settings in a ConfigMap and secrets.yml",
								UsageText: singleline,
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "secrets",
										Aliases:  []string{"s"},
										Usage:    "Path to secrets.yml",
										Required: false,
										Value:    ".docker/secrets.yml",
									},
									&cli.StringFlag{
										Name:     "configmap",
										Aliases:  []string{"c"},
										Usage:    "Path to configmap.yaml",
										Required: true,
									},
									&cli.StringFlag{
										Name:  "secret-suffix",
										Usage: "Suffix used to find ENV variables that denote the Secret Manager Secrets to lookup",
										Value: "_NAME",
									},
								},
								Action: stub,
							},
							{
								Name:      "ansible",
								Aliases:   []string{"legacy"},
								Usage:     "View env values from ansible-vault encrypted Secret file.",
								UsageText: singleline,
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "vault-password-file",
										Usage:    "vault password file `VAULT_PASSWORD_FILE`",
										Required: false,
									},
									&cli.StringFlag{
										Name:     "encrypted-env-file",
										Aliases:  []string{"e"},
										Usage:    "Path to encrypted Kube Secret file",
										Required: true,
									},
									&cli.StringFlag{
										Name:    "accessor",
										Aliases: []string{"a"},
										Usage:   "Accessor key to pull data out of Data block.",
										Value:   ".env",
									},
								},
								Action: stub,
							},
							{
								Name:      "namespace",
								Aliases:   []string{"ns"},
								Usage:     "Interact with env on a running Pod within a Namespace",
								UsageText: multiline,
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "namespace",
										Aliases:  []string{"n"},
										Usage:    "Kube Namespace list Pods from",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "cmd",
										Usage:    "Command to inspect",
										Required: false,
										Value:    "node",
									},
									&cli.StringFlag{
										Name:     "filter-prefix",
										Aliases:  []string{"f"},
										Usage:    "List of prefixes (csv) used to filter values from display. Set to `\"\"` to remove any filters.",
										Required: false,
										Value:    "npm_,KUBERNETES_,API_PORT",
									},
									&cli.StringFlag{
										Name:     "exclude",
										Usage:    "List (csv) of specific env vars to exclude values from display. Set to `\"\"` to remove any exclusions.",
										Required: false,
										Value:    "PATH,SHLVL,HOSTNAME",
									},
								},
								Action: stub,
							},
						},
					},
				},
			},
			{
				Name:  "s3",
				Usage: "simple S3 commands",
				Subcommands: []*cli.Command{
					{
						Name:  "get",
						Usage: "[object path] [destination path]",
						Action: func(c *cli.Context) error {
							fmt.Println("test")
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
					fmt.Println("test")
					return nil
				},
			},
			{
				Name:   "hidden",
				Usage:  "This is hidden",
				Hidden: true,
				Action: func(c *cli.Context) error {
					fmt.Println("test")
					return nil
				},
			},
		},
	}

	return app
}

func Test_ToMarkdown(t *testing.T) {
	is := assert.New(t)

	app := testApp()

	res, err := ToMarkdown(app)

	data, _ := ioutil.ReadFile("data/test.md")

	is.Nil(err)
	is.Equal(res, string(data))
}

func Test_ToMan(t *testing.T) {
	is := assert.New(t)

	app := testApp()

	res, err := ToMan(app)

	data, _ := ioutil.ReadFile("data/test.man")

	is.Nil(err)
	is.Equal(res, string(data))
}

func Test_extractManpageSettings(t *testing.T) {
	is := assert.New(t)

	t.Run("AppName is required", func(t *testing.T) {
		_, _, _, _, err := extractManpageSettings(&InstallManpageCommandInput{
			AppName: "",
			CmdName: "",
			Path:    "",
		})

		is.Error(err, "AppName is required. Options passed in: &cdocs.InstallManpageCommandInput{AppName:\"\", CmdName:\"\", Path:\"\"}")
	})

	t.Run("path, command and hidden defaults", func(t *testing.T) {
		name, cmdname, path, hidden, err := extractManpageSettings(&InstallManpageCommandInput{
			AppName: "test",
			CmdName: "",
			Path:    "",
		})

		is.Nil(err)
		is.Equal(name, "test")
		is.Equal(cmdname, "install-manpage")
		is.Equal(path, "/usr/local/share/man/man8")
		is.False(hidden)
	})

	t.Run("all can be set", func(t *testing.T) {
		name, cmdname, path, hidden, err := extractManpageSettings(&InstallManpageCommandInput{
			AppName: "test",
			CmdName: "test-cmd",
			Path:    "/a/path",
			Hidden:  true,
		})

		is.Nil(err)
		is.Equal(name, "test")
		is.Equal(cmdname, "test-cmd")
		is.Equal(path, "/a/path")
		is.True(hidden)
	})
}

func Test_InstallManpageCommand(t *testing.T) {
	is := assert.New(t)

	t.Run("AppName is required", func(t *testing.T) {
		_, err := InstallManpageCommand(&InstallManpageCommandInput{
			AppName: "",
			CmdName: "",
			Path:    "",
		})

		is.Error(err, "AppName is required. Options passed in: &cdocs.InstallManpageCommandInput{AppName:\"\", CmdName:\"\", Path:\"\"}")
	})

	t.Run("generate a cli.Command", func(t *testing.T) {
		cmd, err := InstallManpageCommand(&InstallManpageCommandInput{
			AppName: "test",
		})

		is.Nil(err)
		is.NotNil(cmd)
		is.IsType(cli.Command{}, *cmd)
		is.Equal("install-manpage", cmd.Name)
		is.Equal("Generate and install man page", cmd.Usage)
		is.Equal("NOTE: Windows is not supported", cmd.UsageText)
	})

	t.Run("generate a man page and install", func(t *testing.T) {
		cmd, err := InstallManpageCommand(&InstallManpageCommandInput{
			AppName: "test-install-manpage",
			Path:    "/tmp",
		})

		is.Nil(err)
		is.NotNil(cmd)

		app := testApp()
		app.Commands = append(app.Commands, cmd)

		err = app.Run([]string{"test-install-manpage", "install-manpage"})
		is.Nil(err)
		is.FileExists("/tmp/test-install-manpage.8")
		err = os.Remove("/tmp/test-install-manpage.8")
		is.Nil(err)
	})

	t.Run("fail to generate a man page when installing to a path that does not exist", func(t *testing.T) {
		cmd, err := InstallManpageCommand(&InstallManpageCommandInput{
			AppName: "test-install-manpage",
			Path:    "/b/a/d/p/a/t/h/",
		})

		is.Nil(err)
		is.NotNil(cmd)

		app := testApp()
		app.Commands = append(app.Commands, cmd)

		err = app.Run([]string{"test-install-manpage", "install-manpage"})

		is.NotNil(err)
		is.Error(err)
		is.Equal("Unable to install man page. /b/a/d/p/a/t/h/ does not exist", err.Error())
	})
}

func Example() {
	im, err := InstallManpageCommand(&InstallManpageCommandInput{
		AppName: "demo",
	})
	if err != nil {
		panic(err)
	}

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
			im,
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

	md, _ := ToMarkdown(app)
	fmt.Println(md)

	man, _ := ToMan(app)
	fmt.Println(man)
}
