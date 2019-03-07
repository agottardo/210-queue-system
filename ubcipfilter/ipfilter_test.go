package ubcipfilter

import (
	"net"
	"testing"
)

var theFilter = Initialize()

func TestIPFilterAccept(t *testing.T) {
	result := theFilter.IsIPAuthorized(net.ParseIP("206.87.122.200"))
	if !result {
		t.Errorf("Should have allowed IP 206.87.122.200")
	}
}

func TestIPFilterReject(t *testing.T) {
	result := theFilter.IsIPAuthorized(net.ParseIP("8.8.8.8"))
	if result {
		t.Errorf("Should have rejected IP 8.8.8.8")
	}
}
