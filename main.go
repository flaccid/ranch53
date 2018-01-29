package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
	d "github.com/flaccid/ranch53/discover"
	r "github.com/flaccid/ranch53/rancher"
)

var (
	VERSION = "v0.0.0-dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "ranch53"
	app.Version = VERSION
	app.Usage = "A tool to managing dns specific to rancher resources"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "sync-host-pools",
			Usage: "synchronise host dns pool records",
		},
		cli.BoolFlag{
			Name:  "sync-lb-services",
			Usage: "synchronise load balancer service dns records",
		},
		cli.StringFlag{
			Name:   "rancher-url",
			Value:  "http://localhost:8080/",
			Usage:  "Provide full URL of rancher server",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "rancher-access-key",
			Usage:  "Rancher Access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "rancher-secret-key",
			Usage:  "Rancher Secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "aws-access-key-id",
			Usage:  "AWS Access Key ID",
			EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "aws-secret-access-key",
			Usage:  "AWS Secret Access Key",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
		cli.IntFlag{
			Name:   "poll-interval,t",
			Usage:  "Polling interval",
			EnvVar: "POLL_INTERVAL",
			Value:  0,
		},
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "enable debug logging",
		},
	}

	app.Run(os.Args)
}

func beforeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func start(c *cli.Context) error {
	log.Info("ranch53 starting up")
	log.Debug("debug mode enabled")

	// user must specify at least one action
	if !c.Bool("sync-host-pools") && !c.Bool("sync-lb-services") {
		log.Errorf("you must specify at least one of --sync-host-pools or --sync-lb-services")
		os.Exit(1)
	}

	// ensure that we have been provided aws credentials
	if len(c.String("aws-access-key-id")) < 1 {
		log.Errorf("aws access key id not provided, exiting")
		os.Exit(1)
	}
	if len(c.String("aws-secret-access-key")) < 1 {
		log.Errorf("aws secret key not provided, exiting")
		os.Exit(1)
	}

	log.Debug("rancher-url: ", c.String("rancher-url"))

	// create the rancher client
	rancherClient := r.CreateClient(c.String("rancher-url"),
		c.String("rancher-access-key"),
		c.String("rancher-secret-key"))

	// create the aws session
	awsSession, err := session.NewSession()
	if err != nil {
		log.Error("failed to create aws session", err)
	}
	r53 := route53.New(awsSession)
	// log.Debug("r53 session ", r53)

	// the integration junction magic factory entrypoint
	d.Discover(rancherClient, r53, c.Bool("sync-host-pools"), c.Bool("sync-lb-services"), c.Int("poll-interval"))

	return nil
}
