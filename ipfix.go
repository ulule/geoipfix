package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

// compilation variables.
var (
	branch   string
	sha      string
	now      string
	compiler string
)

func main() {
	app := cli.NewApp()
	app.Name = "ipfix"
	app.Author = "thoas"
	app.Email = "florent.messa@gmail.com"
	app.Usage = "A webservice to retrieve geolocation information from an ip address"
	app.Version = fmt.Sprintf("%s [git:%s:%s]\ncompiled using %s at %s)", Version, branch, sha, compiler, now)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Config file path",
			EnvVar: "IPFIX_CONFIG_PATH",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "Retrieve the version number",
			Action: func(c *cli.Context) {
				fmt.Printf("ipfix %s\n", Version)
			},
		},
	}
	app.Action = func(c *cli.Context) {
		config := c.String("config")

		if config != "" {
			if _, err := os.Stat(config); err != nil {
				fmt.Fprintf(os.Stderr, "Can't find config file `%s`\n", config)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Can't find config file\n")
			os.Exit(1)
		}

		err := Run(config)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
