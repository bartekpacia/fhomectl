package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/bartekpacia/fhome/api"
	"github.com/bartekpacia/fhome/env"
	"github.com/urfave/cli/v2"
)

var configCommand = cli.Command{
	Name:  "config",
	Usage: "Manage system configuration",
	Subcommands: []*cli.Command{
		{
			Name:  "list",
			Usage: "List all available objects",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "system",
					Usage: "Print obj set by the configurator app",
				},
				&cli.BoolFlag{
					Name:  "user",
					Usage: "Print config set by the configurator app",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("system") && c.Bool("user") {
					return fmt.Errorf("cannot use both --system and --user")
				}

				err := client.OpenCloudSession(e.Email, e.CloudPassword)
				if err != nil {
					return fmt.Errorf("failed to open client session: %v", err)
				}
				log.Println("opened client session")

				_, err = client.GetMyResources()
				if err != nil {
					return fmt.Errorf("failed to get my resources: %v", err)
				}
				log.Println("got my resources")

				err = client.OpenResourceSession(e.ResourcePassword)
				if err != nil {
					return fmt.Errorf("failed to open client to resource session: %v", err)
				}
				log.Println("opened client to resource session")

				sysConfig, err := client.GetSystemConfig()
				if err != nil {
					return fmt.Errorf("failed to get sysConfig: %v", err)
				}
				log.Println("got system config")

				userConfig, err := client.GetUserConfig()
				if err != nil {
					return fmt.Errorf("failed to get user config: %v", err)
				}
				log.Println("got user config")

				if c.Bool("system") {
					w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
					defer w.Flush()

					fmt.Fprintf(w, "id\tdt\tpreset\tstyle\tperm\tstep\tdesc\n")
					fmt.Fprintf(w, "___\t___\t___\t___\t___\t___\t___\n")

					cells := sysConfig.Response.MobileDisplayProperties.Cells
					for _, cell := range cells {
						fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", cell.ID, cell.DisplayType, cell.Preset, cell.Style, cell.Permission, cell.Step, cell.Desc)
					}
				} else if c.Bool("user") {
					panels := map[string]api.UserPanel{}
					for _, panel := range userConfig.Panels {
						panels[panel.ID] = panel
					}

					log.Printf("there are %d cells\n", len(userConfig.Cells))
					for _, cell := range userConfig.Cells {
						log.Printf("id: %3d, name: %s, icon: %s panels:", cell.ObjectID, cell.Name, cell.Icon)
						for _, pos := range cell.PositionInPanel {
							log.Printf(" %s", panels[pos.PanelID].Name)
						}
						log.Println()
					}

					log.Printf("there are %d panels\n", len(userConfig.Panels))
					for _, panel := range userConfig.Panels {
						log.Printf("id: %s, name: %s\n", panel.ID, panel.Name)
					}
				} else {
					config, err := api.MergeConfigs(userConfig, sysConfig)
					if err != nil {
						return fmt.Errorf("failed to merge configs: %v", err)
					}

					log.Printf("there are %d panels and %d cells\n", len(config.Panels), len(config.Cells()))
					for _, panel := range config.Panels {
						log.Printf("panel id: %s, name: %s", panel.ID, panel.Name)
						for _, cell := range panel.Cells {
							log.Printf("\tid: %d, name: %s", cell.ID, cell.Name)
						}
						log.Println()
					}
				}

				return nil
			},
		},
	},
}

var watchCommand = cli.Command{
	Name:  "watch",
	Usage: "Print all incoming messages",
	Action: func(c *cli.Context) error {
		err := client.OpenCloudSession(e.Email, e.CloudPassword)
		if err != nil {
			return fmt.Errorf("failed to open client session: %v", err)
		}

		log.Println("successfully opened client session")

		_, err = client.GetMyResources()
		if err != nil {
			return fmt.Errorf("failed to get my resources: %v", err)
		}

		log.Println("successfully got my resources")

		err = client.OpenResourceSession(e.ResourcePassword)
		if err != nil {
			return fmt.Errorf("failed to open client to resource session: %v", err)
		}

		log.Println("successfully opened client to resource session")

		for {
			msg, err := client.ReadAnyMessage()
			if err != nil {
				return fmt.Errorf("failed to listen: %v", err)
			}

			if msg.ActionName == api.ActionStatusTouchesChanged {
				var touches api.StatusTouchesChangedResponse
				err = json.Unmarshal(msg.Raw, &touches)
				if err != nil {
					return fmt.Errorf("failed to unmarshal touches: %v", err)
				}

				log.Printf("%s\n", api.Pprint(touches))
			}

		}
	},
}

var toggleCommand = cli.Command{
	Name:  "toggle",
	Usage: "Toggle object's state (on/off)",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "object-id",
			Aliases:  []string{"id"},
			Value:    "",
			Usage:    "id of object to toggle",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		objectID := c.Int("object-id")

		err := client.OpenCloudSession(e.Email, e.CloudPassword)
		if err != nil {
			return fmt.Errorf("failed to open client session: %v", err)
		}

		log.Println("successfully opened client session")

		_, err = client.GetMyResources()
		if err != nil {
			return fmt.Errorf("failed to get my resources: %v", err)
		}

		log.Println("successfully got my resources")

		err = client.OpenResourceSession(e.ResourcePassword)
		if err != nil {
			return fmt.Errorf("failed to open client to resource session: %v", err)
		}

		log.Println("successfully opened client to resource session")

		err = client.SendEvent(objectID, api.ValueToggle)
		if err != nil {
			return fmt.Errorf("failed to send xevent to object with id %d: %v", objectID, err)
		}

		log.Println("successfully sent xevent to object with id", objectID)

		return nil
	},
}

var setCommand = cli.Command{
	Name:  "set",
	Usage: "Set object's value (0-100)",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "object-id",
			Aliases:  []string{"id"},
			Value:    "",
			Usage:    "id of object to toggle",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "value",
			Aliases:  []string{"val"},
			Usage:    "value (0-100)",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		objectID := c.Int("object-id")
		value := c.Int("value")

		err := client.OpenCloudSession(e.Email, e.CloudPassword)
		if err != nil {
			return fmt.Errorf("failed to open client session: %v", err)
		}

		log.Println("successfully opened client session")

		_, err = client.GetMyResources()
		if err != nil {
			return fmt.Errorf("failed to get my resources: %v", err)
		}

		log.Println("successfully got my resources")

		err = client.OpenResourceSession(e.ResourcePassword)
		if err != nil {
			return fmt.Errorf("failed to open client to resource session: %v", err)
		}

		log.Println("successfully opened client to resource session")

		err = client.SendEvent(objectID, api.MapLighting(value))
		if err != nil {
			return fmt.Errorf("failed to send xevent to object with id %d: %v", objectID, err)
		}

		log.Println("successfully sent xevent to object with id", objectID)

		return nil
	},
}

var (
	client *api.Client
	e      env.Env
)

func init() {
	log.SetFlags(0)
	var err error

	client, err = api.NewClient()
	if err != nil {
		log.Fatalf("failed to create api api client: %v\n", err)
	}

	e = env.Env{}
	err = e.Load()
	if err != nil {
		log.Fatalf("failed to load env: %v\n", err)
	}
}

func main() {
	app := &cli.App{
		Name:  "fh",
		Usage: "Interact with smart home devices connected to F&Home",
		Commands: []*cli.Command{
			&configCommand,
			&watchCommand,
			&toggleCommand,
			&setCommand,
		},
		CommandNotFound: func(c *cli.Context, command string) {
			log.Printf("invalid command '%s'. See 'fh --help'\n", command)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
