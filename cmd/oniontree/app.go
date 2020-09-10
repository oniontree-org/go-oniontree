package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/onionltd/go-oniontree"
	"github.com/urfave/cli/v2"
	"io/ioutil"
)

const Version = "0.1"

type Application struct {
	ot  *oniontree.OnionTree
	app *cli.App
}

func (a *Application) handleOnionTreeNew() cli.BeforeFunc {
	return func(c *cli.Context) error {
		a.ot = oniontree.New(c.String("C"))
		return nil
	}
}

func (a *Application) handleOnionTreeOpen() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ot, err := oniontree.Open(c.String("C"))
		if err != nil {
			return fmt.Errorf("failed to open OnionTree repository: %s", err)
		}
		a.ot = ot
		return nil
	}
}

func (a *Application) handleInitCommand() cli.ActionFunc {
	return func(c *cli.Context) error {
		return a.ot.Init()
	}
}

func (a *Application) handleAddCommand() cli.ActionFunc {
	return func(c *cli.Context) error {
		id := c.Args().First()
		if id == "" {
			return fmt.Errorf("Missing a service ID")
		}

		service := oniontree.NewService(id)
		service.Name = c.String("name")
		service.Description = c.String("description")
		service.AddURLs(c.StringSlice("url"))

		files := c.StringSlice("public-key")
		publicKeys := make([]*oniontree.PublicKey, 0, len(files))
		for i := range files {
			b, err := ioutil.ReadFile(files[i])
			if err != nil {
				return fmt.Errorf("failed to read public key content: %s", err)
			}

			publicKey, err := oniontree.NewPublicKey(b)
			if err != nil {
				return fmt.Errorf("failed to process public key content: %s", err)
			}

			publicKeys = append(publicKeys, publicKey)
		}
		service.AddPublicKeys(publicKeys)

		if err := a.ot.AddService(service); err != nil {
			return fmt.Errorf("failed to add new service: %s", err)
		}

		return nil
	}
}

func (a *Application) handleUpdateCommand() cli.ActionFunc {
	return func(c *cli.Context) error {
		id := c.Args().First()
		if id == "" {
			return cli.Exit("Missing a service ID", 1)
		}

		service, err := a.ot.GetService(id)
		if err != nil {
			return fmt.Errorf("failed to read service content: %s", err)
		}

		changed := false
		name := c.String("name")
		description := c.String("description")

		if name != "" && name != service.Name {
			service.Name = name
			changed = true
		}
		if description != "" && description != service.Description {
			service.Description = description
			changed = true
		}

		replace := c.Bool("replace")

		addedURLs := 0
		if replace {
			addedURLs = service.SetURLs(c.StringSlice("url"))
		} else {
			addedURLs = service.AddURLs(c.StringSlice("url"))
		}
		if addedURLs > 0 {
			changed = true
		}

		files := c.StringSlice("public-key")
		publicKeys := make([]*oniontree.PublicKey, 0, len(files))
		for i := range files {
			b, err := ioutil.ReadFile(files[i])
			if err != nil {
				return fmt.Errorf("failed to read public key content: %s", err)
			}

			publicKey, err := oniontree.NewPublicKey(b)
			if err != nil {
				return fmt.Errorf("failed to process public key content: %s", err)
			}

			publicKeys = append(publicKeys, publicKey)
		}
		addedPublicKeys := 0
		if replace {
			addedPublicKeys = service.SetPublicKeys(publicKeys)
		} else {
			addedPublicKeys = service.AddPublicKeys(publicKeys)
		}
		if addedPublicKeys > 0 {
			changed = true
		}

		if changed {
			if err := a.ot.UpdateService(service); err != nil {
				return fmt.Errorf("failed to update service: %s", err)
			}
		}

		return nil
	}
}

