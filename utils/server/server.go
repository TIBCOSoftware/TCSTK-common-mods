package server

import "strings"

// GetServer return the LA server for a region
func GetServer(region string, tenant string) (server string) {

	if region == "US" {
		region = ""
	} else {
		region = region + "."
	}

	var tenantURL string
	switch tenant {
	case "bpm":
		tenantURL = "liveapps"
	case "":
		tenantURL = "account"
	default:
		tenantURL = tenant
	}
	return "https://" + strings.ToLower(region) + tenantURL + ".cloud.tibco.com"
}

// GetRegionServer for TIBCO Cloud
func GetRegionServer(region string) (server string) {
	switch region {
	case "EU":
		server = "eu-west-1"
	case "US":
		server = "us-west-2"
	case "AP":
		server = "ap-southeast-2"
	default:
		server = "us-west-2"
	}

	return server

}
