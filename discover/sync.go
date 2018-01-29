package discover

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"

	log "github.com/Sirupsen/logrus"
	dns "github.com/flaccid/ranch53/dns"
	r "github.com/flaccid/ranch53/rancher"
	rancher "github.com/rancher/go-rancher/v2"
)

func syncHostPools(rancherClient *rancher.RancherClient, r53 *route53.Route53) {
	// get all hosts
	hosts := r.ListRancherHosts(rancherClient)
	log.Debug(hosts)

	var pools []*PoolRecord

	// iterate on all hosts
	for _, h := range hosts {
		// only care about hosts with pool_zone and tier labels
		if r.HostHasLabel(h, "pool_zone") && r.HostHasLabel(h, "tier") {
			soa := getLabelValue(h.Labels, "pool_zone")
			name := "rancher-" + getLabelValue(h.Labels, "tier") + "-pool"
			zoneId, _ := dns.GetHostedZoneId(r53, soa)
			ipAddress := h.AgentIpAddress
			pools = addHostPool(pools, &PoolRecord{Soa: soa, Name: name, ZoneId: zoneId})

			for _, p := range pools {
				if p.Soa == soa && p.Name == name && p.ZoneId == zoneId {
					// we only need to append the targets
					// log.Debug("existing targets", p.Targets)
					p.Targets = append(p.Targets, &route53.ResourceRecord{Value: aws.String(ipAddress)})
					// log.Debug("new targets", p.Targets)
				}
			}
		}
	}

	for _, p := range pools {
		log.WithFields(log.Fields{
			"name":    p.Name,
			"soa":     p.Soa,
			"zoneId":  p.ZoneId,
			"targets": p.Targets,
		}).Info("sync pool")
		dns.CreateA(r53, p.ZoneId, p.Name+"."+p.Soa, p.Targets)
	}
}

func syncLbServices(rancherClient *rancher.RancherClient, r53 *route53.Route53) {
	loadBalancerServices := r.ListRancherLoadBalancerServices(rancherClient)
	// log.Debug(loadBalancerServices)

	for _, s := range loadBalancerServices {
		for k, v := range s.LaunchConfig.Labels {
			if k == "dns_alias" {
				// log.Debug("found service with alias: ", s, k, v)

				// assign the zone and get the other label values, params
				dnsAlias := fmt.Sprint(v)
				dnsTarget := getLabelValue(s.LaunchConfig.Labels, "dns_target")
				soa := getDomainFromFqdn(dnsAlias)
				zoneId, _ := dns.GetHostedZoneId(r53, soa)

				log.WithFields(log.Fields{
					"alias":   dnsAlias,
					"target":  dnsTarget,
					"soa":     soa,
					"zone id": zoneId,
				}).Info("update record")

				// TODO: only update the record if changed/needed

				// upsert the record to route 53
				dns.CreateCNAME(r53, zoneId, dnsAlias, dnsTarget)
			}
		}
	}
}
