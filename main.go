package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	log "github.com/Sirupsen/logrus"
	rancher "github.com/rancher/go-rancher/v2"

	"github.com/urfave/cli"
)

var (
	VERSION = "v0.0.0-dev"
)

var withoutPagination *rancher.ListOpts

func main() {
	app := cli.NewApp()
	app.Name = "ranch53"
	app.Version = VERSION
	app.Usage = "ranch53"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug,d",
		},
		cli.StringFlag{
			Name:  "rancher-url",
			Value: "http://localhost:8080/",
			Usage: "Provide full URL of rancher server",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:  "rancher-access-key",
			Usage: "Rancher Access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:  "rancher-secret-key",
			Usage: "Rancher Secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:  "aws-access-key-id",
			Usage: "AWS Access Key ID",
			EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:  "aws-secret-access-key",
			Usage: "AWS Secret Access Key",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
		cli.IntFlag{
			Name: "poll-interval,t",
			Usage: "Polling interval",
			EnvVar: "POLL_INTERVAL",
			Value: 0,
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

	// ensure that we have been provided aws credentials
	if len(c.String("aws-access-key-id")) < 1 {
		log.Errorf("aws access key id not provided, exiting")
		os.Exit(1)
	}
	if len(c.String("aws-secret-access-key")) < 1 {
		log.Errorf("aws secret key not provided, exiting")
		os.Exit(1)
	}

	// create the rancher client
	rancherClient := createClient(c.String("rancher-url"),
	                              c.String("rancher-access-key"),
	                              c.String("rancher-secret-key"))

	// create the aws session
	awsSession, err := session.NewSession()
	if err != nil {
		log.Error("failed to create aws session", err)
	}
	r53 := route53.New(awsSession)

	// the integration junction magic factory entrypoint
	if c.Int("poll-interval") > 0 {
		log.Debug("poll interval: ", c.Int("poll-interval"))
		for {
			discover(rancherClient, r53)
			time.Sleep(time.Duration(c.Int("poll-interval")) * (time.Millisecond * 1000))
		}
	} else {
		discover(rancherClient, r53)
	}

	return nil
}

func createClient(rancherURL, accessKey, secretKey string) (*rancher.RancherClient) {
	client, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       rancherURL,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Timeout:   time.Second * 8,
	})

	if err != nil {
		log.Errorf("Failed to create a client for rancher, error: %s", err)
		os.Exit(1)
	}

	return client
}

func listRancherLoadBalancerServices(client *rancher.RancherClient) []*rancher.LoadBalancerService {
	var servicesList []*rancher.LoadBalancerService

	services, err := client.LoadBalancerService.List(withoutPagination)

	if err != nil {
		log.Errorf("cannot get services: %+v", err)
	}

	for k := range services.Data {
		servicesList = append(servicesList, &services.Data[k])
	}

	return servicesList
}

func createCNAME(svc *route53.Route53, zoneId string, name string, target string) {
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(name),
						Type: aws.String("CNAME"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(target),
							},
						},
						TTL:           aws.Int64(int64(60)),
						Weight:        aws.Int64(int64(1)),
						SetIdentifier: aws.String("Arbitrary Id describing this change set"),
					},
				},
			},
			Comment: aws.String("Sample update."),
		},
		HostedZoneId: aws.String(zoneId),
	}
	resp, err := svc.ChangeResourceRecordSets(params)

	if err != nil {
		log.Error(fmt.Println(err.Error()))
		return
	}

	log.Info(resp)
}

func discover(rancherClient *rancher.RancherClient, r53 *route53.Route53) {
	loadBalancerServices := listRancherLoadBalancerServices(rancherClient)

	for _, s := range loadBalancerServices {
		for k, v := range s.LaunchConfig.Labels {
			// this lb should be in a r53 zone
			if k == "r53_zone_id" {
				log.Debug("found service with r53 zone: ", s, k, v)

				// assign the zone and get the other label values, params
				zoneId := fmt.Sprint(v)
				dnsName := fmt.Sprint(s.LaunchConfig.Labels["dns_name"])
				dnsTarget := fmt.Sprint(s.LaunchConfig.Labels["dns_target"])

				log.WithFields(log.Fields{
					"zone_id": zoneId,
					"name": dnsName,
					"target": dnsTarget,
				}).Info("update record")

				// todo: only update the record if changed/needed

				// upsert the record to r53
				createCNAME(r53, zoneId, dnsName, dnsTarget)
			}
		}
	}
}
