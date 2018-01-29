package dns

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/route53"

	log "github.com/Sirupsen/logrus"
)

func GetHostedZoneId(r *route53.Route53, rootDomainName string) (string, error) {
	soa := rootDomainName + "."
	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(soa),
		MaxItems: aws.String("1"),
	}
	log.Debug(params)

	zones, err := r.ListHostedZonesByName(params)
	if err != nil {
		log.Errorf("could not list hosted zones: %v", err)
	}
	log.Debug(zones)

	if len(zones.HostedZones) == 0 || *zones.HostedZones[0].Name != soa {
		log.Errorf("hosted zone for '%s' not found", rootDomainName)
	}

	return strings.TrimPrefix(*zones.HostedZones[0].Id, "/hostedzone/"), nil
}

func CreateCNAME(r *route53.Route53, zoneId string, name string, target string) {
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
			Comment: aws.String("updated by ranch53."),
		},
		HostedZoneId: aws.String(zoneId),
	}
	resp, err := r.ChangeResourceRecordSets(params)

	if err != nil {
		log.Error(fmt.Println(err.Error()))
		return
	}

	log.Info(resp)
}

func CreateA(r *route53.Route53, zoneId string, name string, target []*route53.ResourceRecord) {
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            aws.String(name),
						Type:            aws.String("A"),
						ResourceRecords: target,
						TTL:             aws.Int64(int64(60)),
					},
				},
			},
			Comment: aws.String("rancher pool"),
		},
		HostedZoneId: aws.String(zoneId),
	}

	result, err := r.ChangeResourceRecordSets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeNoSuchHealthCheck:
				fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
			case route53.ErrCodeInvalidChangeBatch:
				fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			case route53.ErrCodePriorRequestNotComplete:
				fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
