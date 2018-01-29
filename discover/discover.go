package discover

import (
	"time"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/rancher/go-rancher/v2"

	log "github.com/Sirupsen/logrus"
)

func sync(rancherClient *client.RancherClient, r53 *route53.Route53, hostPools bool, lbServices bool) {
	if hostPools {
		syncHostPools(rancherClient, r53)
	}
	if lbServices {
		syncLbServices(rancherClient, r53)
	}

	return
}

func Discover(rancherClient *client.RancherClient, r53 *route53.Route53, syncHostPools bool, syncLbServices bool, pollInterval int) {
	log.WithFields(log.Fields{
		"sync host pools":  syncHostPools,
		"sync lb services": syncLbServices,
		"poll interval":    pollInterval,
	}).Debug("init sync loop")

	if pollInterval > 0 {
		for {
			sync(rancherClient, r53, syncHostPools, syncLbServices)
			time.Sleep(time.Duration(pollInterval) * (time.Millisecond * 1000))
		}
	} else {
		log.Debug("perform a once-off sync")
		sync(rancherClient, r53, syncHostPools, syncLbServices)
	}
}
