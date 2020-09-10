package main

import (
	"github.com/urfave/cli/v2"
)

func (a *Application) commands() {
	a.app = &cli.App{
		Name:    "oniontree",
		Version: Version,
		Usage:   "Manage OnionTree repository",
		Commands: cli.Commands{
			&cli.Command{
				Name:      "init",
				Usage:     "Initialize a new repository",
				ArgsUsage: " ",
				Before:    a.handleOnionTreeNew(),
				Action:    a.handleInitCommand(),
			},
			&cli.Command{
				Name:      "add",
				Usage:     "Add a new service to the repository",
				ArgsUsage: "<id>",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleAddCommand(),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Usage:    "service name",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "description",
						Usage: "service description (supports Markdown)",
					},
					&cli.StringSliceFlag{
						Name:     "url",
						Usage:    "service URL",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:  "public-key",
						Usage: "path to file with PGP public key",
					},
				},
			},
			&cli.Command{
				Name:      "update",
				Usage:     "Update a service",
				ArgsUsage: "<id>",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleUpdateCommand(),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "service name",
					},
					&cli.StringFlag{
						Name:  "description",
						Usage: "service description (supports Markdown)",
					},
					&cli.StringSliceFlag{
						Name:  "url",
						Usage: "service URL",
					},
					&cli.StringSliceFlag{
						Name:  "public-key",
						Usage: "path to file with PGP public key",
					},
					&cli.BoolFlag{
						Name:  "replace",
						Usage: "replace compound values",
					},
				},
			},
			&cli.Command{
				Name:      "show",
				Usage:     "Show service's content",
				ArgsUsage: "<id>",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleShowCommand(),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "json",
						Usage: "switch output to json",
					},
				},
			},
			&cli.Command{
				Name:      "remove",
				Usage:     "Remove services from the repository",
				ArgsUsage: "<id>[ id...]",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleRemoveCommand(),
			},
			&cli.Command{
				Name:      "tag",
				Usage:     "Tag services",
				ArgsUsage: "<id>[ id...]",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleTagCommand(),
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "name",
						Usage:    "tag name",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "replace",
						Usage: "replace tags",
					},
				},
			},
			&cli.Command{
				Name:      "untag",
				Usage:     "Untag services",
				ArgsUsage: "<id>[ id...]",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleUntagCommand(),
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "name",
						Usage:    "tag name",
						Required: true,
					},
				},
			},
			&cli.Command{
				Name:      "lint",
				Usage:     "Lint the repository content",
				ArgsUsage: " ",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleLintCommand(),
			},
		},
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "C",
				Value: ".",
				Usage: "change directory to",
			},
		},
	}
}
