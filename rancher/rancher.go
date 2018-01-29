package rancher

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	rancher "github.com/rancher/go-rancher/v2"
)

type Client struct {
	client *rancher.RancherClient
}

var withoutPagination *rancher.ListOpts

func HostHasLabel(h *rancher.Host, name string) bool {
	for k := range h.Labels {
		if k == name {
			return true
		}
	}
	return false
}

func HostHasLabelValue(h *rancher.Host, name string, value string) bool {
	for k, v := range h.Labels {
		if k == name {
			if v == value {
				return true
			}
		}
	}
	return false
}

func CreateClient(url, accessKey, secretKey string) *rancher.RancherClient {
	client, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       url,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Timeout:   time.Second * 5,
	})

	if err != nil {
		log.Errorf("Failed to create a client for rancher, error: %s", err)
		os.Exit(1)
	}

	return client
}

func ListRancherLoadBalancerServices(client *rancher.RancherClient) []*rancher.LoadBalancerService {
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

func ListRancherHosts(client *rancher.RancherClient) []*rancher.Host {
	var hostsList []*rancher.Host

	hosts, err := client.Host.List(withoutPagination)

	if err != nil {
		log.Errorf("cannot get hosts: %+v", err)
	}

	for k := range hosts.Data {
		hostsList = append(hostsList, &hosts.Data[k])
	}

	return hostsList
}

func ListRancherHostsByLabel(client *rancher.RancherClient, label string) []*rancher.Host {
	var hostsList []*rancher.Host

	hosts, err := client.Host.List(withoutPagination)

	if err != nil {
		log.Errorf("cannot get hosts: %+v", err)
	}

	for k := range hosts.Data {
		hostsList = append(hostsList, &hosts.Data[k])
	}

	return hostsList
}
