package cdocs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/clok/kemba"
	"github.com/cpuguy83/go-md2man/v2/md2man"
	"github.com/urfave/cli/v2"
)

var (
	kl                  = kemba.New("cdocs")
	kim                 = kl.Extend("InstallManpageCommand")
	kman                = kl.Extend("install-manpage")
	markdownDocTemplate = `% {{ .App.Name }} 8
# NAME
{{ .App.Name }}{{ if .App.Usage }} - {{ .App.Usage }}{{ end }}
# SYNOPSIS
{{ .App.Name }}
{{ if .SynopsisArgs }}
` + "```" + `
{{ range $v := .SynopsisArgs }}{{ $v }}{{ end }}` + "```" + `
{{ end }}{{ if .App.UsageText }}
# DESCRIPTION
{{ .App.UsageText }}
{{ end }}

# COMMAND TREE
{{ range $v := .TOC }}
{{ $v }}{{ end }}

**Usage**:
` + "```" + `
{{ .App.Name }} [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
` + "```" + `
{{ if .GlobalArgs }}
# GLOBAL OPTIONS
{{ range $v := .GlobalArgs }}
{{ $v }}{{ end }}
{{ end }}{{ if .Commands }}
# COMMANDS
{{ range $v := .Commands }}
{{ $v }}{{ end }}{{ end }}`
)

// ToMarkdown creates a markdown string for the *cli.App
// The function errors if either parsing or writing of the string fails.
func ToMarkdown(a *cli.App) (string, error) {
	var w bytes.Buffer
	if err := writeDocTemplate(&w, a); err != nil {
		return "", err
	}
	return w.String(), nil
}

// ToMan creates a man page string for the *cli.App
// The function errors if either parsing or writing of the string fails.
func ToMan(a *cli.App) (string, error) {
	var w bytes.Buffer
	if err := writeDocTemplate(&w, a); err != nil {
		return "", err
	}
	man := md2man.Render(w.Bytes())
	return string(man), nil
}

// InstallManpageCommandInput provides an interface to pass in options for the InstallManpageCommand
//
// - AppName is required.
//
// - CmdName defaults to 'install-manpage'
//
// - Path defaults to '/usr/local/share/man/man8'
type InstallManpageCommandInput struct {
	AppName string `required:"true"`
	CmdName string `default:"install-command"`
	Path    string `default:"/usr/local/share/man/man8"`
	Hidden  bool   `default:"false"`
}

// InstallManpageCommand will generate a *cli.Command to be used with a cli.App.
// This will install a manual page (8) to the man-db.
func InstallManpageCommand(opts *InstallManpageCommandInput) (*cli.Command, error) {
	name, cmdname, path, hidden, err := extractManpageSettings(opts)
	if err != nil {
		return nil, err
	}

	cmd := &cli.Command{
		Name:      cmdname,
		Usage:     "Generate and install man page",
		UsageText: "NOTE: Windows is not supported",
		Hidden:    hidden,
		Action: func(c *cli.Context) error {
			kman.Printf("OS detected: %s", runtime.GOOS)
			if runtime.GOOS == "windows" {
				fmt.Println("Windows man page is not supported.")
				return nil
			}

			if _, err := os.Stat(path); os.IsNotExist(err) {
				return cli.Exit(fmt.Sprintf("Unable to install man page. %s does not exist", path), 2)
			}

			mp, _ := ToMan(c.App)
			manpath := filepath.Join(path, fmt.Sprintf("%s.8", name))
			kman.Printf("generated man page path: %s", manpath)
			err := os.WriteFile(manpath, []byte(mp), 0644)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Unable to install man page: %e", err), 2)
			}

			return nil
		},
	}

	return cmd, nil
}

// extractManpageSettings processes the *InstallManpageCommandInput and validates
func extractManpageSettings(opts *InstallManpageCommandInput) (string, string, string, bool, error) {
	kim.Printf("passed opts: %# v", opts)
	name := opts.AppName
	cmdname := opts.CmdName
	path := opts.Path
	hidden := opts.Hidden

	if name == "" {
		return "", "", "", hidden, fmt.Errorf("AppName is required. Options passed in: %# v", opts)
	}

	if path == "" {
		path = "/usr/local/share/man/man8"
	}

	if cmdname == "" {
		cmdname = "install-manpage"
	}
	kim.Printf("name: %s cmdname: %s path: %s hidden: %t", name, cmdname, path, hidden)
	return name, cmdname, path, hidden, nil
}

type cliTemplate struct {
	App          *cli.App
	TOC          []string
	Commands     []string
	GlobalArgs   []string
	SynopsisArgs []string
}

