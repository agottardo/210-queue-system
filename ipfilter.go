package main

import (
	"github.com/yl2chen/cidranger"
	"log"
	"net"
	"strings"
)

var UBC_RANGES = []string{"128.189.16.0/20",
	"128.189.128.0/18",
	"128.189.64.0/19",
	"128.189.192.0/18",
	"206.12.40.0/21",
	"206.12.64.0/21",
	"206.12.136.0/21",
	"206.87.0.0/22",
	"206.87.8.0/22",
	"206.87.12.0/22",
	"206.87.16.0/22",
	"206.87.20.0/22",
	"206.87.112.0/21",
	"206.87.120.0/21",
	"206.87.128.0/19",
	"206.87.192.0/21",
	"206.87.208.0/21",
	"206.87.216.0/21",
	"206.87.232.0/21",
	"128.189.96.0/19",
	"137.82.0.0/16",
	"142.103.0.0/16",
	"198.162.32.0/19",
	"206.12.72.0/22",
	"206.12.118.0/24",
	"206.12.208.0/22",
	"206.87.200.0/21",
	"206.87.224.0/21",
	"207.23.94.0/23",
	"142.103.93.0/24",
	"142.103.165.0/24",
	"206.12.52.0/22",
	"2607:F8F0:0610::/48",
	"2607:F8F0:0400::/52",
}

type UBCIPFilter struct {
	Ranger cidranger.Ranger
}

type IPFilter interface {
	IsIPAuthorized(ip net.IP) bool
}

// Initialize returns an instance of UBCIPFilter,
// which can be used to match IP addresses to filter
// incoming web requests.
// Initialize should only be called once, when starting
// the application.
func Initialize() *UBCIPFilter {
	filter := UBCIPFilter{}
	ranger := cidranger.NewPCTrieRanger()
	for _, cidrRange := range UBC_RANGES {
		_, subnet, _ := net.ParseCIDR(cidrRange)
		_ = ranger.Insert(cidranger.NewBasicRangerEntry(*subnet))
	}
	filter.Ranger = ranger
	return &filter
}

// IsIPAuthorized returns true if the given IP belongs
// to the UBC network.
// This function first tries to match the IP against an
// hard-coded list of CIDR IP ranges (UBC_RANGES).
// If this fails, it executes a reverse DNS lookup, and
// checks whether the PTR DNS record has a ubc.ca. suffix.
func (f *UBCIPFilter) IsIPAuthorized(ip net.IP) bool {
	allowed, err := f.Ranger.Contains(ip)
	if err != nil {
		log.Println("Error determining IP authorization status:", err)
		return false
	} else if allowed {
		return true
	}
	ptr, _ := net.LookupAddr(ip.String())
	for _, rDNS := range ptr {
		if strings.HasSuffix(rDNS, ".ubc.ca.") {
			return true
		}
	}
	return allowed
}
