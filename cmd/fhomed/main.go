package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bartekpacia/fhome/env"
	"github.com/bartekpacia/fhome/fhome"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

var (
	client *fhome.Client
	e      env.Env
	a      *accessory.Switch
)

func init() {
	log.SetFlags(0)

	var err error

	client, err = fhome.NewClient()
	if err != nil {
		log.Fatalf("failed to create fhome client: %v\n", err)
	}

	e = env.Env{}
	e.Load()
}

func main() {
	go setUpHap()

	err := client.OpenClientSession(e.Email, e.CloudPassword)
	if err != nil {
		log.Fatalf("failed to open client session: %v", err)
	}

	log.Println("successfully opened client session")

	_, err = client.GetMyResources()
	if err != nil {
		log.Fatalf("failed to get my resources: %v", err)
	}

	log.Println("successfully got my resources")

	err = client.OpenClientToResourceSession(e.ResourcePassword)
	if err != nil {
		log.Fatalf("failed to open client to resource session: %v", err)
	}

	log.Println("successfully opened client to resource session")

	responses := make(chan fhome.StatusTouchesChangedResponse)
	errors := make(chan error)

	go client.Listen(responses, errors)

	for {
		select {
		case msg := <-responses:
			if len(msg.Response.Cv) == 0 {
				continue
			}

			cellValue := msg.Response.Cv[0]
			// fmt.Println("got msg:", cellValue)
			if cellValue.Voi == "291" {
				if cellValue.Dvs == "100%" {
					log.Println("f&home ON")
					a.Switch.On.SetValue(true)
				} else {
					log.Println("f&home OFF")
					a.Switch.On.SetValue(false)
				}
			}
		case err := <-errors:
			log.Fatalf("failed to listen: %v", err)
		}
	}
}

func setUpHap() {
	// Create the switch accessory.
	a = accessory.NewSwitch(accessory.Info{
		Name: "Bartek's Lamp",
	})

	a.Switch.On.OnValueRemoteUpdate(func(on bool) {
		var newValue string
		if on {
			log.Println("Switch is on")
			newValue = fhome.Value100
		} else {
			log.Println("Switch is off")
			newValue = fhome.Value0
		}

		err := client.XEvent(291, newValue)
		if err != nil {
			log.Fatalf("failed to send event with value %s: %v\n", newValue, err)
		}
	})

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	// Create the hap server.
	server, err := hap.NewServer(fs, a.A)
	if err != nil {
		// stop if an error happens
		log.Panic(err)
	}

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()

	// Run the server.
	server.ListenAndServe(ctx)
}
