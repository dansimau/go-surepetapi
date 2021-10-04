package cli

import (
	"fmt"
	"os"

	"github.com/dansimau/go-surepetapi"
	"github.com/dansimau/go-surepetapi/pkg/config"

	"github.com/jessevdk/go-flags"
)

type Cmd struct {
	ConfigPath string `short:"c" long:"config" description:"path to config"`
}

var (
	cmd    = &Cmd{}
	parser = flags.NewParser(cmd, flags.HelpFlag)
)

// Run executes the program with the specified arguments and returns the code
// the process should exit with.
func Run(args []string) (exitCode int) {
	_, err := parser.ParseArgs(args)
	if err != nil {
		// Handle --help, which is represented as an error by the flags package
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 0
		}

		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return 1
	}

	return 0
}

func (c *Cmd) api() (*surepetapi.Client, error) {
	cfg, err := c.config()
	if err != nil {
		return nil, err
	}

	return surepetapi.NewClient(cfg.API)
}

func (cmd *Cmd) config() (*config.Config, error) {
	var cfg *config.Config

	// use user-provided config, if specified
	if cmd.ConfigPath != "" {
		c, err := config.FromFile(cmd.ConfigPath)
		if err != nil {
			return nil, err
		}

		cfg = c
	} else {
		c, configPath, err := config.SearchAndLoadConfig("surepet.yaml")
		if err != nil {
			return nil, err
		}

		cfg = c
		cmd.ConfigPath = configPath
	}

	return cfg, nil
}
