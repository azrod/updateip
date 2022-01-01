<center><h1>UpdateIP</h1></center>

UpdateIP is a automatic update of a DNS record based on the External IP address

Actual supported DNS providers:

* Amazon Route 53
* Cloudflare
* OVH (Warning: Is not tested)

## How to setup

### **From Docker Run**

```bash
docker run -itd -v "./config.yaml:/config/config.yaml" ghcr.io/azrod/updateip:latest
```

### **From Docker Compose**

```bash
version: '3.8'

services:
  updateip:
    container_name: updateaip
    image: ghcr.io/azrod/updateip:latest
    volumes:
      - ./config.yaml:/config/config.yaml:ro

```

### **From Source**

```bash
go mod tidy 
go build

PATH_CONFIG_DIRECTORY=$(pwd) ./updateip
```

# How to config

Create the **config.yaml** configuration file

```yaml
log:
  level: debug # Available : trace debug info warn error fatal panic
  humanize: true # Default: false

metrics:
  enable: true # Default: false
  port: 8080 # Default : 8080
  host: 0.0.0.0 # Default: 0.0.0.0
  path: /metrics # Default: /metrics

providers:
  aws_account:
    enable: true
    secret:
      access_key_id: "xxx"
      secret_access_key: "xxx"
      region: "eu-west-1"
    record:
      name: "subdomain.domain.com"
      ttl: 60
      domain: "domain.com"

  ovh_account:
    enable: true
    secret:
      application_key: "xxx"
      application_secret: "xxx"
      consumer_key: "xxx"
      region: "eu-west-1"
    record:
      name: "subdomain.domain.com"
      ttl: 60
      zone: "domain.com"

  cloudflare_account:
    enable: true
    secret:
      api_key: "xxx"
      email: "xxx"
    record:
      name: "subdomain.domain.com"
      ttl: 60
      domain: "domain.com"
```
