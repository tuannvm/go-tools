package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "org, o",
			Usage:  "`organization` to generate tf config for",
			EnvVar: "GITHUB_ORGANIZATION",
		},
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "`token` to access GitHub API",
			EnvVar: "GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "error",
			Usage:  "Log level (panic, fatal, error, warn, info, or debug)",
			EnvVar: "LOG_LEVEL",
		},
	}
	app := cli.NewApp()
	app.Name = "github-tf"
	app.Usage = "Download GitHub teams to TF config"
	app.Action = run

	app.Version = fmt.Sprintf("0.1.%s", build)
	app.Author = "Honestbee DevOps"

	app.Flags = flags

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	logLevelString := c.String("log-level")
	logLevel, err := log.ParseLevel(logLevelString)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)

	gitHub := GitHub{
		Organization: c.String("org"),
		Token:        c.String("token"),
	}

	if gitHub.Organization == "" || gitHub.Token == "" {
		cli.ShowAppHelpAndExit(c, 1)
	}

	//gitHub.ListRepos()

	teams, err := gitHub.ListTeams()
	if err != nil {
		return err
	}
	//get first 5 for testing purposes
	for _, t := range teams[:5] {
		fmt.Printf("%v (ID: %v)\n",
			*t.Slug,
			*t.ID,
		)
		teamRoles, _ := gitHub.GetTeamRoles(t)
		for k, v := range teamRoles.UserRoles {
			fmt.Printf("\t%v: %v\n",
				k,
				v,
			)
		}
	}
	return nil
}
