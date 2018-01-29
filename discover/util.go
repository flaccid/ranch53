package discover

import (
	"fmt"
	"strings"
)

func getLabelValue(labels map[string]interface{}, key string) string {
	var labelValue interface{}

	for k, v := range labels {
		if k == key {
			labelValue = v
		}
	}

	return fmt.Sprintf("%v", labelValue)
}

func getDomainFromFqdn(fqdn string) string {
	return strings.TrimPrefix(fqdn, strings.Split(fqdn, ".")[0]+".")
}