func (a *Application) handleRemoveCommand() cli.ActionFunc {
	return func(c *cli.Context) error {
		ids := c.Args().Slice()

		if len(ids) == 0 {
			return fmt.Errorf("Missing service IDs")
		}

		ok := true
		for i := range ids {
			if err := a.ot.RemoveService(ids[i]); err != nil {
				ok = false
				fmt.Printf("%s: %s\n", ids[i], err)
			}
		}

		if !ok {
			return cli.Exit("", 1)
		}

		return nil
	}
}

func (a *Application) handleTagCommand() cli.ActionFunc {
	return func(c *cli.Context) error {
		ids := c.Args().Slice()

		if len(ids) == 0 {
			return fmt.Errorf("Missing service IDs")
		}

		replace := c.Bool("replace")

		tags := make([]oniontree.Tag, len(c.StringSlice("name")))
		for i, tag := range c.StringSlice("name") {
			tags[i] = oniontree.Tag(tag)
		}

		for i := range ids {
			if replace {
				oldTags, err := a.ot.ListServiceTags(ids[i])
				if err != nil {
					return fmt.Errorf("failed to get old tags: %s", err)
				}

				if err := a.ot.UntagService(ids[i], oldTags); err != nil {
					return fmt.Errorf("failed to remove old tags: %s", err)
				}
			}
			if err := a.ot.TagService(ids[i], tags); err != nil {
				return fmt.Errorf("failed to create new tags: %s", err)
			}
		}

		return nil
	}
}

func (a *Application) handleUntagCommand() cli.ActionFunc {
	return func(c *cli.Context) error {
		ids := c.Args().Slice()

		if len(ids) == 0 {
			return fmt.Errorf("Missing service IDs")
		}

		tags := make([]oniontree.Tag, len(c.StringSlice("name")))
		for i, tag := range c.StringSlice("name") {
			tags[i] = oniontree.Tag(tag)
		}

		for i := range ids {
			if err := a.ot.UntagService(ids[i], tags); err != nil {
				return fmt.Errorf("failed to remove tags: %s", err)
			}
		}

		return nil
	}
}

func (a *Application) handleLintCommand() cli.ActionFunc {
	return func(c *cli.Context) error {
		serviceIDs, err := a.ot.ListServices()
		if err != nil {
			return fmt.Errorf("failed to list services: %s", err)
		}
		tags, err := a.ot.ListTags()
		if err != nil {
			return fmt.Errorf("failed to list tags: %s", err)
		}

		ok := true
		for i := range serviceIDs {
			service, err := a.ot.GetService(serviceIDs[i])
			if err != nil {
				return fmt.Errorf("failed to read service content: %s", err)
			}

			if err := service.Validate(); err != nil {
				ok = false
				fmt.Printf("unsorted/%s: %s\n", service.ID(), err)
			}
		}
		for i := range tags {
			if err := tags[i].Validate(); err != nil {
				ok = false
				fmt.Printf("tagged/%s: %s\n", tags[i], err)
			}
		}

		if !ok {
			return cli.Exit("", 1)
		}

		return nil
	}
}

func (a *Application) handleShowCommand() cli.ActionFunc {
	printYAML := func(s *oniontree.Service) {
		b, err := yaml.Marshal(s)
		if err != nil {
			return
		}
		fmt.Printf("%s\n", string(b))
	}
	printJSON := func(s *oniontree.Service) {
		b, err := json.Marshal(s)
		if err != nil {
			return
		}
		fmt.Printf("%s\n", string(b))
	}
	return func(c *cli.Context) error {
		id := c.Args().First()
		if id == "" {
			return cli.Exit("Missing a service ID", 1)
		}

		service, err := a.ot.GetService(id)
		if err != nil {
			return fmt.Errorf("failed to read service content: %s", err)
		}

		if c.Bool("json") {
			printJSON(service)
		} else {
			printYAML(service)
		}

		return nil
	}
}

func (a *Application) Run(args []string) error {
	return a.app.Run(args)
}
