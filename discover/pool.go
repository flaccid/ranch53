package discover

import (
	"github.com/aws/aws-sdk-go/service/route53"
)

type PoolRecord struct {
	Name    string                    // e.g. rancher-private-pool
	Soa     string                    // e.g. dev.mydomain.com
	ZoneId  string                    // e.g. Z66GRIXASCZABC
	Targets []*route53.ResourceRecord // an array of target IPs
}

func addHostPool(pools []*PoolRecord, pool *PoolRecord) []*PoolRecord {
	hasRecord := false
	for _, p := range pools {
		if p.Name == pool.Name && p.Soa == pool.Soa {
			hasRecord = true
		}
	}
	if hasRecord != true {
		pools = append(pools, pool)
	}
	return pools
}
