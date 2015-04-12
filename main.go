package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/ulule/ipfix/application"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipfix"
	app.Author = "thoas"
	app.Email = "florent.messa@gmail.com"
	app.Usage = "A webservice to retrieve geolocation information from an ip address"
	app.Version = application.Version
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
				fmt.Printf("ipfix %s\n", application.Version)
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

		err := application.Run(config)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
