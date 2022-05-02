package main

// import testing
import (
	"testing"

	"github.com/azrod/updateip/pkg/ip"
)

func TestGetMyExternalIP(t *testing.T) {
	ip, err := ip.GetMyExternalIP()
	if err != nil {
		t.Error(err)
	}

	// if ip is empty
	if ip.String() == "" {
		t.Error("IP is empty")
	}

	// if ip is not an IPv4 address
	if ip.To4() == nil {
		t.Error("IP is not an IPv4 address")
	}

	t.Log(ip)
}
