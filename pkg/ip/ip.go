package ip

import (
	"context"
	"net"
	"time"
)

// GetMyExternalIP returns the external IP address of the local machine.
func GetMyExternalIP() (ip net.IP, err error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "208.67.222.222:53")
		},
	}
	xIP, err := r.LookupHost(context.Background(), "myip.opendns.com")
	if err != nil {
		return net.IP{}, err
	}
	return net.ParseIP(xIP[0]), err
}