func writeDocTemplate(w io.Writer, a *cli.App) error {
	const name = "cli"
	t, err := template.New(name).Parse(markdownDocTemplate)
	if err != nil {
		return err
	}

	toc := generateCommandTree(a.Commands, 0)

	return t.ExecuteTemplate(w, name, &cliTemplate{
		App:          a,
		TOC:          toc,
		Commands:     prepareCommands(a.Commands, 0),
		GlobalArgs:   prepareArgsWithValues(a.VisibleFlags()),
		SynopsisArgs: prepareArgsSynopsis(a.VisibleFlags()),
	})
}

func generateCommandTree(commands []*cli.Command, level int) []string {
	var coms []string
	for _, command := range commands {
		if command.Hidden {
			continue
		}

		prepared := fmt.Sprintf("%s- [%s](#%s)", strings.Repeat("    ", level), strings.Join(command.Names(), ", "), strings.Join(command.Names(), "-"))

		coms = append(coms, prepared)

		// recursively iterate subcommands
		if len(command.Subcommands) > 0 {
			coms = append(
				coms,
				generateCommandTree(command.Subcommands, level+1)...,
			)
		}
	}

	return coms
}

func prepareCommands(commands []*cli.Command, level int) []string {
	var coms []string
	for _, command := range commands {
		if command.Hidden {
			continue
		}

		usageText := prepareUsageText(command)

		usage := prepareUsage(command, usageText)

		prepared := fmt.Sprintf("%s %s\n\n%s%s",
			strings.Repeat("#", level+2),
			strings.Join(command.Names(), ", "),
			usage,
			usageText,
		)

		flags := prepareArgsWithValues(command.Flags)
		if len(flags) > 0 {
			prepared += fmt.Sprintf("\n%s", strings.Join(flags, "\n"))
		}

		coms = append(coms, prepared)

		// recursevly iterate subcommands
		if len(command.Subcommands) > 0 {
			coms = append(
				coms,
				prepareCommands(command.Subcommands, level+1)...,
			)
		}
	}

	return coms
}

func prepareArgsWithValues(flags []cli.Flag) []string {
	return prepareFlags(flags, ", ", "**", "**", `""`, true)
}

func prepareArgsSynopsis(flags []cli.Flag) []string {
	return prepareFlags(flags, "|", "[", "]", "[value]", false)
}

func prepareFlags(
	flags []cli.Flag,
	sep, opener, closer, value string,
	addDetails bool,
) []string {
	args := []string{}
	for _, f := range flags {
		flag, ok := f.(cli.DocGenerationFlag)
		if !ok {
			continue
		}
		modifiedArg := opener

		for _, s := range flag.Names() {
			trimmed := strings.TrimSpace(s)
			if len(modifiedArg) > len(opener) {
				modifiedArg += sep
			}
			if len(trimmed) > 1 {
				modifiedArg += fmt.Sprintf("--%s", trimmed)
			} else {
				modifiedArg += fmt.Sprintf("-%s", trimmed)
			}
		}
		modifiedArg += closer
		if flag.TakesValue() {
			modifiedArg += fmt.Sprintf("=%s", value)
		}

		if addDetails {
			modifiedArg += flagDetails(flag)
		}

		args = append(args, modifiedArg+"\n")

	}
	sort.Strings(args)
	return args
}

// flagDetails returns a string containing the flags metadata
func flagDetails(flag cli.DocGenerationFlag) string {
	description := flag.GetUsage()
	value := flag.GetValue()
	if value != "" {
		description += " (default: " + value + ")"
	}
	return ": " + description
}

func prepareUsageText(command *cli.Command) string {
	usageText := ""
	if command.UsageText != "" {
		// Remove leading and trailing newlines
		preparedUsageText := strings.TrimSuffix(command.UsageText, "\n")
		preparedUsageText = strings.TrimPrefix(preparedUsageText, "\n")

		if strings.Contains(preparedUsageText, "\n") {
			// Format multi-line string as a code block
			usageText = fmt.Sprintf("```\n%s\n```\n", preparedUsageText)
		} else {
			// Style a single line as a note
			usageText = fmt.Sprintf(">%s\n", preparedUsageText)
		}
	}
	return usageText
}

func prepareUsage(command *cli.Command, usageText string) string {
	usage := ""
	if command.Usage != "" {
		usage = fmt.Sprintf("%s\n", command.Usage)
	}

	// Add a newline to the Usage IFF there is a UsageText
	if usageText != "" && usage != "" {
		usage = fmt.Sprintf("%s\n", usage)
	}

	return usage
}
