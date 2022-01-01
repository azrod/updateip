# Advanced - DNS

OpenDNS offer service to get your external IP address by DNS query. UpdateIP use this service to get your external IP address.
This is a representation of the DNS Client in bash. 

```bash
dig +short myip.opendns.com @208.67.222.222
```

```go
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
```

UpdateIP use also the DNS Client to get status of the DNS record.
The configuration of the DNS server used is therefore important.

If you running UpdateIP from a Docker container, you can use the following command to get the configuration of the DNS server used:

On `docker run` command use `--dns x.x.x.x` to set manual DNS server for this container.
On Docker compose use a following parameters to set manual DNS server for this container.

```bash
version: '3.8'

services:
  updateip:
    container_name: updateaip
    image: ghcr.io/azrod/updateip:latest
    volumes:
      - ./config.yaml:/config/config.yaml:ro
    dns:
      - x.x.x.x
    environnement:
      - LOG_LEVEL=debug
      - LOG_HUMANIZE=true
      - METRICS_ENABLE=true
      - METRICS_PORT=8080
      - METRICS_HOST=0.0.0.0
    port:
      - "8080:8080"
    restart: always
```
